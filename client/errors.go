package client

import "errors"

var (
	// Configuration errors
	ErrStorageKeyRequired  = errors.New("storage key is required")
	ErrStoragePathRequired = errors.New("storage path is required")

	// Client state errors
	ErrClientClosed     = errors.New("client is closed")
	ErrClientNotStarted = errors.New("client not started")

	// Discovery errors
	ErrDiscoveryTimeout = errors.New("discovery request timed out")
	ErrInvalidQRCode    = errors.New("invalid QR code")

	// Chat errors
	ErrInvalidPeerDID  = errors.New("invalid peer DID")
	ErrMessageTooLarge = errors.New("message too large")

	// Request errors
	ErrRequestNotFound = errors.New("request not found")
	ErrInvalidResponse = errors.New("invalid response")
)
