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

// DemoMultiMoveCall demonstrates multiple Move calls in one transaction
func main() {
	fmt.Println("🚀 MGO Go SDK - Multi Move Call Demo")
	fmt.Println("====================================")

	// Initialize client (for demo purposes, we won't execute)
	cli := client.NewMgoClient(config.RpcMgoTestnetEndpoint)

	// Create a demo keypair
	key, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		log.Fatal("Failed to create keypair:", err)
	}

	fmt.Printf("📝 Signer Address: %s\n\n", key.MgoAddress())

	// Demo 1: Multiple Independent Move Calls
	fmt.Println("📦 Demo 1: Multiple Independent Move Calls")
	fmt.Println("------------------------------------------")
	demoMultipleIndependentMoveCalls(cli, key)

	// Demo 2: Chained Move Calls
	fmt.Println("\n📦 Demo 2: Chained Move Calls")
	fmt.Println("-----------------------------")
	demoChainedMoveCalls(cli, key)

	// Demo 3: Move Calls with Coin Operations
	fmt.Println("\n📦 Demo 3: Move Calls with Coin Operations")
	fmt.Println("------------------------------------------")
	demoMoveCallsWithCoinOps(cli, key)

	fmt.Println("\n✅ Demo completed! Multiple Move calls in one transaction are fully supported.")
	fmt.Println("📚 See docs/MULTI_MOVE_CALL_GUIDE.md for detailed documentation.")
}

func demoMultipleIndependentMoveCalls(cli *client.Client, key *keypair.Keypair) {
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

	fmt.Println("Building transaction with multiple independent Move calls:")

	// Move Call 1: Split coins
	fmt.Println("  1. MoveCall: pay::split (Split 0.01 MGO)")
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

	// Move Call 2: Another split operation
	fmt.Println("  2. MoveCall: pay::split (Split 0.02 MGO)")
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
			tx.Pure(uint64(1000000000 * 0.02)), // 0.02 MGO
		},
	)

	// Move Call 3: Example custom package call
	fmt.Println("  3. MoveCall: custom_package::example_function")
	tx.MoveCall(
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"example_module",
		"example_function",
		[]transaction.TypeTag{}, // No type arguments
		[]transaction.Argument{
			tx.Pure(uint64(42)), // Some parameter
			tx.Pure("0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"), // Address parameter
		},
	)

	showTransactionStructure(tx, "Multiple Independent Move Calls")
}

func demoChainedMoveCalls(cli *client.Client, key *keypair.Keypair) {
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

	fmt.Println("Building transaction with chained Move calls:")

	// Move Call 1: Split coins to get a coin object
	fmt.Println("  1. MoveCall: pay::split (Create coin for chaining)")
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
	fmt.Println("  2. MoveCall: pay::join (Use result from call 1)")
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

	// Move Call 3: Another operation using the split result
	fmt.Println("  3. MoveCall: defi::deposit (Use split result in DeFi)")
	tx.MoveCall(
		"0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
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
			tx.Object("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"), // Pool object
			splitResult, // Coin to deposit (reusing from call 1)
		},
	)

	showTransactionStructure(tx, "Chained Move Calls")
}

func demoMoveCallsWithCoinOps(cli *client.Client, key *keypair.Keypair) {
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

	fmt.Println("Building transaction mixing Move calls with coin operations:")

	// Operation 1: Split coins using built-in operation
	fmt.Println("  1. SplitCoins: Built-in split operation (0.1 MGO)")
	splitCoin1 := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.1)), // 0.1 MGO
	})

	// Operation 2: Move call to split more coins
	fmt.Println("  2. MoveCall: pay::split (0.05 MGO)")
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
	fmt.Println("  3. MoveCall: processor::process_coins")
	processedCoin := tx.MoveCall(
		"0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
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
	fmt.Println("  4. MergeCoins: Built-in merge operation")
	mergedCoin := tx.MergeCoins(processedCoin, []transaction.Argument{splitCoin1})

	// Operation 5: Final Move call
	fmt.Println("  5. MoveCall: mgo::transfer (Final transfer)")
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
