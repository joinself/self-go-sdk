package common

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joinself/self-go-sdk-next/account"
	"github.com/joinself/self-go-sdk-next/keypair/signing"
	"github.com/joinself/self-go-sdk-next/message"
)

// to determine which user we are interacting with, we can generate a
// discovery request and encode it to a qr code that your users can scan.
// as we may not have interacted with this user before, we need to prepare
// information they need to establish an encrypted group

type DiscoveryManager struct {
	requests sync.Map
}

func NewDiscoveryManager() *DiscoveryManager {
	return &DiscoveryManager{}
}

func (d *DiscoveryManager) GenerateAndDisplayQRCode(selfAccount *account.Account, inboxAddress *signing.PublicKey) (chan *message.Message, error) {
	// get a key package that the responder can use to create an encryped group
	// for us, if there is not already an existing one.
	keyPackage, err := selfAccount.ConnectionNegotiateOutOfBand(
		inboxAddress,
		time.Now().Add(time.Minute*5),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate key package: %w", err)
	}

	// build the key package into a discovery request
	content, err := message.NewDiscoveryRequest().
		KeyPackage(keyPackage).
		Expires(time.Now().Add(time.Minute * 5)).
		Finish()

	if err != nil {
		return nil, fmt.Errorf("failed to build discovery request: %w", err)
	}

	// create a channel to track the response from our qr code
	completer := make(chan *message.Message, 1)

	d.requests.Store(
		hex.EncodeToString(content.ID()),
		completer,
	)

	// encode it as a QR code. This can be encoded as either an SVG
	// for use in rendering on a web page, or Unicode, for encoding
	// in text based environments like a terminal
	qrCode, err := message.NewAnonymousMessage(content).
		EncodeToQR(message.QREncodingUnicode)

	if err != nil {
		return nil, fmt.Errorf("failed to encode anonymous message: %w", err)
	}

	log.Info("scan the qr code to complete the discovery request")
	fmt.Println(string(qrCode))

	log.Info(
		"waiting for response to discovery request",
		"requestId", hex.EncodeToString(content.ID()),
	)

	return completer, nil
}

func (d *DiscoveryManager) HandleDiscoveryResponse(selfAccount *account.Account, msg *message.Message) {
	log.Info(
		"received response to discovery request",
		"from", msg.FromAddress().String(),
		"requestId", hex.EncodeToString(msg.ID()),
	)

	discoveryResponse, err := message.DecodeDiscoveryResponse(msg)
	if err != nil {
		log.Warn("failed to decode discovery response", "error", err)
		return
	}

	completer, ok := d.requests.LoadAndDelete(hex.EncodeToString(discoveryResponse.ResponseTo()))
	if !ok {
		log.Warn(
			"received response to unknown request",
			"requestId", hex.EncodeToString(msg.ID()),
		)
		return
	}

	completer.(chan *message.Message) <- msg
}
