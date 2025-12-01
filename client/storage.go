package client

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/joinself/self-go-sdk/keypair/signing"
)

// StorageItem represents a stored item with metadata
type StorageItem struct {
	Key       string
	Value     []byte
	ExpiresAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Storage handles key-value storage functionality
type Storage struct {
	client *Client
	mu     sync.RWMutex
}

// newStorage creates a new storage component
func newStorage(client *Client) *Storage {
	return &Storage{
		client: client,
	}
}

// Store stores a value with the given key
func (s *Storage) Store(key string, value []byte) error {
	if s.client.isClosed() {
		return ErrClientClosed
	}

	return s.client.account.ValueStore(key, value)
}

// StoreWithExpiry stores a value with the given key and expiry time
func (s *Storage) StoreWithExpiry(key string, value []byte, expires time.Time) error {
	if s.client.isClosed() {
		return ErrClientClosed
	}

	return s.client.account.ValueStoreWithExpiry(key, value, expires)
}

// StoreString stores a string value
func (s *Storage) StoreString(key, value string) error {
	return s.Store(key, []byte(value))
}

// StoreStringWithExpiry stores a string value with expiry
func (s *Storage) StoreStringWithExpiry(key, value string, expires time.Time) error {
	return s.StoreWithExpiry(key, []byte(value), expires)
}

// StoreJSON stores a JSON-serializable value
func (s *Storage) StoreJSON(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return s.Store(key, data)
}

// StoreJSONWithExpiry stores a JSON-serializable value with expiry
func (s *Storage) StoreJSONWithExpiry(key string, value interface{}, expires time.Time) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return s.StoreWithExpiry(key, data, expires)
}

// Lookup retrieves a value by key
func (s *Storage) Lookup(key string) ([]byte, error) {
	if s.client.isClosed() {
		return nil, ErrClientClosed
	}

	return s.client.account.ValueLookup(key)
}

