package multi_move_call

import (
	"context"
	"testing"

	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model"
	"github.com/mangonet-labs/mgo-go-sdk/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	devCli = client.NewMgoClient(config.RpcMgoTestnetEndpoint)
	ctx    = context.Background()
)

// getSigner creates a test keypair for signing transactions
func getSigner() (*keypair.Keypair, error) {
	return keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
}

// TestMultipleMoveCallsBasic tests basic multiple Move calls in one transaction
func TestMultipleMoveCallsBasic(t *testing.T) {
	key, err := getSigner()
	require.NoError(t, err)

	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(100000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	require.NoError(t, err)

	// Move Call 1: Split coins
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
			tx.Pure(uint64(1000000000 * 0.01)),
		},
	)

	// Move Call 2: Another split operation
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
			tx.Pure(uint64(1000000000 * 0.02)),
		},
	)

	// Verify transaction structure
	txData, err := tx.GetTransactionData()
	require.NoError(t, err)

	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	assert.Len(t, commands, 2)

	// Both commands should be Move calls
	assert.NotNil(t, commands[0].MoveCall)
	assert.NotNil(t, commands[1].MoveCall)

	// Verify Move call details
	assert.Equal(t, "pay", commands[0].MoveCall.Module)
	assert.Equal(t, "split", commands[0].MoveCall.Function)
	assert.Equal(t, "pay", commands[1].MoveCall.Module)
	assert.Equal(t, "split", commands[1].MoveCall.Function)

	t.Log("Multiple Move calls structure validated successfully")
}

// TestChainedMoveCalls tests Move calls that use results from previous calls
func TestChainedMoveCalls(t *testing.T) {
	key, err := getSigner()
	require.NoError(t, err)

	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(100000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	require.NoError(t, err)

	// Move Call 1: Split coins
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
			tx.Pure(uint64(1000000000 * 0.05)),
		},
	)

	// Move Call 2: Use result from first call
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
			tx.Gas(),
			splitResult, // Using result from previous Move call
		},
	)

	// Verify transaction structure
	txData, err := tx.GetTransactionData()
	require.NoError(t, err)

	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	assert.Len(t, commands, 2)

	// Both should be Move calls
	assert.NotNil(t, commands[0].MoveCall)
	assert.NotNil(t, commands[1].MoveCall)

	// Verify the second call uses result from first
	secondCallArgs := commands[1].MoveCall.Arguments
	assert.Len(t, secondCallArgs, 2)
	// The second argument should be a Result reference
	assert.NotNil(t, secondCallArgs[1].Result)

	t.Log("Chained Move calls validated successfully")
}

// TestMoveCallsWithCoinOperations tests mixing Move calls with built-in coin operations
func TestMoveCallsWithCoinOperations(t *testing.T) {
	key, err := getSigner()
	require.NoError(t, err)

	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(100000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	require.NoError(t, err)

	// Operation 1: Built-in SplitCoins
	splitCoin := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.1)),
	})

	// Operation 2: Move call using the split coin
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
			tx.Pure(uint64(1000000000 * 0.01)),
		},
	)

	// Operation 3: Built-in TransferObjects using Move call result
	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"
	tx.TransferObjects([]transaction.Argument{moveCallResult}, tx.Pure(recipient))

	// Verify transaction structure
	txData, err := tx.GetTransactionData()
	require.NoError(t, err)

	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	assert.Len(t, commands, 3)

	// Verify command types
	assert.NotNil(t, commands[0].SplitCoins)      // Built-in operation
	assert.NotNil(t, commands[1].MoveCall)       // Move call
	assert.NotNil(t, commands[2].TransferObjects) // Built-in operation

	t.Log("Move calls with coin operations validated successfully")
}

// TestComplexMultiMoveCall tests a complex scenario with many Move calls
func TestComplexMultiMoveCall(t *testing.T) {
	key, err := getSigner()
	require.NoError(t, err)

	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(200000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	require.NoError(t, err)

	// Add 5 different Move calls
	for i := 0; i < 5; i++ {
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
				tx.Pure(uint64(1000000000 * 0.001 * float64(i+1))), // Different amounts
			},
		)
	}

	// Verify transaction structure
	txData, err := tx.GetTransactionData()
	require.NoError(t, err)

	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	assert.Len(t, commands, 5)

	// All should be Move calls
	moveCallCount := 0
	for _, cmd := range commands {
		if cmd.MoveCall != nil {
			moveCallCount++
		}
	}
	assert.Equal(t, 5, moveCallCount)

	t.Log("Complex multi-Move call transaction validated successfully")
}

// TestMoveCallArgumentTypes tests different argument types in Move calls
func TestMoveCallArgumentTypes(t *testing.T) {
	key, err := getSigner()
	require.NoError(t, err)

	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(100000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	// Move call with different argument types
	tx.MoveCall(
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"test_module",
		"test_function",
		[]transaction.TypeTag{}, // No type arguments
		[]transaction.Argument{
			tx.Pure(uint64(42)),                                                                 // Pure value
			tx.Pure("0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"), // Pure address
			tx.Gas(),                                                                            // Gas coin reference
			tx.Object("0xe586e913e413de2df45e3eca0f0adf342a1f5d8d71e61805e14fd4872529f727"), // Object reference
		},
	)

	// Verify transaction structure
	txData, err := tx.GetTransactionData()
	require.NoError(t, err)

	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	assert.Len(t, commands, 1)
	assert.NotNil(t, commands[0].MoveCall)

	// Verify arguments
	args := commands[0].MoveCall.Arguments
	assert.Len(t, args, 4)

	t.Log("Move call with different argument types validated successfully")
}

// BenchmarkMultiMoveCall benchmarks multiple Move calls in one transaction
func BenchmarkMultiMoveCall(b *testing.B) {
	key, err := getSigner()
	require.NoError(b, err)

	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx := transaction.NewTransaction()
		tx.SetMgoClient(devCli).
			SetSigner(key).
			SetSender(model.MgoAddress(key.MgoAddress())).
			SetGasPrice(1000).
			SetGasBudget(100000000).
			SetGasOwner(model.MgoAddress(key.MgoAddress()))

		// Add 3 Move calls
		for j := 0; j < 3; j++ {
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
					tx.Pure(uint64(1000000000 * 0.01)),
				},
			)
		}

		// Build transaction
		_, err := tx.Build(true)
		require.NoError(b, err)
	}
}
