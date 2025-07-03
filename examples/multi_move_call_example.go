package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model"
	"github.com/mangonet-labs/mgo-go-sdk/transaction"
)

// MultiMoveCallExample demonstrates how to make multiple Move calls
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

	fmt.Println("🚀 Multi Move Call Examples")
	fmt.Println("===========================")

	// Example 1: Multiple Independent Move Calls
	fmt.Println("\n📦 Example 1: Multiple Independent Move Calls")
	multipleIndependentMoveCalls(ctx, cli, key)

	// Example 2: Chained Move Calls
	fmt.Println("\n📦 Example 2: Chained Move Calls")
	chainedMoveCalls(ctx, cli, key)

	// Example 3: Move Calls with Coin Operations
	fmt.Println("\n📦 Example 3: Move Calls with Coin Operations")
	moveCallsWithCoinOps(ctx, cli, key)

	// Example 4: Complex Multi-Move Call Transaction
	fmt.Println("\n📦 Example 4: Complex Multi-Move Call Transaction")
	complexMultiMoveCall(ctx, cli, key)
}

// multipleIndependentMoveCalls demonstrates multiple Move calls that don't depend on each other
func multipleIndependentMoveCalls(ctx context.Context, cli *client.Client, key *keypair.Keypair) {
	tx := transaction.NewTransaction()

	// Set up transaction parameters
	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(100000000). // Higher gas budget for multiple Move calls
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	// Get address bytes for type arguments
	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		log.Fatal("Failed to convert address:", err)
	}

	// Move Call 1: Transfer MGO
	fmt.Println("  Adding Move Call 1: MGO Transfer")
	tx.MoveCall(
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"mgo",
		"transfer",
		[]transaction.TypeTag{}, // No type arguments for this call
		[]transaction.Argument{
			tx.Object("0xe586e913e413de2df45e3eca0f0adf342a1f5d8d71e61805e14fd4872529f727"), // Object to transfer
			tx.Pure("0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"),   // Recipient
		},
	)

	// Move Call 2: Split coins
	fmt.Println("  Adding Move Call 2: Coin Split")
	tx.MoveCall(
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
			tx.Gas(),
			tx.Pure(uint64(1000000000 * 0.01)), // 0.01 MGO
		},
	)

	// Move Call 3: Another operation (example with different package)
	fmt.Println("  Adding Move Call 3: Custom Package Call")
	tx.MoveCall(
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", // Example package ID
		"example_module",
		"example_function",
		[]transaction.TypeTag{}, // No type arguments
		[]transaction.Argument{
			tx.Pure(uint64(42)), // Some parameter
			tx.Pure("0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"), // Address parameter
		},
	)

	// Show transaction structure without executing
	showTransactionStructure(tx, "Multiple Independent Move Calls")
}

// chainedMoveCalls demonstrates Move calls where outputs from one call are used in another
func chainedMoveCalls(ctx context.Context, cli *client.Client, key *keypair.Keypair) {
	tx := transaction.NewTransaction()

	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(100000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		log.Fatal("Failed to convert address:", err)
	}

	// Move Call 1: Split coins to get a coin object
	fmt.Println("  Move Call 1: Split coins")
	splitResult := tx.MoveCall(
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
			tx.Gas(),
			tx.Pure(uint64(1000000000 * 0.05)), // 0.05 MGO
		},
	)

	// Move Call 2: Use the result from the first call
	fmt.Println("  Move Call 2: Use split result in another operation")
	tx.MoveCall(
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"pay",
		"join",
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
			tx.Gas(),    // Destination coin
			splitResult, // Source coin from previous call
		},
	)

	// Move Call 3: Another operation using a different result
	fmt.Println("  Move Call 3: Custom function using previous results")
	tx.MoveCall(
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"defi_module",
		"deposit",
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
			tx.Object("0xpool_object_id"), // Pool object
			splitResult,                   // Coin to deposit (reusing from call 1)
		},
	)

	showTransactionStructure(tx, "Chained Move Calls")
}

// moveCallsWithCoinOps demonstrates mixing Move calls with coin operations
func moveCallsWithCoinOps(ctx context.Context, cli *client.Client, key *keypair.Keypair) {
	tx := transaction.NewTransaction()

	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(100000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"
	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		log.Fatal("Failed to convert address:", err)
	}

	// Operation 1: Split coins using built-in operation
	fmt.Println("  Operation 1: SplitCoins (built-in)")
	splitCoin1 := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.1)), // 0.1 MGO
	})

	// Operation 2: Move call to split more coins
	fmt.Println("  Operation 2: MoveCall split")
	splitCoin2 := tx.MoveCall(
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
			tx.Gas(),
			tx.Pure(uint64(1000000000 * 0.05)), // 0.05 MGO
		},
	)

	// Operation 3: Move call to process the coins
	fmt.Println("  Operation 3: MoveCall to process coins")
	processedCoin := tx.MoveCall(
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"processor",
		"process_coins",
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
			splitCoin1,
			splitCoin2,
			tx.Pure(uint64(100)), // Processing parameter
		},
	)

	// Operation 4: Built-in merge operation
	fmt.Println("  Operation 4: MergeCoins (built-in)")
	mergedCoin := tx.MergeCoins(processedCoin, []transaction.Argument{splitCoin1})

	// Operation 5: Final Move call
	fmt.Println("  Operation 5: Final MoveCall")
	tx.MoveCall(
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"mgo",
		"transfer",
		[]transaction.TypeTag{},
		[]transaction.Argument{
			mergedCoin,
			tx.Pure(recipient),
		},
	)

	showTransactionStructure(tx, "Move Calls with Coin Operations")
}

