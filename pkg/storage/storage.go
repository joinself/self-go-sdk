package storage

import (
	"crypto/ed25519"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	selfcrypto "github.com/joinself/self-crypto-go"
	"github.com/joinself/self-go-sdk/pkg/siggraph"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// ErrInvalidRecipients returned when an empty list of recipients is provided
	ErrInvalidRecipients = errors.New("recipient list is empty")
	// ErrSessionNotFound session was not found for a given recipient and sender
	ErrSessionNotFound = errors.New("there is no existing session that matches the given sender and recipient")
	// ErrInvalidGroupMessageRecipient a received message is not intended for this identity
	ErrInvalidGroupMessageRecipient = errors.New("group message does not contain a recipient header for this identity")
)

type oneTimeKeys []oneTimeKey

type oneTimeKey struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

// PKI the public key infrastructure provider used to retrieve and store keys
type PKI interface {
	GetHistory(selfID string) ([]json.RawMessage, error)
	GetDeviceKey(selfID, deviceID string) ([]byte, error)
	SetDeviceKeys(selfID, deviceID string, pkb []byte) error
}

// Config messaging configuration for connecting to self messaging
type Config struct {
	StorageDir    string
	EncryptionKey string
	AccountID     string
	PKI           PKI
}

// Stoprage the default storage implementation
// based on sqlite
type Storage struct {
	db *sql.DB
	mu sync.Mutex
	pk PKI
	ec string
}

func New(cfg *Config) (*Storage, error) {
	err := os.MkdirAll(cfg.StorageDir, 0744)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", filepath.Join(cfg.StorageDir, "self.db"))
	if err != nil {
		return nil, err
	}

	// TODO handle signnals and graceful shutdown
	s := &Storage{
		db: db,
		pk: cfg.PKI,
		ec: cfg.EncryptionKey,
	}

	err = s.setPragmas()
	if err != nil {
		return nil, err
	}

	err = s.createAccountsTable()
	if err != nil {
		return nil, err
	}

	err = s.createSessionsTable()
	if err != nil {
		return nil, err
	}

	err = s.migrateLegacyStorage(cfg.StorageDir, cfg.AccountID)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) setPragmas() error {
	pragmaStatement := `
		PRAGMA synchronous = NORMAL;
		PRAGMA journal_mode = WAL;
		PRAGMA temp_store = MEMORY;
	`
	_, err := s.db.Exec(pragmaStatement)

	return err
}

func (s *Storage) createAccountsTable() error {
	sessionTableStatement := `
		CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			as_identifier TEXT NOT NULL,
			offset INTEGER NOT NULL,
			olm_account BLOB NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_as_identifier
		ON accounts (as_identifier);
	`
	_, err := s.db.Exec(sessionTableStatement)

	return err
}

func (s *Storage) createSessionsTable() error {
	// TODO we could deduplicate as_identifier and with_identifier here
	// by creating a record for each on a new identifier table,
	// but this is only temporary
	sessionTableStatement := `
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			as_identifier TEXT NOT NULL,
			with_identifier TEXT NOT NULL,
			olm_session BLOB NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_sessions_with_identifier
		ON sessions (as_identifier, with_identifier);
	`
	_, err := s.db.Exec(sessionTableStatement)

	return err
}

func (s *Storage) AccountCreate(inboxID string, secretKey ed25519.PrivateKey) error {
	txn, err := s.db.Begin()
	if err != nil {
		return err
	}

	// check the account does not already exist
	row := txn.QueryRow("SELECT olm_account FROM accounts WHERE as_identifier = ?;", inboxID)
	if row.Err() != nil {
		txn.Rollback()
		return row.Err()
	}

	var accountPickle string

	err = row.Scan(&accountPickle)
	if err == nil {
		return nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		txn.Rollback()
		return err
	}

	account, err := selfcrypto.AccountFromSeed(inboxID, secretKey.Seed())
	if err != nil {
		txn.Rollback()
		return err
	}

	otks, err := s.generateOneTimeKeys(inboxID, account)
	if err != nil {
		txn.Rollback()
		return err
	}

	accountPickle, err = account.Pickle(s.ec)
	if err != nil {
		txn.Rollback()
		return err
	}

	_, err = txn.Exec("INSERT INTO accounts (as_identifier, offset, olm_account) VALUES (?, ?, ?);", inboxID, 0, accountPickle)
	if err != nil {
		txn.Rollback()
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	// publish the keys after they have been successfuly saved to the db
	// to avoid a situation where keys are published to the network, but
	// forgotten by the account. attempt to retry the upload if it fails
	for i := 0; i < 60; i++ {
		err = s.publishOneTimeKeys(inboxID, otks)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 5)
	}

	return err
}

