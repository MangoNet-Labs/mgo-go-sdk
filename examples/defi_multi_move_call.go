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

// DeFiMultiMoveCallExample demonstrates a practical DeFi use case
// with multiple Move calls in one transaction
func main() {
	fmt.Println("🏦 DeFi Multi Move Call Example")
	fmt.Println("===============================")
	fmt.Println("This example shows how to perform complex DeFi operations")
	fmt.Println("atomically using multiple Move calls in one transaction.\n")

	// Initialize client
	cli := client.NewMgoClient(config.RpcMgoTestnetEndpoint)
	ctx := context.Background()

	// Create keypair (replace with your private key)
	key, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		log.Fatal("Failed to create keypair:", err)
	}

	fmt.Printf("📝 Trader Address: %s\n\n", key.MgoAddress())

	// Example 1: Liquidity Pool Operations
	fmt.Println("💧 Example 1: Add Liquidity to Multiple Pools")
	addLiquidityToMultiplePools(ctx, cli, key)

	// Example 2: Token Swap Chain
	fmt.Println("\n🔄 Example 2: Multi-Hop Token Swap")
	multiHopTokenSwap(ctx, cli, key)

	// Example 3: Yield Farming Strategy
	fmt.Println("\n🌾 Example 3: Yield Farming Strategy")
	yieldFarmingStrategy(ctx, cli, key)

	fmt.Println("\n✅ All DeFi operations can be executed atomically!")
	fmt.Println("💡 This ensures consistency and reduces gas costs.")
}