// complexMultiMoveCall demonstrates a complex scenario with many Move calls
func complexMultiMoveCall(ctx context.Context, cli *client.Client, key *keypair.Keypair) {
	tx := transaction.NewTransaction()

	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(200000000). // Even higher gas budget for complex operations
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		log.Fatal("Failed to convert address:", err)
	}

	fmt.Println("  Building complex multi-Move call transaction...")

	// Step 1: Initialize resources
	fmt.Println("    Step 1: Initialize resources")
	resource1 := tx.MoveCall(
		"0xpackage1",
		"factory",
		"create_resource",
		[]transaction.TypeTag{},
		[]transaction.Argument{
			tx.Pure(uint64(1000)),
			tx.Pure("resource_type_1"),
		},
	)

	// Step 2: Create another resource
	fmt.Println("    Step 2: Create another resource")
	resource2 := tx.MoveCall(
		"0xpackage2",
		"factory",
		"create_resource",
		[]transaction.TypeTag{},
		[]transaction.Argument{
			tx.Pure(uint64(2000)),
			tx.Pure("resource_type_2"),
		},
	)

	// Step 3: Process resources together
	fmt.Println("    Step 3: Process resources together")
	processedResult := tx.MoveCall(
		"0xpackage3",
		"processor",
		"combine_resources",
		[]transaction.TypeTag{},
		[]transaction.Argument{
			resource1,
			resource2,
			tx.Pure(uint64(42)), // Processing parameter
		},
	)

	// Step 4: Split coins for fees
	fmt.Println("    Step 4: Prepare coins for fees")
	feeCoin := tx.MoveCall(
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
			tx.Gas(),
			tx.Pure(uint64(1000000000 * 0.01)), // 0.01 MGO for fees
		},
	)

	// Step 5: Execute main operation with fee
	fmt.Println("    Step 5: Execute main operation")
	finalResult := tx.MoveCall(
		"0xpackage4",
		"main_contract",
		"execute_with_fee",
		[]transaction.TypeTag{},
		[]transaction.Argument{
			processedResult,
			feeCoin,
			tx.Pure("0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"), // Beneficiary
		},
	)

	// Step 6: Finalize and cleanup
	fmt.Println("    Step 6: Finalize and cleanup")
	tx.MoveCall(
		"0xpackage5",
		"cleanup",
		"finalize",
		[]transaction.TypeTag{},
		[]transaction.Argument{
			finalResult,
			tx.Pure(uint64(1)), // Cleanup mode
		},
	)

	showTransactionStructure(tx, "Complex Multi-Move Call Transaction")
}

// Helper function to show transaction structure
func showTransactionStructure(tx *transaction.Transaction, description string) {
	txData, err := tx.GetTransactionData()
	if err != nil {
		log.Printf("Error getting transaction data: %v", err)
		return
	}

	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	fmt.Printf("\n✅ %s\n", description)
	fmt.Printf("   Transaction block contains %d commands:\n", len(commands))

	moveCallCount := 0
	for i, cmd := range commands {
		if cmd.MoveCall != nil {
			moveCallCount++
			packageStr := fmt.Sprintf("%x", cmd.MoveCall.Package)
			if len(packageStr) > 10 {
				packageStr = packageStr[:10] + "..."
			}
			fmt.Printf("   Command %d: MoveCall (%s::%s::%s)\n",
				i+1,
				packageStr,
				cmd.MoveCall.Module,
				cmd.MoveCall.Function)
		} else if cmd.SplitCoins != nil {
			fmt.Printf("   Command %d: SplitCoins\n", i+1)
		} else if cmd.MergeCoins != nil {
			fmt.Printf("   Command %d: MergeCoins\n", i+1)
		} else if cmd.TransferObjects != nil {
			fmt.Printf("   Command %d: TransferObjects\n", i+1)
		} else {
			fmt.Printf("   Command %d: Other operation\n", i+1)
		}
	}

	fmt.Printf("   📊 Total Move calls: %d\n", moveCallCount)
	fmt.Println("   🔗 All operations will execute atomically")
}