// AccountExecute executes an action on an account
func (s *Storage) AccountExecute(inboxID string, action func(account *selfcrypto.Account) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	txn, err := s.db.Begin()
	if err != nil {
		return err
	}

	row := txn.QueryRow("SELECT olm_account FROM accounts WHERE as_identifier = ?;", inboxID)
	if row.Err() != nil {
		return row.Err()
	}

	var accountPickle string

	err = row.Scan(&accountPickle)
	if err != nil {
		txn.Rollback()
		return err
	}

	account, err := selfcrypto.AccountFromPickle(inboxID, s.ec, accountPickle)
	if err != nil {
		txn.Rollback()
		return err
	}

	err = action(account)
	if err != nil {
		txn.Rollback()
		return err
	}

	accountPickle, err = account.Pickle(s.ec)
	if err != nil {
		txn.Rollback()
		return err
	}

	_, err = txn.Exec("UPDATE accounts SET olm_account = ? WHERE as_identifier = ?", accountPickle, inboxID)
	if err != nil {
		txn.Rollback()
		return err
	}

	return txn.Commit()
}

// AccountOffset returns the latest offset for an account that messages will be resumed from
func (s *Storage) AccountOffset(inboxID string) (int64, error) {
	var offset int64

	txn, err := s.db.Begin()
	if err != nil {
		return offset, err
	}

	row := txn.QueryRow("SELECT offset FROM accounts WHERE as_identifier = ?;", inboxID)
	if row.Err() != nil {
		txn.Rollback()
		return offset, row.Err()
	}

	err = row.Scan(&offset)
	if err != nil {
		txn.Rollback()
		return offset, err
	}

	return offset, txn.Commit()
}

func (s *Storage) Encrypt(from string, to []string, plaintext []byte) ([]byte, error) {
	if len(to) < 1 {
		return nil, ErrInvalidRecipients
	}

	sessions := make([]*selfcrypto.Session, len(to))

	s.mu.Lock()
	defer s.mu.Unlock()

	txn, err := s.db.Begin()
	if err != nil {
		txn.Rollback()
		return nil, err
	}

	statement := fmt.Sprintf(
		"SELECT with_identifier, olm_session FROM sessions WHERE as_identifier = ? AND with_identifier IN(?%s)",
		strings.Repeat(",?", len(to)-1),
	)

	args := make([]any, len(to)+1)
	args[0] = from

	for i := range to {
		args[i+1] = to[i]
	}

	rows, err := txn.Query(statement, args...)
	if err != nil {
		txn.Rollback()
		return nil, err
	}

	foundSessions := make(map[string]*selfcrypto.Session)

	for rows.Next() {
		var with string
		var sessionPickle string

		err := rows.Scan(&with, &sessionPickle)
		if err != nil {
			txn.Rollback()
			return nil, err
		}

		session, err := selfcrypto.SessionFromPickle(with, s.ec, sessionPickle)
		if err != nil {
			txn.Rollback()
			return nil, err
		}

		foundSessions[with] = session
	}

	if rows.Err() != nil {
		txn.Rollback()
		return nil, rows.Err()
	}

	for i := range to {
		session, ok := foundSessions[to[i]]
		if !ok {
			session, err = s.createOutboundSession(txn, from, to[i])
			if err != nil {
				txn.Rollback()
				return nil, err
			}
		}

		sessions[i] = session
	}

	group, err := selfcrypto.CreateGroupSession(from, sessions)
	if err != nil {
		txn.Rollback()
		return nil, err
	}

	ciphertext, err := group.Encrypt(plaintext)
	if err != nil {
		txn.Rollback()
		return nil, err
	}

	for i := range sessions {
		sessionPickle, err := sessions[i].Pickle(s.ec)
		if err != nil {
			txn.Rollback()
			return nil, err
		}

		_, ok := foundSessions[to[i]]
		if ok {
			_, err = txn.Exec("UPDATE sessions SET olm_session = ? WHERE as_identifier = ? AND with_identifier = ?;", sessionPickle, from, to[i])
		} else {
			_, err = txn.Exec("INSERT INTO sessions (as_identifier, with_identifier, olm_session) VALUES (?, ?, ?);", from, to[i], sessionPickle)
		}

		if err != nil {
			txn.Rollback()
			return nil, err
		}
	}

	return ciphertext, txn.Commit()
}

