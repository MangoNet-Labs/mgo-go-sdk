package transaction

import (
	"encoding/json"
	"testing"

	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model"
	"github.com/mangonet-labs/mgo-go-sdk/transaction"
	"github.com/stretchr/testify/assert"
)

func TestTypeArgumentSerialization(t *testing.T) {
	devCli := client.NewMgoClient(config.RpcMgoTestnetEndpoint)

	// Create a keypair for testing
	key, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		typeArgs     []transaction.TypeTag
		expectedJSON []string
	}{
		{
			name: "Primitive Types",
			typeArgs: []transaction.TypeTag{
				{Bool: new(bool)},
				{U8: new(bool)},
				{U16: new(bool)},
				{U32: new(bool)},
				{U64: new(bool)},
				{U128: new(bool)},
				{U256: new(bool)},
				{Address: new(bool)},
				{Signer: new(bool)},
			},
			expectedJSON: []string{
				"bool", "u8", "u16", "u32", "u64", "u128", "u256", "address", "signer",
			},
		},
		{
			name: "Vector Types",
			typeArgs: []transaction.TypeTag{
				{Vector: &transaction.TypeTag{U8: new(bool)}},
				{Vector: &transaction.TypeTag{U64: new(bool)}},
				{Vector: &transaction.TypeTag{Address: new(bool)}},
			},
			expectedJSON: []string{
				"vector<u8>", "vector<u64>", "vector<address>",
			},
		},
		{
			name: "Nested Vector Types",
			typeArgs: []transaction.TypeTag{
				{Vector: &transaction.TypeTag{Vector: &transaction.TypeTag{U8: new(bool)}}},
				{Vector: &transaction.TypeTag{Vector: &transaction.TypeTag{Address: new(bool)}}},
			},
			expectedJSON: []string{
				"vector<vector<u8>>", "vector<vector<address>>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a transaction with MoveCall containing type arguments
			tx := transaction.NewTransaction()
			tx.SetMgoClient(devCli).
				SetSigner(key).
				SetSender(model.MgoAddress(key.MgoAddress())).
				SetGasPrice(1000).
				SetGasBudget(50000000).
				SetGasOwner(model.MgoAddress(key.MgoAddress()))

			// Add a MoveCall with the test type arguments
			tx.MoveCall(
				"0x0000000000000000000000000000000000000000000000000000000000000002",
				"test_module",
				"test_function",
				tt.typeArgs,
				[]transaction.Argument{},
			)

			// Get the transaction data and serialize it
			txData, err := tx.GetTransactionData()
			assert.NoError(t, err)

			serializedJSON, err := txData.Serialize()
			assert.NoError(t, err)

			// Parse the JSON to verify type arguments
			var parsedJSON map[string]interface{}
			err = json.Unmarshal([]byte(serializedJSON), &parsedJSON)
			assert.NoError(t, err)

			transactions := parsedJSON["transactions"].([]interface{})
			assert.Len(t, transactions, 1)

			moveCall := transactions[0].(map[string]interface{})
			assert.Equal(t, "MoveCall", moveCall["kind"])

			typeArguments := moveCall["typeArguments"].([]interface{})
			assert.Len(t, typeArguments, len(tt.expectedJSON))

			for i, expectedType := range tt.expectedJSON {
				assert.Equal(t, expectedType, typeArguments[i])
			}

			// Test deserialization
			deserializedTx, err := transaction.DeserializeFromJSON(serializedJSON)
			assert.NoError(t, err)

			// Re-serialize and compare
			reserializedJSON, err := deserializedTx.Serialize()
			assert.NoError(t, err)

			var originalJSON, reserializedJSONParsed map[string]interface{}
			err = json.Unmarshal([]byte(serializedJSON), &originalJSON)
			assert.NoError(t, err)
			err = json.Unmarshal([]byte(reserializedJSON), &reserializedJSONParsed)
			assert.NoError(t, err)

			// Compare type arguments specifically
			origTxs := originalJSON["transactions"].([]interface{})
			reserTxs := reserializedJSONParsed["transactions"].([]interface{})

			origMoveCall := origTxs[0].(map[string]interface{})
			reserMoveCall := reserTxs[0].(map[string]interface{})

			assert.Equal(t, origMoveCall["typeArguments"], reserMoveCall["typeArguments"])
		})
	}
}

