// Copyright 2020 Self Group Ltd. All Rights Reserved.

package selfsdk

import (
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/joinself/self-go-sdk/pkg/messaging"
	"github.com/joinself/self-go-sdk/pkg/object"
	"github.com/joinself/self-go-sdk/pkg/pki"
	"github.com/joinself/self-go-sdk/pkg/storage"
	"github.com/joinself/self-go-sdk/pkg/transport"
	"golang.org/x/crypto/ed25519"
)

var (
	defaultAPIURL               = "https://api.joinself.com"
	defaultMessagingURL         = "wss://messaging.joinself.com/v2/messaging"
	defaultReconnectionAttempts = 10
	defaultTCPDeadline          = time.Second * 90
	defaultRequestTimeout       = time.Second * 5
	defaultInboxSize            = 256

	decoder = base64.RawStdEncoding
)

// Connectors stores all connectors for working with different self api's
type Connectors struct {
	Rest           RestTransport
	Websocket      WebsocketTransport
	Messaging      MessagingClient
	PKI            PKIClient
	Storage        Storage
	FileInteractor remoteFile
}

// Config configuration options for the sdk
type Config struct {
	SelfAppID            string
	SelfAppDeviceSecret  string
	StorageKey           string
	DeviceID             string
	StorageDir           string
	APIURL               string
	MessagingURL         string
	Environment          string
	OnConnect            func()
	OnDisconnect         func(err error)
	OnPing               func()
	ReconnectionAttempts int
	TCPDeadline          time.Duration
	RequestTimeout       time.Duration
	Connectors           *Connectors
	kid                  string
	sk                   ed25519.PrivateKey
}

func (c Config) validate() error {
	if c.SelfAppID == "" {
		return errors.New("config must specify the self app id")
	}

	if c.SelfAppDeviceSecret == "" {
		return errors.New("config must specify an app device secret key")
	}

	if len(strings.Split(c.SelfAppDeviceSecret, ":")) < 2 {
		return errors.New("config must specify an app device secret key")
	}

	if c.StorageKey == "" {
		return errors.New("config must specify a key to encrypt storage")
	}

	if c.StorageDir == "" {
		return errors.New("config must specify a storage directory")
	}

	return nil
}

func (c *Config) load() error {
	if strings.Contains(c.SelfAppDeviceSecret, "_") {
		keyParts := strings.Split(c.SelfAppDeviceSecret, "_")
		if keyParts[0] != "sk" {
			return errors.New("the device secret key provided is not valid")
		}
		c.SelfAppDeviceSecret = keyParts[1]
	}

	if c.Connectors == nil {
		c.Connectors = &Connectors{}
	}

	if c.DeviceID == "" {
		c.DeviceID = "1"
	}

	if c.Environment != "" {
		if c.APIURL == "" {
			c.APIURL = "https://api." + c.Environment + ".joinself.com"
		}

		if c.MessagingURL == "" {
			c.MessagingURL = "wss://messaging." + c.Environment + ".joinself.com/v2/messaging"
		}
	}

	if c.APIURL == "" {
		c.APIURL = defaultAPIURL
	}

	if c.MessagingURL == "" {
		c.MessagingURL = defaultMessagingURL
	}

	if c.ReconnectionAttempts != -1 {
		c.ReconnectionAttempts = defaultReconnectionAttempts
	}

	if c.TCPDeadline == 0 {
		c.TCPDeadline = defaultTCPDeadline
	}

	if c.RequestTimeout == 0 {
		c.RequestTimeout = defaultRequestTimeout
	}

	kp := strings.Split(c.SelfAppDeviceSecret, ":")

	skData, err := decoder.DecodeString(kp[1])
	if err != nil {
		return errors.New("could not decode private key")
	}

	c.sk = ed25519.NewKeyFromSeed(skData)
	c.kid = kp[0]

	// loading connectors should be done in order due to dependencies
	err = c.loadRestConnector()
	if err != nil {
		return err
	}

	err = c.loadWebsocketConnector()
	if err != nil {
		return err
	}

	err = c.loadPKIConnector()
	if err != nil {
		return err
	}

	err = c.loadStorageConnector()
	if err != nil {
		return err
	}

	c.loadRemoteFileInteractor()

	return c.loadMessagingConnector()
}

func (c Config) loadRestConnector() error {
	if c.Connectors.Rest != nil {
		return nil
	}

	cfg := transport.RestConfig{
		Client: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   defaultTCPDeadline,
					KeepAlive: defaultTCPDeadline / 2,
				}).Dial,
			},
		},
		APIURL:     c.APIURL,
		SelfID:     c.SelfAppID,
		KeyID:      c.kid,
		PrivateKey: c.sk,
	}

	rest, err := transport.NewRest(cfg)
	if err != nil {
		return err
	}

	c.Connectors.Rest = rest

	return nil
}

func (c Config) loadRemoteFileInteractor() {
	if c.Connectors.FileInteractor == nil {
		c.Connectors.FileInteractor = object.NewRemoteFileInteractor(c.Connectors.Rest)
	}
}

func (c Config) loadWebsocketConnector() error {
	if c.Connectors.Websocket != nil {
		return nil
	}

	cfg := transport.WebsocketConfig{
		MessagingURL: c.MessagingURL,
		SelfID:       c.SelfAppID,
		KeyID:        c.kid,
		DeviceID:     c.DeviceID,
		PrivateKey:   c.sk,
		TCPDeadline:  defaultTCPDeadline,
		InboxSize:    defaultInboxSize,
		OnConnect:    c.OnConnect,
		OnDisconnect: c.OnDisconnect,
		OnPing:       c.OnPing,
	}

	ws, err := transport.NewWebsocket(cfg)
	if err != nil {
		return err
	}

	c.Connectors.Websocket = ws

	return nil
}

func (c Config) loadStorageConnector() error {
	if c.Connectors.Storage != nil {
		return nil
	}

	cfg := storage.Config{
		StorageDir:    c.StorageDir,
		EncryptionKey: c.StorageKey,
		PKI:           c.Connectors.PKI,
	}

	storage, err := storage.New(&cfg)
	if err != nil {
		return err
	}

	c.Connectors.Storage = storage

	return nil
}

func (c Config) loadPKIConnector() error {
	if c.Connectors.PKI != nil {
		return nil
	}

	cfg := pki.Config{
		APIURL:     c.APIURL,
		SelfID:     c.SelfAppID,
		PrivateKey: c.sk,
		Transport:  c.Connectors.Rest,
	}

	client, err := pki.New(cfg)
	if err != nil {
		return err
	}

	c.Connectors.PKI = client

	return nil
}

func (c Config) loadMessagingConnector() error {
	if c.Connectors.Messaging != nil {
		return nil
	}

	cfg := messaging.Config{
		PrivateKey: c.sk,
		Storage:    c.Connectors.Storage,
		Transport:  c.Connectors.Websocket,
	}

	client, err := messaging.New(cfg)
	if err != nil {
		return err
	}

	c.Connectors.Messaging = client

	return nil
}

func (c Config) privateKey() ed25519.PrivateKey {
	return c.sk
}
