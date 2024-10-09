package common

import (
	"github.com/charmbracelet/log"
	"github.com/joinself/self-go-sdk-next/account"
	"github.com/joinself/self-go-sdk-next/message"
)

func NewExamplesDefaultConfig() *account.Config {
	// initialize an account that will be used to interact with self and other entities on
	// the network. the account provides storage of all cryptographic key material, as well
	// as credentials and all state used for e2e encrypted messaging with other entitites
	return &account.Config{
		// provide a secure storage key that will  be used to encrypt your local account
		// state. this should be replaced with a securely generated key!
		StorageKey: make([]byte, 32),
		// provide a storage path to the directory where your local account state will be
		// stored
		StoragePath: "./storage",
		// provide an environment to target [Develop, Sandbox]
		Environment: account.TargetSandbox,
		// provide the level of log granularity [Error, Warn, Info, Debug, Trace]
		LogLevel: account.LogWarn,
		// specify callbacks to handle events
		Callbacks: account.Callbacks{
			// invoked when the messaging socket connects
			OnConnect: func() {
				log.Info("messaging socket connected")
			},
			// invoked when the messaging socket disconnects. if there is no error
			OnDisconnect: func(err error) {
				if err != nil {
					log.Warn("messaging socket disconnected", "error", err)
				} else {
					log.Info("messaging socket disconnected")
				}
			},
			// invoked when there is a response to a discovery request from a new address.
			OnWelcome: func(selfAccount *account.Account, wlc *message.Welcome) {
				// we have received a response to our discovery request that is from a new
				// user/address that we do not have an  end to end encrypted session.
				// accept the invite to join the encrypted group created by the user.
				groupAddres, err := selfAccount.ConnectionAccept(
					wlc.ToAddress(),
					wlc,
				)

				if err != nil {
					log.Warn("failed to accept connection to encrypted group", "error", err.Error())
					return
				}

				log.Info(
					"accepted connection encrypted group",
					"from", wlc.FromAddress().String(),
					"group", groupAddres.String(),
				)
			},
		},
	}
}
