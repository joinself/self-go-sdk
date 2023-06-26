# Self Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/joinself/self-go-sdk.svg)](https://pkg.go.dev/github.com/joinself/self-go-sdk)
[![CI](https://github.com/joinself/self-go-sdk/actions/workflows/ci.yml/badge.svg)](https://github.com/joinself/self-go-sdk/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/joinself/self-go-sdk)](https://goreportcard.com/report/github.com/joinself/self-go-sdk)

The official Self SDK for Go.

This SDK provides access to the following self services:

- **Authentication**: For authenticating users
- **Identity**: For looking up identities, apps, devices, public keys
- **Fact**: For requesting information from identities or intemediaries
- **Messaging**: For building services to interact with other entities

## Installation

### Dependencies

- [Go](https://go.dev) 1.13 or later
- [Self OMEMO](https://github.com/joinself/self-omemo)

##### Debian/Ubuntu
```bash
curl -LO https://github.com/joinself/self-omemo/releases/download/0.4.0/self-omemo_0.4.0_amd64.deb
apt install -y ./self-omemo_0.4.0_amd64.deb
```

##### CentOS/Fedora/RedHat
```bash
rpm -Uvh https://github.com/joinself/self-omemo/releases/download/0.4.0/self-omemo-0.4.0-1.x86_64.rpm
```

##### MacOS - AMD64
```bash
brew tap joinself/crypto
brew install libself-omemo
```

##### MacOS - ARM64
Brew on M1 macs currently lacks environment variables needed for the SDK to find the `omemo` library, so you will need to add some additional configuration to your system:

In your `~/.zshrc`, add:
```bash
export C_INCLUDE_PATH=/opt/homebrew/include/
export LIBRARY_PATH=$LIBRARY_PATH:/opt/homebrew/lib
```

You should then be able to run:

```bash
source ~/.zshrc
brew tap joinself/crypto
brew install --build-from-source libself-omemo
```

Note, you may also need to create `/usr/local/lib` if it does not exist:
```bash
sudo mkdir /usr/local/lib
```

### Install

```bash
go get github.com/joinself/self-go-sdk
```

## Usage

### Register Application

Before the SDK can be used you must first register an application on the Self Developer Portal. Once registered, the portal will generate credentials for the application that the SDK will use to authenticate against the Self network.

Self provides two isolated networks:

[Developer Portal (production network)](https://developer.joinself.com) - Suitable for production services  
[Developer Portal (sandbox network)](https://developer.sandbox.joinself.com) - Suitable for testing and experimentation

Register your application using one of the links above ([further information](https://docs.joinself.com/quickstart/app-setup/)).

### Examples

#### Client setup

```go
import "github.com/joinself/self-go-sdk"

func main() {
    cfg := selfsdk.Config{
        SelfAppID:           "<application-id>",
        SelfAppDeviceSecret: "<application-secret-key>",
        StorageDir:          "/data",
        StorageKey:          "random-secret-string",
        Environment:         "sandbox",  // optional (defaults to production)
    }

    client, err := selfsdk.New(cfg)
    client.Start()
}
```

#### Identities

The identities service provides functionality for looking up identities, devices and public keys.

To query an identity:

```go
import "github.com/joinself/self-go-sdk"

func main() {
    svc := client.IdentityService()

    identity, err := svc.GetIdentity("<self-id>")
    ...
}
```

#### Facts

The fact service can be used to ask for specific attested facts from an identity. These requests can be sent to the identity directly, or via an intermediary if you would prefer not to see the users personal information directly, but would like to know it satisfies a given criteria.

For detailed examples of fact requests:
- [Fact Request](_examples/fact_request/fact.go)
- [Fact Request via an Intermediary](_examples/fact_request_intermediary/fact.go)
- [Fact Request QR code](_examples/fact_request_qr/fact.go)

To directly ask an identity for facts:

```go
import (
    "github.com/joinself/self-go-sdk"
    "github.com/joinself/self-go-sdk/fact"
)

func main() {
    svc := client.FactService()

    req := fact.FactRequest{
        ...
    }

    resp, err := svc.Request(&req)
    ...
}
```

#### Authentication

The authentication service can be used to send an authentication challenge to a users device. The response the user sends will be signed by their identity and can be validated. You can authenticate a client by two means; If you know the self id of the user you wish to authenticate, you can do it directly. Alternatively, if you do not know the identity of the user, you can generate and display a qr code that can be read by the users device.

For detailed examples of authentication requests:
- [Authentication Request](_examples/authentication/authentication.go)
- [Authentication Request QR code](_examples/authentication_qr/authentication.go)

To authenticate a user directly:

```go
import (
    "github.com/joinself/self-go-sdk"
)

func main() {
    svc := client.AuthenticationService()

    err = svc.Request("<self-id>")
    ...
}
```

## Documentation

- [Documentation](https://docs.joinself.com/)
- [Examples](_examples)

## Support

Looking for help? Reach out to us at [support@joinself.zendesk.com](mailto:support@joinself.zendesk.com)

## Contributing

See [Contributing](CONTRIBUTING.md).

## License

See [License](LICENSE).
