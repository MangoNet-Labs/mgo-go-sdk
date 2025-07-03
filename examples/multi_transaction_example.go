package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/transaction"
)

// MultiTransactionExample demonstrates how to create and execute multiple transactions
// within a single transaction block using the MGO Go SDK
func main() {
	// Initialize client
	cli := client.NewMgoClient(config.RpcMgoTestnetEndpoint)
	ctx := context.Background()

	// Create a keypair for signing transactions
	key, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "your_private_key_here")
	if err != nil {
		log.Fatal("Failed to create keypair:", err)
	}

	// Example 1: Simple Multi-Transaction Block
	fmt.Println("=== Example 1: Split Coins and Transfer ===")
	simpleMultiTransaction(ctx, cli, key)

	// Example 2: Complex Multi-Transaction Block
	fmt.Println("\n=== Example 2: Complex Multi-Transaction Block ===")
	complexMultiTransaction(ctx, cli, key)

	// Example 3: MoveCall with Other Operations
	fmt.Println("\n=== Example 3: MoveCall with Coin Operations ===")
	moveCallWithCoinOps(ctx, cli, key)
}

// simpleMultiTransaction demonstrates a basic multi-transaction block
// that splits coins and transfers them to multiple recipients
func simpleMultiTransaction(ctx context.Context, cli *client.Client, key *keypair.Keypair) {
	// Create a new transaction
	tx := transaction.NewTransaction()

	// Set up transaction parameters
	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	// Recipients for the transfers
	recipient1 := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"
	recipient2 := "0xbb3888e6c078a8ccedde58394873584ba39878984f1f8da4cba870de7eb5c3d2"

	// Transaction 1: Split coins into multiple amounts
	splitResult := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.01)), // 0.01 MGO
		tx.Pure(uint64(1000000000 * 0.02)), // 0.02 MGO
	})

	// Transaction 2: Transfer first split coin to recipient1
	tx.TransferObjects([]transaction.Argument{splitResult}, tx.Pure(recipient1))

	// Transaction 3: Get the second split result and transfer to recipient2
	// Note: splitResult refers to the first result, we need to get the second one
	secondSplit := transaction.Argument{Result: &[]uint16{0}[0]} // First command result, second output
	tx.TransferObjects([]transaction.Argument{secondSplit}, tx.Pure(recipient2))

	// Execute the multi-transaction block
	executeTransaction(ctx, tx, "Simple Multi-Transaction")
}

// complexMultiTransaction demonstrates a more complex multi-transaction block
// with multiple coin operations
func complexMultiTransaction(ctx context.Context, cli *client.Client, key *keypair.Keypair) {
	tx := transaction.NewTransaction()

	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"

	// Transaction 1: Split gas coin into multiple parts
	splitResult1 := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.05)), // 0.05 MGO
		tx.Pure(uint64(1000000000 * 0.03)), // 0.03 MGO
	})

	// Transaction 2: Split one of the results again
	splitResult2 := tx.SplitCoins(splitResult1, []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.01)), // 0.01 MGO from the 0.05 MGO
	})

	// Transaction 3: Merge some coins back together
	// Get the second result from first split
	secondCoin := transaction.Argument{Result: &[]uint16{0}[0]} // Adjust index as needed
	mergedCoin := tx.MergeCoins(splitResult2, []transaction.Argument{secondCoin})

	// Transaction 4: Transfer the merged coin
	tx.TransferObjects([]transaction.Argument{mergedCoin}, tx.Pure(recipient))

	executeTransaction(ctx, tx, "Complex Multi-Transaction")
}

// moveCallWithCoinOps demonstrates combining MoveCall with coin operations
func moveCallWithCoinOps(ctx context.Context, cli *client.Client, key *keypair.Keypair) {
	tx := transaction.NewTransaction()

	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	// Transaction 1: Split coins for use in MoveCall
	splitCoin := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.1)), // 0.1 MGO
	})

	// Transaction 2: Make a MoveCall using the split coin
	// This example calls a hypothetical function that requires a coin input
	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		log.Fatal("Failed to convert address:", err)
	}

	moveCallResult := tx.MoveCall(
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"pay",
		"split",
		[]transaction.TypeTag{
			{
				Struct: &transaction.StructTag{
					Address: *addressBytes,
					Module:  "mgo",
					Name:    "MGO",
				},
			},
		},
		[]transaction.Argument{
			splitCoin,
			tx.Pure(uint64(1000000000 * 0.05)), // Split the 0.1 MGO into 0.05 MGO
		},
	)

	// Transaction 3: Transfer the result of the MoveCall
	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"
	tx.TransferObjects([]transaction.Argument{moveCallResult}, tx.Pure(recipient))

	executeTransaction(ctx, tx, "MoveCall with Coin Operations")
}

// executeTransaction is a helper function to execute and log transaction results
func executeTransaction(ctx context.Context, tx *transaction.Transaction, description string) {
	fmt.Printf("Executing: %s\n", description)

	// Build the transaction to see the structure
	txData, err := tx.GetTransactionData()
	if err != nil {
		log.Printf("Failed to get transaction data: %v", err)
		return
	}

	// Print number of commands in this transaction block
	commandCount := len(txData.V1.Kind.ProgrammableTransaction.Commands)
	fmt.Printf("Transaction block contains %d commands\n", commandCount)

	// Execute the transaction
	resp, err := tx.Execute(
		ctx,
		request.MgoTransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
			ShowEvents:   true,
		},
		"WaitForLocalExecution",
	)
	if err != nil {
		log.Printf("Failed to execute transaction: %v", err)
		return
	}

	fmt.Printf("Transaction executed successfully!\n")
	fmt.Printf("Digest: %s\n", resp.Digest)
	fmt.Printf("Status: %s\n", resp.Effects.Status.Status)
	
	if resp.Effects.GasUsed != nil {
		fmt.Printf("Gas used: %s\n", resp.Effects.GasUsed.ComputationCost)
	}
	
	fmt.Println("---")
}