func (s *Storage) Decrypt(from, to string, offset int64, ciphertext []byte) ([]byte, error) {
	var gm selfcrypto.GroupMessage

	err := json.Unmarshal(ciphertext, &gm)
	if err != nil {
		return nil, err
	}

	otkm, ok := gm.Recipients[to]
	if !ok {
		return nil, ErrInvalidGroupMessageRecipient
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	txn, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	row := txn.QueryRow("SELECT olm_session FROM sessions WHERE as_identifier = ? AND with_identifier = ?", to, from)
	if row.Err() != nil {
		txn.Rollback()
		return nil, row.Err()
	}

	var session *selfcrypto.Session
	var sessionPickle string
	var sessionExisting bool
	var otks []byte

	err = row.Scan(&sessionPickle)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			txn.Rollback()
			return nil, err
		}

		session, otks, err = s.createInboundSession(txn, from, to, otkm)
		if err != nil {
			txn.Rollback()
			return nil, err
		}
	} else {
		sessionExisting = true

		session, err = selfcrypto.SessionFromPickle(from, s.ec, sessionPickle)
		if err != nil {
			txn.Rollback()
			return nil, err
		}

		matches, err := session.MatchesInboundSession(otkm)
		if err != nil {
			txn.Rollback()
			return nil, err
		}

		// sender is attempting to renegotiate the session
		// so create a new inbound session as this is a
		// one time key message
		if otkm.Type == 0 && !matches {
			session, otks, err = s.createInboundSession(txn, from, to, otkm)
			if err != nil {
				txn.Rollback()
				return nil, err
			}
		}
	}

	gs, err := selfcrypto.CreateGroupSession(to, []*selfcrypto.Session{session})
	if err != nil {
		txn.Rollback()
		return nil, err
	}

	plaintext, err := gs.Decrypt(from, ciphertext)
	if err != nil {
		txn.Rollback()
		return nil, err
	}

	sessionPickle, err = session.Pickle(s.ec)
	if err != nil {
		txn.Rollback()
		return nil, err
	}

	if sessionExisting {
		_, err = txn.Exec("UPDATE sessions SET olm_session = ? WHERE as_identifier = ? AND with_identifier = ?;", sessionPickle, to, from)
	} else {
		_, err = txn.Exec("INSERT INTO sessions (as_identifier, with_identifier, olm_session) VALUES (?, ?, ?);", to, from, sessionPickle)
	}

	if err != nil {
		txn.Rollback()
		return nil, err
	}

	_, err = txn.Exec("UPDATE accounts SET offset = ? WHERE as_identifier = ?;", offset, to)
	if err != nil {
		txn.Rollback()
		return nil, err
	}

	err = txn.Commit()
	if err != nil {
		return nil, err
	}

	return plaintext, s.publishOneTimeKeys(to, otks)
}

// Close closes the storage connection
func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) createInboundSession(txn *sql.Tx, from, to string, otkm *selfcrypto.Message) (*selfcrypto.Session, []byte, error) {
	row := txn.QueryRow("SELECT olm_account FROM accounts WHERE as_identifier = ?;", to)
	if row.Err() != nil {
		return nil, nil, row.Err()
	}

	var accountPickle string

	err := row.Scan(&accountPickle)
	if err != nil {
		return nil, nil, err
	}

	account, err := selfcrypto.AccountFromPickle(to, s.ec, accountPickle)
	if err != nil {
		return nil, nil, err
	}

	session, err := selfcrypto.CreateInboundSession(account, from, otkm)
	if err != nil {
		return nil, nil, err
	}

	err = account.RemoveOneTimeKeys(session)
	if err != nil {
		return nil, nil, err
	}

	otks, err := account.OneTimeKeys()
	if err != nil {
		return nil, nil, err
	}

	var otkd []byte

	if len(otks.Curve25519) < 10 {
		otkd, err = s.generateOneTimeKeys(to, account)
		if err != nil {
			return nil, nil, err
		}
	}

	accountPickle, err = account.Pickle(s.ec)
	if err != nil {
		return nil, nil, err
	}

	_, err = txn.Exec("UPDATE accounts SET olm_account = ? WHERE as_identifier = ?", accountPickle, to)
	if err != nil {
		return nil, nil, err
	}

	return session, otkd, nil
}

