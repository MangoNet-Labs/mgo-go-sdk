# Mango Blockchain Go SDK

## Introduction

`github.com/MangoNet-Labs/mgo-go-sdk` is the official Go SDK for Mango Blockchain, providing capabilities to interact with the Mango chain, including account management, transaction building, token operations, and more.

## Features

- **Account Management**: Supports `ed25519` and `secp256k1` key pair generation and signing.
- **Transaction Operations**: Provides transaction creation, signing, and submission.
- **On-Chain Data Queries**: Supports querying token balances, events, objects, and more.
- **WebSocket Subscription**: Enables real-time event push notifications.

## Installation

```sh
# Install using go mod
go get github.com/MangoNet-Labs/mgo-go-sdk
```

## Directory Structure

```
├─ account          # Account management
│  ├─ keypair       # Key pair management
│  └─ signer        # Signer 
├─ bcs              # Serialization and deserialization
├─ client           # Core client functionalities
│  ├─ httpconn      # HTTP connection management
│  ├─ wsconn        # WebSocket connection management
├─ config           # Configuration 
├─ model            # Data models
│  ├─ request       # Request data structures
│  └─ response      # Response data structures
├─ test             # Unit tests and usage examples
├─ utils            # Utility functions
```

## Quick Start

### 1. Initialize the Client

```go
package main

import (
    "fmt"
    "github.com/MangoNet-Labs/mgo-go-sdk/client"
	"github.com/MangoNet-Labs/mgo-go-sdk/config"
)

func main() {
    cli := client.NewMgoClient(config.RpcMgoDevnetEndpoint)
    fmt.Println("Mango Client Initialized:", cli)
}
```

### 2. Generate a Key Pair

```go
package main

import (
    "fmt"
    "github.com/MangoNet-Labs/mgo-go-sdk/account/keypair"
)

func main() {
	kp, err := keypair.NewKeypair(config.Ed25519Flag)
    if err != nil {
		log.Fatalf("%v", err)
		return
	}
    fmt.Println("Public Key:", kp.PublicKeyHex())
}
```

### 3. Send a Transaction

```go
package main

import (
    "fmt"
    "github.com/MangoNet-Labs/mgo-go-sdk/client"
    "github.com/MangoNet-Labs/mgo-go-sdk/account/keypair"
)

func main() {
	var ctx = context.Background()
	cli := client.NewMgoClient(config.RpcMgoDevnetEndpoint)
	kp, err := keypair.NewKeypairWithMgoPrivateKey("your mgo privateKey")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	recipient := "recipient mgo addreess"
	tx, err := cli.TransferMgo(ctx, request.TransferMgoRequest{
		Signer:      kp.MgoAddress(),
		Recipient:   recipient,
		MgoObjectId: "sender mgo coin object address",
		GasBudget:   "10000000",
		Amount:      "10000000000",
	})
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	executeRes, err := cli.SignAndExecuteTransactionBlock(ctx, request.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: tx,
		Keypair:     kp,
		Options: request.TransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
		},
		RequestType: "WaitForLocalExecution",
	})
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	log.Println("Digest", executeRes.Digest)
}
```

## Usage Examples

For detailed usage examples, refer to the `test` directory, which includes sample implementations.

## Contributing

If you wish to contribute to the SDK, please follow these steps:

1. Fork this repository.
2. Create a new branch for development.
3. Submit a Pull Request.

## License

[Apache 2.0 license](LICENSE)
