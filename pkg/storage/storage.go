package storage

import (
	"crypto/ed25519"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	selfcrypto "github.com/joinself/self-crypto-go"
	"github.com/joinself/self-go-sdk/pkg/siggraph"
	_ "github.com/mattn/go-sqlite3"
)

var (
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
			as_identifier INTEGER NOT NULL,
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
			as_identifier INTEGER NOT NULL,
			with_identifier INTEGER NOT NULL,
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

	// TODO split this so we're not publishing prekeys before we have committed them
	// this will cause one time keys that the account to be recognised by other senders
	s.generateAndPublishOneTimeKeys(inboxID, account)

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

	return txn.Commit()
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
	sessions := make([]*selfcrypto.Session, len(to))

	s.mu.Lock()
	defer s.mu.Unlock()

	txn, err := s.db.Begin()
	if err != nil {
		txn.Rollback()
		return nil, err
	}

	statement := fmt.Sprintf(
		"SELECT with_identifier, olm_session FROM sessions WHERE with_identifier IN(?%s)",
		strings.Repeat(",?", len(to)-1),
	)

	recipients := make([]any, len(to))
	for i := range to {
		recipients[i] = to[i]
	}

	rows, err := txn.Query(statement, recipients...)
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

	err = row.Scan(&sessionPickle)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			txn.Rollback()
			return nil, err
		}

		session, err = s.createInboundSession(txn, from, to, otkm)
		if err != nil {
			txn.Rollback()
			return nil, err
		}
	} else {
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
			session, err = s.createInboundSession(txn, from, to, otkm)
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

	_, err = txn.Exec("UPDATE sessions SET olm_session = ? WHERE as_identifier = ? AND with_identifier = ?;", sessionPickle, to, from)
	if err != nil {
		txn.Rollback()
		return nil, err
	}

	_, err = txn.Exec("UPDATE accounts SET offset = ? WHERE as_identifier = ?;", offset, to)
	if err != nil {
		txn.Rollback()
		return nil, err
	}

	return plaintext, txn.Commit()
}

// Close closes the storage connection
func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) createInboundSession(txn *sql.Tx, from, to string, otkm *selfcrypto.Message) (*selfcrypto.Session, error) {
	row := txn.QueryRow("SELECT olm_account FROM accounts WHERE as_identifier = ?;", to)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var accountPickle string

	err := row.Scan(&accountPickle)
	if err != nil {
		return nil, err
	}

	account, err := selfcrypto.AccountFromPickle(to, s.ec, accountPickle)
	if err != nil {
		return nil, err
	}

	session, err := selfcrypto.CreateInboundSession(account, from, otkm)
	if err != nil {
		return nil, err
	}

	err = account.RemoveOneTimeKeys(session)
	if err != nil {
		return nil, err
	}

	otks, err := account.OneTimeKeys()
	if err != nil {
		return nil, err
	}

	if len(otks.Curve25519) < 10 {
		err = s.generateAndPublishOneTimeKeys(to, account)
		if err != nil {
			return nil, err
		}
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

func (s *Storage) generateAndPublishOneTimeKeys(as string, account *selfcrypto.Account) error {
	var otkb oneTimeKeys

	potks, err := account.OneTimeKeys()
	if err != nil {
		return err
	}

	// TODO configure this based on usecase
	err = account.GenerateOneTimeKeys(100)
	if err != nil {
		return err
	}

	otks, err := account.OneTimeKeys()
	if err != nil {
		return err
	}

	for k, v := range otks.Curve25519 {
		_, ok := potks.Curve25519[k]
		if ok {
			continue
		}
		otkb = append(otkb, oneTimeKey{ID: k, Key: v})
	}

	otkd, err := json.Marshal(otkb)
	if err != nil {
		return err
	}

	identifier, device := idsplit(as)

	return s.pk.SetDeviceKeys(identifier, device, otkd)
}

func idsplit(id string) (string, string) {
	i := strings.Index(id, ":")

	return id[:i], id[i+1:]
}
