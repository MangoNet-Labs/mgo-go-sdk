package multi_transaction

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

// TestMultiTransactionBasic tests basic multi-transaction functionality
func TestMultiTransactionBasic(t *testing.T) {
	key, err := getSigner()
	require.NoError(t, err)

	// Create a transaction with multiple commands
	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"

	// Command 1: Split coins
	splitResult := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.01)), // 0.01 MGO
	})

	// Command 2: Transfer split coins
	tx.TransferObjects([]transaction.Argument{splitResult}, tx.Pure(recipient))

	// Verify transaction structure
	txData, err := tx.GetTransactionData()
	require.NoError(t, err)

	// Should have 2 commands
	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	assert.Len(t, commands, 2)

	// First command should be SplitCoins
	assert.NotNil(t, commands[0].SplitCoins)
	assert.Nil(t, commands[0].TransferObjects)

	// Second command should be TransferObjects
	assert.NotNil(t, commands[1].TransferObjects)
	assert.Nil(t, commands[1].SplitCoins)

	t.Log("Multi-transaction structure validated successfully")
}

// TestMultiTransactionChaining tests chaining multiple operations
func TestMultiTransactionChaining(t *testing.T) {
	key, err := getSigner()
	require.NoError(t, err)

	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"

	// Command 1: Split gas coin
	split1 := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.05)), // 0.05 MGO
	})

	// Command 2: Split the result again
	split2 := tx.SplitCoins(split1, []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.01)), // 0.01 MGO from the 0.05 MGO
	})

	// Command 3: Transfer the final split
	tx.TransferObjects([]transaction.Argument{split2}, tx.Pure(recipient))

	// Verify transaction structure
	txData, err := tx.GetTransactionData()
	require.NoError(t, err)

	// Should have 3 commands
	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	assert.Len(t, commands, 3)

	// Verify command types
	assert.NotNil(t, commands[0].SplitCoins)
	assert.NotNil(t, commands[1].SplitCoins)
	assert.NotNil(t, commands[2].TransferObjects)

	t.Log("Multi-transaction chaining validated successfully")
}

// TestMultiTransactionWithMoveCall tests combining MoveCall with other operations
func TestMultiTransactionWithMoveCall(t *testing.T) {
	key, err := getSigner()
	require.NoError(t, err)

	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	// Command 1: Split coins for use in MoveCall
	splitCoin := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.1)), // 0.1 MGO
	})

	// Command 2: MoveCall using the split coin
	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	require.NoError(t, err)

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
	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"
	tx.TransferObjects([]transaction.Argument{moveCallResult}, tx.Pure(recipient))

	// Verify transaction structure
	txData, err := tx.GetTransactionData()
	require.NoError(t, err)

	// Should have 3 commands
	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	assert.Len(t, commands, 3)

	// Verify command types
	assert.NotNil(t, commands[0].SplitCoins)
	assert.NotNil(t, commands[1].MoveCall)
	assert.NotNil(t, commands[2].TransferObjects)

	t.Log("Multi-transaction with MoveCall validated successfully")
}

// TestMultiTransactionMergeOperations tests merge operations in multi-transaction
func TestMultiTransactionMergeOperations(t *testing.T) {
	key, err := getSigner()
	require.NoError(t, err)

	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	// Command 1: Split gas coin into multiple parts
	split1 := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.02)), // 0.02 MGO
		tx.Pure(uint64(1000000000 * 0.03)), // 0.03 MGO
	})

	// Command 2: Split again to get more coins
	split2 := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.01)), // 0.01 MGO
	})

	// Command 3: Merge some coins together
	// Note: This is a simplified example - in practice you'd need to handle
	// multiple results from split operations correctly
	mergedCoin := tx.MergeCoins(split1, []transaction.Argument{split2})

	// Command 4: Transfer the merged coin
	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"
	tx.TransferObjects([]transaction.Argument{mergedCoin}, tx.Pure(recipient))

	// Verify transaction structure
	txData, err := tx.GetTransactionData()
	require.NoError(t, err)

	// Should have 4 commands
	commands := txData.V1.Kind.ProgrammableTransaction.Commands
	assert.Len(t, commands, 4)

	// Verify command types
	assert.NotNil(t, commands[0].SplitCoins)
	assert.NotNil(t, commands[1].SplitCoins)
	assert.NotNil(t, commands[2].MergeCoins)
	assert.NotNil(t, commands[3].TransferObjects)

	t.Log("Multi-transaction with merge operations validated successfully")
}

// TestMultiTransactionSerialization tests serialization of multi-transaction blocks
func TestMultiTransactionSerialization(t *testing.T) {
	key, err := getSigner()
	require.NoError(t, err)

	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"

	// Add multiple commands
	splitResult := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.01)),
	})
	tx.TransferObjects([]transaction.Argument{splitResult}, tx.Pure(recipient))

	// Test building the transaction
	txBytes, err := tx.Build(true)
	require.NoError(t, err)
	assert.NotEmpty(t, txBytes)

	// Test getting transaction data
	txData, err := tx.GetTransactionData()
	require.NoError(t, err)
	assert.NotNil(t, txData)

	// Test serialization
	serialized, err := txData.Marshal()
	require.NoError(t, err)
	assert.NotEmpty(t, serialized)

	t.Log("Multi-transaction serialization validated successfully")
}

// BenchmarkMultiTransaction benchmarks multi-transaction creation
func BenchmarkMultiTransaction(b *testing.B) {
	key, err := getSigner()
	require.NoError(b, err)

	recipient := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx := transaction.NewTransaction()
		tx.SetMgoClient(devCli).
			SetSigner(key).
			SetSender(model.MgoAddress(key.MgoAddress())).
			SetGasPrice(1000).
			SetGasBudget(50000000).
			SetGasOwner(model.MgoAddress(key.MgoAddress()))

		// Add multiple commands
		splitResult := tx.SplitCoins(tx.Gas(), []transaction.Argument{
			tx.Pure(uint64(1000000000 * 0.01)),
		})
		tx.TransferObjects([]transaction.Argument{splitResult}, tx.Pure(recipient))

		// Build transaction
		_, err := tx.Build(true)
		require.NoError(b, err)
	}
}