// LookupString retrieves a string value by key
func (s *Storage) LookupString(key string) (string, error) {
	data, err := s.Lookup(key)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LookupJSON retrieves and unmarshals a JSON value by key
func (s *Storage) LookupJSON(key string, target interface{}) error {
	data, err := s.Lookup(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

// Exists checks if a key exists in storage
func (s *Storage) Exists(key string) bool {
	_, err := s.Lookup(key)
	return err == nil
}

// Delete removes a value by key
func (s *Storage) Delete(key string) error {
	if s.client.isClosed() {
		return ErrClientClosed
	}

	return s.client.account.ValueRemove(key)
}

// StoreTemporary stores a value with a relative expiry duration
func (s *Storage) StoreTemporary(key string, value []byte, duration time.Duration) error {
	expires := time.Now().Add(duration)
	return s.StoreWithExpiry(key, value, expires)
}

// StoreTemporaryString stores a string value with a relative expiry duration
func (s *Storage) StoreTemporaryString(key, value string, duration time.Duration) error {
	expires := time.Now().Add(duration)
	return s.StoreStringWithExpiry(key, value, expires)
}

// StoreTemporaryJSON stores a JSON value with a relative expiry duration
func (s *Storage) StoreTemporaryJSON(key string, value interface{}, duration time.Duration) error {
	expires := time.Now().Add(duration)
	return s.StoreJSONWithExpiry(key, value, expires)
}

// StorageNamespace provides namespaced storage operations
type StorageNamespace struct {
	storage   *Storage
	namespace string
}

// Namespace returns a namespaced storage interface
func (s *Storage) Namespace(namespace string) *StorageNamespace {
	return &StorageNamespace{
		storage:   s,
		namespace: namespace,
	}
}

// namespacedKey creates a namespaced key
func (sn *StorageNamespace) namespacedKey(key string) string {
	return fmt.Sprintf("%s:%s", sn.namespace, key)
}

// Store stores a value in the namespace
func (sn *StorageNamespace) Store(key string, value []byte) error {
	return sn.storage.Store(sn.namespacedKey(key), value)
}

// StoreWithExpiry stores a value with expiry in the namespace
func (sn *StorageNamespace) StoreWithExpiry(key string, value []byte, expires time.Time) error {
	return sn.storage.StoreWithExpiry(sn.namespacedKey(key), value, expires)
}

// StoreString stores a string value in the namespace
func (sn *StorageNamespace) StoreString(key, value string) error {
	return sn.storage.StoreString(sn.namespacedKey(key), value)
}

// StoreJSON stores a JSON value in the namespace
func (sn *StorageNamespace) StoreJSON(key string, value interface{}) error {
	return sn.storage.StoreJSON(sn.namespacedKey(key), value)
}

// Lookup retrieves a value from the namespace
func (sn *StorageNamespace) Lookup(key string) ([]byte, error) {
	return sn.storage.Lookup(sn.namespacedKey(key))
}

// LookupString retrieves a string value from the namespace
func (sn *StorageNamespace) LookupString(key string) (string, error) {
	return sn.storage.LookupString(sn.namespacedKey(key))
}

// LookupJSON retrieves a JSON value from the namespace
func (sn *StorageNamespace) LookupJSON(key string, target interface{}) error {
	return sn.storage.LookupJSON(sn.namespacedKey(key), target)
}

// Exists checks if a key exists in the namespace
func (sn *StorageNamespace) Exists(key string) bool {
	return sn.storage.Exists(sn.namespacedKey(key))
}

// Delete removes a value from the namespace
func (sn *StorageNamespace) Delete(key string) error {
	return sn.storage.Delete(sn.namespacedKey(key))
}

// StoreTemporary stores a value with relative expiry in the namespace
func (sn *StorageNamespace) StoreTemporary(key string, value []byte, duration time.Duration) error {
	return sn.storage.StoreTemporary(sn.namespacedKey(key), value, duration)
}

// StoreJSONWithExpiry stores a JSON value with expiry in the namespace
func (sn *StorageNamespace) StoreJSONWithExpiry(key string, value interface{}, expires time.Time) error {
	return sn.storage.StoreJSONWithExpiry(sn.namespacedKey(key), value, expires)
}

// Cache provides caching functionality with automatic expiry
type Cache struct {
	storage *Storage
	prefix  string
}

// Cache returns a cache interface with the given prefix
func (s *Storage) Cache(prefix string) *Cache {
	return &Cache{
		storage: s,
		prefix:  prefix,
	}
}

// cacheKey creates a cache key with prefix
func (c *Cache) cacheKey(key string) string {
	return fmt.Sprintf("cache:%s:%s", c.prefix, key)
}

// Set stores a value in the cache with default expiry (1 hour)
func (c *Cache) Set(key string, value []byte) error {
	return c.SetWithTTL(key, value, time.Hour)
}

// SetWithTTL stores a value in the cache with custom TTL
func (c *Cache) SetWithTTL(key string, value []byte, ttl time.Duration) error {
	expires := time.Now().Add(ttl)
	return c.storage.StoreWithExpiry(c.cacheKey(key), value, expires)
}

// SetString stores a string value in the cache
func (c *Cache) SetString(key, value string) error {
	return c.Set(key, []byte(value))
}

// SetJSON stores a JSON value in the cache
func (c *Cache) SetJSON(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return c.Set(key, data)
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) ([]byte, error) {
	return c.storage.Lookup(c.cacheKey(key))
}

// GetString retrieves a string value from the cache
func (c *Cache) GetString(key string) (string, error) {
	data, err := c.Get(key)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetJSON retrieves a JSON value from the cache
func (c *Cache) GetJSON(key string, target interface{}) error {
	data, err := c.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// Has checks if a key exists in the cache
func (c *Cache) Has(key string) bool {
	return c.storage.Exists(c.cacheKey(key))
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) error {
	return c.storage.Delete(c.cacheKey(key))
}

// Internal methods for handling events

func (s *Storage) onConnect() {
	// Connection established - no specific action needed
}

func (s *Storage) onDisconnect(err error) {
	// Connection lost - no specific action needed
}

func (s *Storage) onWelcome(from *signing.PublicKey, groupAddress *signing.PublicKey) {
	// New connection established - no specific action needed
}

func (s *Storage) onKeyPackage(from *signing.PublicKey) {
	// Key package received - no specific action needed
}

func (s *Storage) onIntroduction(from *signing.PublicKey, tokenCount int) {
	// Introduction received - no specific action needed
}

func (s *Storage) close() {
	// Clean up any resources if needed
}
