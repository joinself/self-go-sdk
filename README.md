# Self Go SDK

[![GoDoc](https://godoc.org/github.com/joinself/self-go-sdk?status.svg)](https://godoc.org/github.com/joinself/self-go-sdk) [![Go Report Card](https://goreportcard.com/badge/github.com/joinself/self-go-sdk)](https://goreportcard.com/report/github.com/joinself/self-go-sdk) [![Build Status](https://travis-ci.com/joinself/self-go-sdk?branch=master)](https://travis-ci.com/joinself/self-go-sdk)


The official self sdk for golang.

## Overview


This sdk provides access to the following self services:

- **Authentication**: For authenticating users
- **Identity**: For looking up identities, apps, devices, public keys
- **Fact**: For requesting information from identities or intemediaries
- **Messaging**: For building services to interact with other entities

## Requirements

- [libself_olm](github.com/joinself/olm)
- [libself_omemo](github.com/joinself/omemo)
- [go 1.13 or higher](golang.org)


#### Debian/Ubuntu
```sh
$ curl -O http://download.joinself.com/olm/libself-olm_0.1.17_amd64.deb
$ curl -O http://download.joinself.com/omemo/libself-omemo_0.1.2_amd64.deb
$ apt install libsodium-dev
$ apt install ./libself-olm_0.1.17_amd64.deb ./libself-omemo_0.1.2_amd64.deb
```

#### Redhat/Centos
```sh
$ rpm -Uvh http://download.joinself.com/olm/libself-olm-0.1.14-1.x86_64.rpm
$ rpm -Uvh http://download.joinself.com/omemo/libself-omemo-0.1.2-1.x86_64.rpm
```

#### Mac
```sh
$ brew tap joinself/crypto
$ brew install libself_olm libself_omemo
```

## Quick Start

To install the go sdk:
```sh
$ go get github.com/joinself/self-go-sdk
```


Using the credentials obtained from the [developer portal](developer.joinself.com), you can configure a new self client as follows:

```go
import "github.com/joinself/self-go-sdk"

func main() {
    cfg := selfsdk.Config{
		SelfAppID:           os.Getenv("SELF_APP_ID"),
		SelfAppDeviceSecret: os.Getenv("SELF_APP_DEVICE_SECRET"),
		StorageDir:          "/opt/self/crypto",
		StorageKey:          "my-secret-crypto-storage-key",
	}

    client, err := selfsdk.New(cfg)
}
```

## Identities

The identities service provides functionality for looking up identities, devices and public keys.

To query an identity:

```go
import "github.com/joinself/self-go-sdk"

func main() {
    svc := client.Identities()

    identity, err := svc.GetIdentity("selfID")
    ...
}
```

## Facts

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
    client, err := selfsdk.New(cfg)
    ...

    svc := client.Facts()

    req := fact.FactRequest{
        ...
    }

    resp, err := svc.Request(&req)
    ...
}
```

## Authentication

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
    client, err := selfsdk.New("appID", "privateKey")
    ...

    svc := client.Authentication()

    err = svc.Request("selfID")
    ...
}
```

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/joinself/self-go-sdk.


## License

The gem is available as open source under the terms of the [MIT License](LICENSE).
