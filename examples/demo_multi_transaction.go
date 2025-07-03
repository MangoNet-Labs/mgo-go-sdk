package main

import (
	"fmt"
	"log"

	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model"
	"github.com/mangonet-labs/mgo-go-sdk/transaction"
)

// DemoMultiTransaction demonstrates the multi-transaction capabilities
// of the MGO Go SDK without actually executing transactions
func main() {
	fmt.Println("🚀 MGO Go SDK - Multi-Transaction Block Demo")
	fmt.Println("============================================")

	// Initialize client (for demo purposes, we won't execute)
	cli := client.NewMgoClient(config.RpcMgoTestnetEndpoint)

	// Create a demo keypair
	key, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		log.Fatal("Failed to create keypair:", err)
	}

	fmt.Printf("📝 Signer Address: %s\n\n", key.MgoAddress())

	// Demo 1: Basic Multi-Transaction
	fmt.Println("📦 Demo 1: Basic Multi-Transaction Block")
	fmt.Println("----------------------------------------")
	demoBasicMultiTransaction(cli, key)

	// Demo 2: Complex Chaining
	fmt.Println("\n📦 Demo 2: Complex Transaction Chaining")
	fmt.Println("---------------------------------------")
	demoComplexChaining(cli, key)

	// Demo 3: MoveCall Integration
	fmt.Println("\n📦 Demo 3: MoveCall with Coin Operations")
	fmt.Println("----------------------------------------")
	demoMoveCallIntegration(cli, key)

	fmt.Println("\n✅ Demo completed! The MGO Go SDK already supports multiple transactions in one block.")
	fmt.Println("📚 See docs/MULTI_TRANSACTION_GUIDE.md for detailed documentation.")
}

func demoBasicMultiTransaction(cli *client.Client, key *keypair.Keypair) {
	// Create a new transaction block
	tx := transaction.NewTransaction()

	// Set up transaction parameters
	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"

	fmt.Println("Building transaction with multiple commands:")

	// Command 1: Split coins
	fmt.Println("  1. SplitCoins: Split 0.01 MGO from gas coin")
	splitResult := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.01)), // 0.01 MGO
	})

	// Command 2: Transfer split coins
	fmt.Println("  2. TransferObjects: Transfer split coins to recipient")
	tx.TransferObjects([]transaction.Argument{splitResult}, tx.Pure(recipient))

	// Get transaction data to show structure
	txData, err := tx.GetTransactionData()
	if err != nil {
		log.Printf("Error getting transaction data: %v", err)
		return
	}

	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	fmt.Printf("\n✅ Transaction block created with %d commands\n", len(commands))
	
	for i, cmd := range commands {
		if cmd.SplitCoins != nil {
			fmt.Printf("   Command %d: SplitCoins\n", i+1)
		} else if cmd.TransferObjects != nil {
			fmt.Printf("   Command %d: TransferObjects\n", i+1)
		} else if cmd.MoveCall != nil {
			fmt.Printf("   Command %d: MoveCall\n", i+1)
		} else if cmd.MergeCoins != nil {
			fmt.Printf("   Command %d: MergeCoins\n", i+1)
		}
	}
}

func demoComplexChaining(cli *client.Client, key *keypair.Keypair) {
	tx := transaction.NewTransaction()
	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"

	fmt.Println("Building complex chained transaction:")

	// Command 1: Split gas coin
	fmt.Println("  1. SplitCoins: Split 0.05 MGO from gas coin")
	split1 := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.05)), // 0.05 MGO
	})

	// Command 2: Split the result again
	fmt.Println("  2. SplitCoins: Split 0.01 MGO from the previous result")
	split2 := tx.SplitCoins(split1, []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.01)), // 0.01 MGO from the 0.05 MGO
	})

	// Command 3: Split gas again for merge operation
	fmt.Println("  3. SplitCoins: Split another 0.02 MGO from gas coin")
	split3 := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.02)), // 0.02 MGO
	})

	// Command 4: Merge coins
	fmt.Println("  4. MergeCoins: Merge split results together")
	merged := tx.MergeCoins(split2, []transaction.Argument{split3})

	// Command 5: Transfer the final result
	fmt.Println("  5. TransferObjects: Transfer merged coins to recipient")
	tx.TransferObjects([]transaction.Argument{merged}, tx.Pure(recipient))

	// Show transaction structure
	txData, err := tx.GetTransactionData()
	if err != nil {
		log.Printf("Error getting transaction data: %v", err)
		return
	}

	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	fmt.Printf("\n✅ Complex transaction block created with %d commands\n", len(commands))
	fmt.Println("   This demonstrates chaining where outputs from one command")
	fmt.Println("   become inputs to subsequent commands.")
}

func demoMoveCallIntegration(cli *client.Client, key *keypair.Keypair) {
	tx := transaction.NewTransaction()
	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"

	fmt.Println("Building transaction with MoveCall integration:")

	// Command 1: Split coins for use in MoveCall
	fmt.Println("  1. SplitCoins: Split 0.1 MGO for MoveCall")
	splitCoin := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.1)), // 0.1 MGO
	})

	// Command 2: MoveCall using the split coin
	fmt.Println("  2. MoveCall: Call pay::split function with split coin")
	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		log.Printf("Error converting address: %v", err)
		return
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

	// Command 3: Transfer the result
	fmt.Println("  3. TransferObjects: Transfer MoveCall result to recipient")
	tx.TransferObjects([]transaction.Argument{moveCallResult}, tx.Pure(recipient))

	// Show transaction structure
	txData, err := tx.GetTransactionData()
	if err != nil {
		log.Printf("Error getting transaction data: %v", err)
		return
	}

	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	fmt.Printf("\n✅ MoveCall transaction block created with %d commands\n", len(commands))
	fmt.Println("   This demonstrates integrating Move function calls")
	fmt.Println("   with coin operations in a single atomic transaction.")
}
