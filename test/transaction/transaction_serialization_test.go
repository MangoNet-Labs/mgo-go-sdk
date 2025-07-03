package transaction

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model"
	"github.com/mangonet-labs/mgo-go-sdk/transaction"
	"github.com/stretchr/testify/assert"
)

func TestSerializeAndDeserialize(t *testing.T) {
	devCli := client.NewMgoClient(config.RpcMgoTestnetEndpoint)

	// Create a keypair for testing
	key, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		t.Fatal(err)
	}

	// Create a transaction
	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	// Add a simple transfer operation
	receiver := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"
	splitCoin := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.01)),
	})
	tx.TransferObjects([]transaction.Argument{splitCoin}, tx.Pure(receiver))

	// Get the transaction data
	txData, err := tx.GetTransactionData()
	if err != nil {
		t.Fatal(err)
	}

	// Serialize the transaction data to JSON
	serializedJSON, err := txData.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	// Deserialize the JSON back to transaction data
	deserializedTx, err := transaction.DeserializeFromJSON(serializedJSON)
	if err != nil {
		t.Fatal(err)
	}

	// Re-serialize the deserialized transaction data
	reserializedJSON, err := deserializedTx.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(reserializedJSON)

	// Build the original transaction
	originalBytes, err := tx.Build(true)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new transaction from the deserialized data
	newTx := transaction.NewTransaction()
	newTx.SetMgoClient(devCli).
		SetSigner(key).
		SetTransactionData(deserializedTx)

	// Build the new transaction
	newBytes, err := newTx.Build(true)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the binary outputs (this is the most important test)
	assert.Equal(t, originalBytes, newBytes, "Binary outputs should match")

	// For JSON comparison, we'll parse both and compare the important fields
	// since expiration field might be auto-added during serialization
	var originalJSON, reserializedJSONParsed map[string]interface{}
	err = json.Unmarshal([]byte(serializedJSON), &originalJSON)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(reserializedJSON), &reserializedJSONParsed)
	assert.NoError(t, err)

	// Compare the important fields that should be identical
	assert.Equal(t, originalJSON["version"], reserializedJSONParsed["version"], "Version should match")
	assert.Equal(t, originalJSON["sender"], reserializedJSONParsed["sender"], "Sender should match")
	assert.Equal(t, originalJSON["gasConfig"], reserializedJSONParsed["gasConfig"], "Gas config should match")
	assert.Equal(t, originalJSON["inputs"], reserializedJSONParsed["inputs"], "Inputs should match")
	assert.Equal(t, originalJSON["transactions"], reserializedJSONParsed["transactions"], "Transactions should match")
}