// addLiquidityToMultiplePools demonstrates adding liquidity to multiple pools atomically
func addLiquidityToMultiplePools(ctx context.Context, cli *client.Client, key *keypair.Keypair) {
	tx := transaction.NewTransaction()

	// Set up transaction
	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(200000000). // Higher gas for complex DeFi operations
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	// Get type arguments for MGO
	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		log.Fatal("Failed to convert address:", err)
	}

	mgoType := transaction.TypeTag{
		Struct: &transaction.StructTag{
			Address: *addressBytes,
			Module:  "mgo",
			Name:    "MGO",
		},
	}

	fmt.Println("  Building liquidity operations:")

	// Step 1: Split MGO for multiple pools
	fmt.Println("    1. Split MGO for Pool A (0.5 MGO)")
	mgoForPoolA := tx.MoveCall(
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"pay", "split",
		[]transaction.TypeTag{mgoType},
		[]transaction.Argument{
			tx.Gas(),
			tx.Pure(uint64(500000000)), // 0.5 MGO
		},
	)

	// Step 2: Split MGO for Pool B
	fmt.Println("    2. Split MGO for Pool B (0.3 MGO)")
	mgoForPoolB := tx.MoveCall(
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"pay", "split",
		[]transaction.TypeTag{mgoType},
		[]transaction.Argument{
			tx.Gas(),
			tx.Pure(uint64(300000000)), // 0.3 MGO
		},
	)

	// Step 3: Add liquidity to Pool A (MGO/USDC)
	fmt.Println("    3. Add liquidity to MGO/USDC pool")
	lpTokenA := tx.MoveCall(
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"liquidity_pool",
		"add_liquidity",
		[]transaction.TypeTag{mgoType, {
			Struct: &transaction.StructTag{
				Address: *addressBytes,
				Module:  "usdc",
				Name:    "USDC",
			},
		}},
		[]transaction.Argument{
			mgoForPoolA,
			tx.Object("0x1111111111111111111111111111111111111111111111111111111111111111"), // USDC coin
			tx.Object("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), // Pool A object
		},
	)

	// Step 4: Add liquidity to Pool B (MGO/ETH)
	fmt.Println("    4. Add liquidity to MGO/ETH pool")
	lpTokenB := tx.MoveCall(
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"liquidity_pool",
		"add_liquidity",
		[]transaction.TypeTag{mgoType, {
			Struct: &transaction.StructTag{
				Address: *addressBytes,
				Module:  "eth",
				Name:    "ETH",
			},
		}},
		[]transaction.Argument{
			mgoForPoolB,
			tx.Object("0x2222222222222222222222222222222222222222222222222222222222222222"), // ETH coin
			tx.Object("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), // Pool B object
		},
	)

	// Step 5: Stake LP tokens for rewards
	fmt.Println("    5. Stake LP tokens in farming contract")
	tx.MoveCall(
		"0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		"farming",
		"stake_lp_tokens",
		[]transaction.TypeTag{},
		[]transaction.Argument{
			lpTokenA,
			lpTokenB,
			tx.Object("0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"), // Farming pool
		},
	)

	showTransactionSummary(tx, "Add Liquidity to Multiple Pools")
}

// multiHopTokenSwap demonstrates a multi-hop token swap
func multiHopTokenSwap(ctx context.Context, cli *client.Client, key *keypair.Keypair) {
	tx := transaction.NewTransaction()

	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(200000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		log.Fatal("Failed to convert address:", err)
	}

	fmt.Println("  Building multi-hop swap: MGO → USDC → ETH → BTC:")

	// Step 1: Prepare MGO for swapping
	fmt.Println("    1. Prepare 1.0 MGO for swapping")
	swapAmount := tx.MoveCall(
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"pay", "split",
		[]transaction.TypeTag{{
			Struct: &transaction.StructTag{
				Address: *addressBytes,
				Module:  "mgo",
				Name:    "MGO",
			},
		}},
		[]transaction.Argument{
			tx.Gas(),
			tx.Pure(uint64(1000000000)), // 1.0 MGO
		},
	)

	// Step 2: Swap MGO → USDC
	fmt.Println("    2. Swap MGO → USDC")
	usdcAmount := tx.MoveCall(
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"swap",
		"swap_exact_input",
		[]transaction.TypeTag{
			{Struct: &transaction.StructTag{Address: *addressBytes, Module: "mgo", Name: "MGO"}},
			{Struct: &transaction.StructTag{Address: *addressBytes, Module: "usdc", Name: "USDC"}},
		},
		[]transaction.Argument{
			swapAmount,
			tx.Pure(uint64(0)), // Min output amount
			tx.Object("0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"), // MGO/USDC pool
		},
	)

	// Step 3: Swap USDC → ETH
	fmt.Println("    3. Swap USDC → ETH")
	ethAmount := tx.MoveCall(
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"swap",
		"swap_exact_input",
		[]transaction.TypeTag{
			{Struct: &transaction.StructTag{Address: *addressBytes, Module: "usdc", Name: "USDC"}},
			{Struct: &transaction.StructTag{Address: *addressBytes, Module: "eth", Name: "ETH"}},
		},
		[]transaction.Argument{
			usdcAmount,
			tx.Pure(uint64(0)), // Min output amount
			tx.Object("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"), // USDC/ETH pool
		},
	)

	// Step 4: Swap ETH → BTC
	fmt.Println("    4. Swap ETH → BTC")
	btcAmount := tx.MoveCall(
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"swap",
		"swap_exact_input",
		[]transaction.TypeTag{
			{Struct: &transaction.StructTag{Address: *addressBytes, Module: "eth", Name: "ETH"}},
			{Struct: &transaction.StructTag{Address: *addressBytes, Module: "btc", Name: "BTC"}},
		},
		[]transaction.Argument{
			ethAmount,
			tx.Pure(uint64(0)), // Min output amount
			tx.Object("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"), // ETH/BTC pool
		},
	)

	// Step 5: Transfer final BTC to user
	fmt.Println("    5. Transfer BTC to user")
	tx.TransferObjects(
		[]transaction.Argument{btcAmount},
		tx.Pure(key.MgoAddress()),
	)

	showTransactionSummary(tx, "Multi-Hop Token Swap")
}

// yieldFarmingStrategy demonstrates a complex yield farming strategy
func yieldFarmingStrategy(ctx context.Context, cli *client.Client, key *keypair.Keypair) {
	tx := transaction.NewTransaction()

	tx.SetMgoClient(cli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(300000000). // Even higher gas for complex strategy
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	fmt.Println("  Building yield farming strategy:")

	// Step 1: Claim existing rewards
	fmt.Println("    1. Claim pending rewards from existing positions")
	rewards := tx.MoveCall(
		"0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		"farming",
		"claim_rewards",
		[]transaction.TypeTag{},
		[]transaction.Argument{
			tx.Object("0x3333333333333333333333333333333333333333333333333333333333333333"), // User's farming position
		},
	)

	// Step 2: Compound rewards by adding to liquidity
	fmt.Println("    2. Compound rewards into new liquidity position")
	newLpTokens := tx.MoveCall(
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"liquidity_pool",
		"add_single_asset_liquidity",
		[]transaction.TypeTag{},
		[]transaction.Argument{
			rewards,
			tx.Object("0x4444444444444444444444444444444444444444444444444444444444444444"), // Pool object
		},
	)

	// Step 3: Stake new LP tokens
	fmt.Println("    3. Stake new LP tokens for additional rewards")
	tx.MoveCall(
		"0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		"farming",
		"stake_lp_tokens",
		[]transaction.TypeTag{},
		[]transaction.Argument{
			newLpTokens,
			tx.Object("0x5555555555555555555555555555555555555555555555555555555555555555"), // Farming pool
		},
	)

	// Step 4: Update user's strategy parameters
	fmt.Println("    4. Update strategy parameters")
	tx.MoveCall(
		"0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
		"strategy",
		"update_parameters",
		[]transaction.TypeTag{},
		[]transaction.Argument{
			tx.Object("0x6666666666666666666666666666666666666666666666666666666666666666"), // Strategy object
			tx.Pure(uint64(80)), // New allocation percentage
			tx.Pure(uint64(24)), // Rebalance frequency (hours)
		},
	)

	showTransactionSummary(tx, "Yield Farming Strategy")
}

// Helper function to show transaction summary
func showTransactionSummary(tx *transaction.Transaction, description string) {
	txData, err := tx.GetTransactionData()
	if err != nil {
		log.Printf("Error getting transaction data: %v", err)
		return
	}

	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	moveCallCount := 0
	otherOpCount := 0

	for _, cmd := range commands {
		if cmd.MoveCall != nil {
			moveCallCount++
		} else {
			otherOpCount++
		}
	}

	fmt.Printf("\n  ✅ %s Transaction Summary:\n", description)
	fmt.Printf("     📊 Total commands: %d\n", len(commands))
	fmt.Printf("     🔧 Move calls: %d\n", moveCallCount)
	fmt.Printf("     ⚙️  Other operations: %d\n", otherOpCount)
	fmt.Printf("     💰 All operations execute atomically\n")
	fmt.Printf("     ⛽ Single gas payment for entire transaction\n")

	// Uncomment to actually execute:
	/*
		fmt.Printf("     🚀 Executing transaction...\n")
		resp, err := tx.Execute(
			context.Background(),
			request.MgoTransactionBlockOptions{
				ShowInput:    true,
				ShowEffects:  true,
				ShowEvents:   true,
			},
			"WaitForLocalExecution",
		)
		if err != nil {
			log.Printf("     ❌ Transaction failed: %v\n", err)
			return
		}
		fmt.Printf("     ✅ Transaction successful! Digest: %s\n", resp.Digest)
	*/
}
