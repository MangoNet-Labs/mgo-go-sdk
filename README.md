# Mango Blockchain Go SDK

## Introduction

`github.com/mangonet-labs/mgo-go-sdk` is the official Go SDK for Mango Blockchain, providing capabilities to interact with the Mango chain, including account management, transaction building, token operations, and more.

## Features

- **Account Management**: Supports `ed25519` and `secp256k1` key pair generation and signing.
- **Transaction Operations**: Provides transaction creation, signing, and submission.
- **On-Chain Data Queries**: Supports querying token balances, events, objects, and more.
- **WebSocket Subscription**: Enables real-time event push notifications.

## Installation

```sh
# Install using go mod
go get github.com/mangonet-labs/mgo-go-sdk
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
    "github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
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
    "github.com/mangonet-labs/mgo-go-sdk/account/keypair"
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
    "github.com/mangonet-labs/mgo-go-sdk/client"
    "github.com/mangonet-labs/mgo-go-sdk/account/keypair"
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

## Multi-Transaction Blocks

The MGO Go SDK supports **Programmable Transaction Blocks (PTBs)**, allowing you to execute multiple transactions atomically within a single transaction block:

```go
// Create a transaction with multiple operations
tx := transaction.NewTransaction()
tx.SetMgoClient(cli).SetSigner(key).SetSender(sender).SetGasPrice(1000).SetGasBudget(50000000)

// Operation 1: Split coins
splitResult := tx.SplitCoins(tx.Gas(), []transaction.Argument{
    tx.Pure(uint64(1000000000 * 0.01)), // 0.01 MGO
})

// Operation 2: Transfer split coins
tx.TransferObjects([]transaction.Argument{splitResult}, tx.Pure(recipient))

// Execute all operations atomically
resp, err := tx.Execute(ctx, options, "WaitForLocalExecution")
```

### Benefits of Multi-Transaction Blocks:

- **Atomicity**: All operations succeed or fail together
- **Gas Efficiency**: Lower gas costs compared to separate transactions
- **Chaining**: Use outputs from one operation as inputs to another
- **Consistency**: Ensure state consistency across multiple operations

### Multiple Move Calls in One Transaction

You can also make multiple Move calls within a single transaction:

```go
// Multiple Move calls in one transaction
tx := transaction.NewTransaction()
tx.SetMgoClient(cli).SetSigner(key).SetSender(sender).SetGasPrice(1000).SetGasBudget(100000000)

// Move Call 1: Split coins
splitResult := tx.MoveCall(
    "0x0000000000000000000000000000000000000000000000000000000000000002",
    "pay", "split", typeArgs,
    []transaction.Argument{tx.Gas(), tx.Pure(uint64(1000000000 * 0.05))},
)

// Move Call 2: Use result from first call
tx.MoveCall(
    "0x0000000000000000000000000000000000000000000000000000000000000002",
    "pay", "join", typeArgs,
    []transaction.Argument{tx.Gas(), splitResult}, // Chain the calls
)

// Move Call 3: Custom contract interaction
tx.MoveCall(
    "0xpackage_id", "module_name", "function_name", typeArgs,
    []transaction.Argument{splitResult, tx.Pure(params)},
)

// Execute all Move calls atomically
resp, err := tx.Execute(ctx, options, "WaitForLocalExecution")
```

For detailed examples and patterns, see:

- `docs/MULTI_TRANSACTION_GUIDE.md` - Comprehensive guide
- `docs/MULTI_MOVE_CALL_GUIDE.md` - Multi Move call patterns
- `examples/multi_transaction_example.go` - Working examples
- `examples/demo_multi_move_call.go` - Move call demonstrations
- `test/multi_transaction/` - Test cases and validation
- `test/multi_move_call/` - Move call test cases

## Usage Examples

For detailed usage examples, refer to the `test` directory, which includes sample implementations.

## Contributing

If you wish to contribute to the SDK, please follow these steps:

1. Fork this repository.
2. Create a new branch for development.
3. Submit a Pull Request.

## License

[Apache 2.0 license](LICENSE)
