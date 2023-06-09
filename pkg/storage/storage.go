package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	selfcrypto "github.com/joinself/self-crypto-go"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// ErrSessionNotFound session was not found for a given recipient and sender
	ErrSessionNotFound = errors.New("there is no existing session that matches the given sender and recipient")
	// ErrInvalidGroupMessageRecipient a received message is not intended for this identity
	ErrInvalidGroupMessageRecipient = errors.New("group message does not contain a recipient header for this identity")
)

type prekeys []prekey

type prekey struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

// PKI the public key infrastructure provider used to retrieve and store keys
type PKI interface {
	GetHistory(selfID string) ([]json.RawMessage, error)
	GetDeviceKey(selfID, deviceID string) ([]byte, error)
	SetDeviceKeys(selfID, deviceID string, pkb []byte) error
}

// Stoprage the default storage implementation
// based on sqlite
type Storage struct {
	db *sql.DB
	mu sync.Mutex
	pk PKI
	ec string
}

func New(path, encryptionKey string, pki PKI) (*Storage, error) {
	db, err := sql.Open("sqlite3", filepath.Join(path, "self.db"))
	if err != nil {
		return nil, err
	}

	// TODO handle signnals and graceful shutdown

	s := &Storage{
		db: db,
		pk: pki,
		ec: encryptionKey,
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
		CREATE UNIQUE INDEX idx_accounts_as_identifier
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
		CREATE UNIQUE INDEX idx_sessions_with_identifier
		ON sessions (as_identifier, with_identifier);
	`
	_, err := s.db.Exec(sessionTableStatement)

	return err
}

func (s *Storage) AccountCreate(inboxID string, account *selfcrypto.Account) error {
	accountPickle, err := account.Pickle(s.ec)
	if err != nil {
		return err
	}

	txn, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = txn.Exec("INSERT INTO accounts (as_identifier, offset, olm_account) VALUES (?, ?, ?);", inboxID, 0, accountPickle)
	if err != nil {
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
		return err
	}

	account, err := selfcrypto.AccountFromPickle(inboxID, s.ec, accountPickle)
	if err != nil {
		return err
	}

	err = action(account)
	if err != nil {
		return err
	}

	accountPickle, err = account.Pickle(s.ec)
	if err != nil {
		return err
	}

	_, err = txn.Exec("UPDATE accounts SET olm_account = ? WHERE as_identifier = ?", accountPickle, inboxID)
	if err != nil {
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
		return offset, row.Err()
	}

	err = row.Scan(&offset)
	if err != nil {
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
		return nil, err
	}

	statement := fmt.Sprintf(
		"SELECT olm_session FROM sessions WHERE with_identifier IN(?%s)",
		strings.Repeat(",?", len(to)-1),
	)

	rows, err := txn.Query(statement, to)
	if err != nil {
		return nil, err
	}

	var i int

	for rows.Next() {
		var sessionPickle string

		err := rows.Scan(&sessionPickle)
		if err != nil {
			return nil, err
		}

		session, err := selfcrypto.SessionFromPickle(to[i], s.ec, sessionPickle)
		if err != nil {
			return nil, err
		}

		sessions[i] = session

		i++
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	group, err := selfcrypto.CreateGroupSession(from, sessions)
	if err != nil {
		return nil, err
	}

	ciphertext, err := group.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}

	for x := range sessions {
		sessionPickle, err := sessions[x].Pickle(s.ec)
		if err != nil {
			return nil, err
		}

		_, err = txn.Exec("UPDATE sessions SET olm_session = ? WHERE as_identifier = ? AND with_identifier = ?;", sessionPickle, from, to[x])
		if err != nil {
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

	row := txn.QueryRow("SEELCT olm_session FROM sessions WHERE as_identifier = ? AND with_identifier = ?", to, from)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var session *selfcrypto.Session
	var sessionPickle string

	err = row.Scan(&sessionPickle)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		session, err = s.createInboundSession(txn, from, to, otkm)
		if err != nil {
			return nil, err
		}
	} else {
		session, err = selfcrypto.SessionFromPickle(from, s.ec, sessionPickle)
		if err != nil {
			return nil, err
		}

		matches, err := session.MatchesInboundSession(otkm)
		if err != nil {
			return nil, err
		}

		// sender is attempting to renegotiate the session
		// so create a new inbound session as this is a
		// one time key message
		if otkm.Type == 0 && !matches {
			session, err = s.createInboundSession(txn, from, to, otkm)
			if err != nil {
				return nil, err
			}
		}
	}

	gs, err := selfcrypto.CreateGroupSession(to, []*selfcrypto.Session{session})
	if err != nil {
		return nil, err
	}

	plaintext, err := gs.Decrypt(from, ciphertext)
	if err != nil {
		return nil, err
	}

	sessionPickle, err = session.Pickle(s.ec)
	if err != nil {
		return nil, err
	}

	_, err = txn.Exec("UPDATE sessions SET olm_session = ? WHERE as_identifier = ? AND with_identifier = ?;", sessionPickle, to, from)
	if err != nil {
		return nil, err
	}

	return plaintext, txn.Commit()
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
		err = s.generateAndPublishOneTimeKeys(txn, to, account)
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

func (s *Storage) generateAndPublishOneTimeKeys(txn *sql.Tx, as string, account *selfcrypto.Account) error {
	var pkb prekeys

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
		pkb = append(pkb, prekey{ID: k, Key: v})
	}

	pkbd, err := json.Marshal(pkb)
	if err != nil {
		return err
	}

	parts := strings.Split(as, ":")

	return s.pk.SetDeviceKeys(parts[0], parts[1], pkbd)
}