func TestStructTypeArguments(t *testing.T) {
	devCli := client.NewMgoClient(config.RpcMgoTestnetEndpoint)

	// Create a keypair for testing
	key, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		t.Fatal(err)
	}

	// Helper function to create address bytes
	createAddressBytes := func(addr string) model.MgoAddressBytes {
		addressBytes, err := transaction.ConvertMgoAddressStringToBytes(model.MgoAddress(addr))
		if err != nil {
			t.Fatal(err)
		}
		return *addressBytes
	}

	tests := []struct {
		name         string
		typeArgs     []transaction.TypeTag
		expectedJSON []string
	}{
		{
			name: "Simple Struct Types",
			typeArgs: []transaction.TypeTag{
				{
					Struct: &transaction.StructTag{
						Address: createAddressBytes("0x0000000000000000000000000000000000000000000000000000000000000002"),
						Module:  "mgo",
						Name:    "MGO",
					},
				},
				{
					Struct: &transaction.StructTag{
						Address: createAddressBytes("0x0000000000000000000000000000000000000000000000000000000000000002"),
						Module:  "coin",
						Name:    "Coin",
						TypeParams: []*transaction.TypeTag{
							{
								Struct: &transaction.StructTag{
									Address: createAddressBytes("0x0000000000000000000000000000000000000000000000000000000000000002"),
									Module:  "mgo",
									Name:    "MGO",
								},
							},
						},
					},
				},
			},
			expectedJSON: []string{
				"0x0000000000000000000000000000000000000000000000000000000000000002::mgo::MGO",
				"0x0000000000000000000000000000000000000000000000000000000000000002::coin::Coin<0x0000000000000000000000000000000000000000000000000000000000000002::mgo::MGO>",
			},
		},
		{
			name: "Complex Nested Struct Types",
			typeArgs: []transaction.TypeTag{
				{
					Struct: &transaction.StructTag{
						Address: createAddressBytes("0x0000000000000000000000000000000000000000000000000000000000000002"),
						Module:  "option",
						Name:    "Option",
						TypeParams: []*transaction.TypeTag{
							{
								Struct: &transaction.StructTag{
									Address: createAddressBytes("0x0000000000000000000000000000000000000000000000000000000000000002"),
									Module:  "coin",
									Name:    "Coin",
									TypeParams: []*transaction.TypeTag{
										{
											Struct: &transaction.StructTag{
												Address: createAddressBytes("0x0000000000000000000000000000000000000000000000000000000000000002"),
												Module:  "mgo",
												Name:    "MGO",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedJSON: []string{
				"0x0000000000000000000000000000000000000000000000000000000000000002::option::Option<0x0000000000000000000000000000000000000000000000000000000000000002::coin::Coin<0x0000000000000000000000000000000000000000000000000000000000000002::mgo::MGO>>",
			},
		},
		{
			name: "Vector of Struct Types",
			typeArgs: []transaction.TypeTag{
				{
					Vector: &transaction.TypeTag{
						Struct: &transaction.StructTag{
							Address: createAddressBytes("0x0000000000000000000000000000000000000000000000000000000000000002"),
							Module:  "mgo",
							Name:    "MGO",
						},
					},
				},
				{
					Vector: &transaction.TypeTag{
						Struct: &transaction.StructTag{
							Address: createAddressBytes("0x0000000000000000000000000000000000000000000000000000000000000002"),
							Module:  "coin",
							Name:    "Coin",
							TypeParams: []*transaction.TypeTag{
								{U64: new(bool)},
							},
						},
					},
				},
			},
			expectedJSON: []string{
				"vector<0x0000000000000000000000000000000000000000000000000000000000000002::mgo::MGO>",
				"vector<0x0000000000000000000000000000000000000000000000000000000000000002::coin::Coin<u64>>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a transaction with MoveCall containing type arguments
			tx := transaction.NewTransaction()
			tx.SetMgoClient(devCli).
				SetSigner(key).
				SetSender(model.MgoAddress(key.MgoAddress())).
				SetGasPrice(1000).
				SetGasBudget(50000000).
				SetGasOwner(model.MgoAddress(key.MgoAddress()))

			// Add a MoveCall with the test type arguments
			tx.MoveCall(
				"0x0000000000000000000000000000000000000000000000000000000000000002",
				"test_module",
				"test_function",
				tt.typeArgs,
				[]transaction.Argument{},
			)

			// Get the transaction data and serialize it
			txData, err := tx.GetTransactionData()
			assert.NoError(t, err)

			serializedJSON, err := txData.Serialize()
			assert.NoError(t, err)

			// Parse the JSON to verify type arguments
			var parsedJSON map[string]interface{}
			err = json.Unmarshal([]byte(serializedJSON), &parsedJSON)
			assert.NoError(t, err)

			transactions := parsedJSON["transactions"].([]interface{})
			assert.Len(t, transactions, 1)

			moveCall := transactions[0].(map[string]interface{})
			assert.Equal(t, "MoveCall", moveCall["kind"])

			typeArguments := moveCall["typeArguments"].([]interface{})
			assert.Len(t, typeArguments, len(tt.expectedJSON))

			for i, expectedType := range tt.expectedJSON {
				if i < len(typeArguments) {
					actualType := typeArguments[i].(string)
					assert.Equal(t, expectedType, actualType, "Type argument %d should match", i)
				} else {
					t.Errorf("Missing type argument at index %d", i)
				}
			}

			// Test deserialization
			deserializedTx, err := transaction.DeserializeFromJSON(serializedJSON)
			assert.NoError(t, err)

			// Re-serialize and compare
			reserializedJSON, err := deserializedTx.Serialize()
			assert.NoError(t, err)

			var originalJSON, reserializedJSONParsed map[string]interface{}
			err = json.Unmarshal([]byte(serializedJSON), &originalJSON)
			assert.NoError(t, err)
			err = json.Unmarshal([]byte(reserializedJSON), &reserializedJSONParsed)
			assert.NoError(t, err)

			// Compare type arguments specifically
			origTxs := originalJSON["transactions"].([]interface{})
			reserTxs := reserializedJSONParsed["transactions"].([]interface{})

			origMoveCall := origTxs[0].(map[string]interface{})
			reserMoveCall := reserTxs[0].(map[string]interface{})

			assert.Equal(t, origMoveCall["typeArguments"], reserMoveCall["typeArguments"])
		})
	}
}

func TestTypeArgumentEdgeCases(t *testing.T) {
	devCli := client.NewMgoClient(config.RpcMgoTestnetEndpoint)

	// Create a keypair for testing
	key, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		typeArgs []transaction.TypeTag
		testFunc func(t *testing.T, typeArgs []transaction.TypeTag)
	}{
		{
			name:     "Empty Type Arguments",
			typeArgs: []transaction.TypeTag{},
			testFunc: func(t *testing.T, typeArgs []transaction.TypeTag) {
				// Test that empty type arguments work correctly
				tx := transaction.NewTransaction()
				tx.SetMgoClient(devCli).
					SetSigner(key).
					SetSender(model.MgoAddress(key.MgoAddress())).
					SetGasPrice(1000).
					SetGasBudget(50000000).
					SetGasOwner(model.MgoAddress(key.MgoAddress()))

				tx.MoveCall(
					"0x0000000000000000000000000000000000000000000000000000000000000002",
					"test_module",
					"test_function",
					typeArgs,
					[]transaction.Argument{},
				)

				txData, err := tx.GetTransactionData()
				assert.NoError(t, err)

				serializedJSON, err := txData.Serialize()
				assert.NoError(t, err)

				var parsedJSON map[string]interface{}
				err = json.Unmarshal([]byte(serializedJSON), &parsedJSON)
				assert.NoError(t, err)

				transactions := parsedJSON["transactions"].([]interface{})
				moveCall := transactions[0].(map[string]interface{})
				typeArguments, exists := moveCall["typeArguments"]
				if !exists || typeArguments == nil {
					// No type arguments field or nil - this is expected for empty type arguments
					assert.Len(t, typeArgs, 0)
				} else {
					typeArgsSlice := typeArguments.([]interface{})
					assert.Len(t, typeArgsSlice, 0)
				}
			},
		},
		{
			name: "Nil Type Parameters in Struct",
			typeArgs: []transaction.TypeTag{
				{
					Struct: &transaction.StructTag{
						Address:    mustCreateAddressBytes("0x0000000000000000000000000000000000000000000000000000000000000002"),
						Module:     "mgo",
						Name:       "MGO",
						TypeParams: nil, // Test nil type params
					},
				},
			},
			testFunc: func(t *testing.T, typeArgs []transaction.TypeTag) {
				tx := transaction.NewTransaction()
				tx.SetMgoClient(devCli).
					SetSigner(key).
					SetSender(model.MgoAddress(key.MgoAddress())).
					SetGasPrice(1000).
					SetGasBudget(50000000).
					SetGasOwner(model.MgoAddress(key.MgoAddress()))

				tx.MoveCall(
					"0x0000000000000000000000000000000000000000000000000000000000000002",
					"test_module",
					"test_function",
					typeArgs,
					[]transaction.Argument{},
				)

				txData, err := tx.GetTransactionData()
				assert.NoError(t, err)

				serializedJSON, err := txData.Serialize()
				assert.NoError(t, err)

				// Test deserialization
				deserializedTx, err := transaction.DeserializeFromJSON(serializedJSON)
				assert.NoError(t, err)
				assert.NotNil(t, deserializedTx)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.typeArgs)
		})
	}
}

func TestTypeArgumentRoundTrip(t *testing.T) {
	devCli := client.NewMgoClient(config.RpcMgoTestnetEndpoint)

	// Create a keypair for testing
	key, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		t.Fatal(err)
	}

	// Test that type arguments can be serialized and deserialized correctly through JSON
	testCases := []struct {
		name     string
		typeArgs []transaction.TypeTag
	}{
		{
			name: "All Primitive Types",
			typeArgs: []transaction.TypeTag{
				{Bool: new(bool)},
				{U8: new(bool)},
				{U16: new(bool)},
				{U32: new(bool)},
				{U64: new(bool)},
				{U128: new(bool)},
				{U256: new(bool)},
				{Address: new(bool)},
				{Signer: new(bool)},
			},
		},
		{
			name: "Vector Types",
			typeArgs: []transaction.TypeTag{
				{Vector: &transaction.TypeTag{U8: new(bool)}},
				{Vector: &transaction.TypeTag{Vector: &transaction.TypeTag{U64: new(bool)}}},
			},
		},
		{
			name: "Complex Struct Types",
			typeArgs: []transaction.TypeTag{
				{
					Struct: &transaction.StructTag{
						Address:    mustCreateAddressBytes("0x0000000000000000000000000000000000000000000000000000000000000002"),
						Module:     "mgo",
						Name:       "MGO",
						TypeParams: []*transaction.TypeTag{},
					},
				},
				{
					Struct: &transaction.StructTag{
						Address: mustCreateAddressBytes("0x0000000000000000000000000000000000000000000000000000000000000002"),
						Module:  "coin",
						Name:    "Coin",
						TypeParams: []*transaction.TypeTag{
							{
								Struct: &transaction.StructTag{
									Address:    mustCreateAddressBytes("0x0000000000000000000000000000000000000000000000000000000000000002"),
									Module:     "mgo",
									Name:       "MGO",
									TypeParams: []*transaction.TypeTag{},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create transaction with type arguments
			tx := transaction.NewTransaction()
			tx.SetMgoClient(devCli).
				SetSigner(key).
				SetSender(model.MgoAddress(key.MgoAddress())).
				SetGasPrice(1000).
				SetGasBudget(50000000).
				SetGasOwner(model.MgoAddress(key.MgoAddress()))

			tx.MoveCall(
				"0x0000000000000000000000000000000000000000000000000000000000000002",
				"test_module",
				"test_function",
				tc.typeArgs,
				[]transaction.Argument{},
			)

			// Get transaction data and serialize
			txData, err := tx.GetTransactionData()
			assert.NoError(t, err)

			serializedJSON, err := txData.Serialize()
			assert.NoError(t, err)

			// Deserialize back
			deserializedTx, err := transaction.DeserializeFromJSON(serializedJSON)
			assert.NoError(t, err)

			// Re-serialize and compare
			reserializedJSON, err := deserializedTx.Serialize()
			assert.NoError(t, err)

			// Parse both JSON strings and compare type arguments
			var originalJSON, reserializedJSONParsed map[string]interface{}
			err = json.Unmarshal([]byte(serializedJSON), &originalJSON)
			assert.NoError(t, err)
			err = json.Unmarshal([]byte(reserializedJSON), &reserializedJSONParsed)
			assert.NoError(t, err)

			// Compare type arguments specifically
			origTxs := originalJSON["transactions"].([]interface{})
			reserTxs := reserializedJSONParsed["transactions"].([]interface{})

			origMoveCall := origTxs[0].(map[string]interface{})
			reserMoveCall := reserTxs[0].(map[string]interface{})

			assert.Equal(t, origMoveCall["typeArguments"], reserMoveCall["typeArguments"], "Type arguments should be preserved through round trip")

			// Verify that the type arguments are correctly parsed and serialized
			origTypeArgs := origMoveCall["typeArguments"].([]interface{})
			reserTypeArgs := reserMoveCall["typeArguments"].([]interface{})

			assert.Len(t, reserTypeArgs, len(origTypeArgs), "Number of type arguments should match")
			for i, origTypeArg := range origTypeArgs {
				assert.Equal(t, origTypeArg, reserTypeArgs[i], "Type argument %d should match exactly", i)
			}
		})
	}
}

// Helper function for tests
func mustCreateAddressBytes(addr string) model.MgoAddressBytes {
	addressBytes, err := transaction.ConvertMgoAddressStringToBytes(model.MgoAddress(addr))
	if err != nil {
		panic(err)
	}
	return *addressBytes
}