func (s *Storage) createOutboundSession(txn *sql.Tx, from, to string) (*selfcrypto.Session, error) {
	identifier, device := idsplit(to)

	var otk oneTimeKey

	otkd, err := s.pk.GetDeviceKey(identifier, device)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(otkd, &otk)
	if err != nil {
		return nil, err
	}

	history, err := s.pk.GetHistory(identifier)
	if err != nil {
		return nil, err
	}

	sg, err := siggraph.New(history)
	if err != nil {
		return nil, err
	}

	pkd, err := sg.ActiveDevice(device)
	if err != nil {
		return nil, err
	}

	pkr, err := selfcrypto.Ed25519PKToCurve25519(pkd)
	if err != nil {
		return nil, err
	}

	pk := base64.RawStdEncoding.EncodeToString(pkr)

	row := txn.QueryRow("SELECT olm_account FROM accounts WHERE as_identifier = ?;", from)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var accountPickle string

	err = row.Scan(&accountPickle)
	if err != nil {
		return nil, err
	}

	account, err := selfcrypto.AccountFromPickle(to, s.ec, accountPickle)
	if err != nil {
		return nil, err
	}

	session, err := selfcrypto.CreateOutboundSession(account, to, pk, otk.Key)
	if err != nil {
		return nil, err
	}

	accountPickle, err = account.Pickle(s.ec)
	if err != nil {
		return nil, err
	}

	_, err = txn.Exec("UPDATE accounts SET olm_account = ? WHERE as_identifier = ?", accountPickle, to)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Storage) generateOneTimeKeys(as string, account *selfcrypto.Account) ([]byte, error) {
	var otkb oneTimeKeys

	potks, err := account.OneTimeKeys()
	if err != nil {
		return nil, err
	}

	// TODO configure this based on usecase
	err = account.GenerateOneTimeKeys(100)
	if err != nil {
		return nil, err
	}

	otks, err := account.OneTimeKeys()
	if err != nil {
		return nil, err
	}

	for k, v := range otks.Curve25519 {
		_, ok := potks.Curve25519[k]
		if ok {
			continue
		}
		otkb = append(otkb, oneTimeKey{ID: k, Key: v})
	}

	return json.Marshal(otkb)
}

func (s *Storage) publishOneTimeKeys(as string, otks []byte) error {
	if len(otks) < 1 {
		return nil
	}
	identifier, device := idsplit(as)
	return s.pk.SetDeviceKeys(identifier, device, otks)
}

func (s *Storage) migrateLegacyStorage(dir, accountID string) error {
	basePath := strings.Replace(dir, "identities", "apps", 1)

	_, err := os.Stat(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// TODO rethink key revocation and consider not allowing a device to rotate keys
	// in favour of revoking the device and creating a new one
	accounts := make(map[string]*struct {
		account  string
		offset   int64
		sessions []struct {
			with    string
			session string
		}
	})

	// check for any files stored in the old structure and move them into the new database
	err = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		switch filepath.Ext(info.Name()) {
		case ".offset":
			fn := strings.TrimSuffix(info.Name(), ".offset")

			od, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			offset, err := strconv.Atoi(string(od[:19]))
			if err != nil {
				return err
			}

			account, ok := accounts[fn]
			if ok {
				account.offset = int64(offset)
				return nil
			}

			accounts[fn] = &struct {
				account  string
				offset   int64
				sessions []struct {
					with    string
					session string
				}
			}{
				offset: int64(offset),
			}
		case ".pickle":
			fn := strings.TrimSuffix(info.Name(), ".pickle")

			pd, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if fn == "account" {
				account, ok := accounts[accountID]
				if ok {
					account.account = string(pd)
					return nil
				}

				accounts[accountID] = &struct {
					account  string
					offset   int64
					sessions []struct {
						with    string
						session string
					}
				}{
					account: string(pd),
				}
			} else {
				account, ok := accounts[accountID]
				if ok {
					account.sessions = append(account.sessions, struct {
						with    string
						session string
					}{
						with:    strings.TrimSuffix(fn, "-session"),
						session: string(pd),
					})
					return nil
				}

				accounts[accountID] = &struct {
					account  string
					offset   int64
					sessions []struct {
						with    string
						session string
					}
				}{
					account: string(pd),
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	txn, err := s.db.Begin()
	if err != nil {
		return err
	}

	for inboxID, account := range accounts {
		_, err = txn.Exec("INSERT INTO accounts (as_identifier, offset, olm_account) VALUES (?, ?, ?);", inboxID, account.offset, account.account)
		if err != nil {
			txn.Rollback()
			return err
		}

		for _, session := range account.sessions {
			_, err = txn.Exec("INSERT INTO sessions (as_identifier, with_identifier, olm_session) VALUES (?, ?, ?);", inboxID, session.with, session.session)
			if err != nil {
				txn.Rollback()
				return err
			}
		}
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return os.Rename(basePath, basePath+"-depreciated")
}

func idsplit(id string) (string, string) {
	i := strings.Index(id, ":")
	return id[:i], id[i+1:]
}
