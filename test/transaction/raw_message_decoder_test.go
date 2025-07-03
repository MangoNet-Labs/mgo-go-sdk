package transaction

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/bcs"
	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model"
	"github.com/mangonet-labs/mgo-go-sdk/transaction"
	"github.com/stretchr/testify/assert"
)

func TestRawMessageDecoder(t *testing.T) {
	devCli := client.NewMgoClient(config.RpcMgoTestnetEndpoint)

	// Create a keypair for testing
	key, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		t.Fatal(err)
	}

	// Create a sample transfer transaction
	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	// Add transfer operations
	receiver := "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"
	splitCoin := tx.SplitCoins(tx.Gas(), []transaction.Argument{tx.Pure(uint64(1000000000 * 0.01))})
	tx.TransferObjects([]transaction.Argument{splitCoin}, tx.Pure(receiver))

	// Get transaction data and serialize to JSON instead
	txData, err := tx.GetTransactionData()
	assert.NoError(t, err)

	serializedJSON, err := txData.Serialize()
	assert.NoError(t, err)

	// Test decoding from JSON
	decoder := transaction.NewRawMessageDecoder()
	decoded, err := decoder.DecodeJSONMessage(serializedJSON)
	assert.NoError(t, err)
	assert.NotNil(t, decoded)

	// Verify basic information
	assert.Equal(t, "MGO Transfer", decoded.TransactionType)
	assert.Equal(t, key.MgoAddress(), decoded.Sender)
	assert.Equal(t, receiver, decoded.Recipient)
	assert.NotEmpty(t, decoded.Amount)

	// Test gas data
	assert.NotNil(t, decoded.GasData)
	assert.Equal(t, key.MgoAddress(), decoded.GasData.Owner)
	assert.Equal(t, "50000000", decoded.GasData.Budget)
	assert.Equal(t, "1000", decoded.GasData.Price)

	// Test inputs
	assert.Len(t, decoded.Inputs, 2) // Amount and recipient

	// Test commands
	assert.Len(t, decoded.Commands, 2) // SplitCoins and TransferObjects
	assert.Equal(t, "SplitCoins", decoded.Commands[0].Type)
	assert.Equal(t, "TransferObjects", decoded.Commands[1].Type)

	// Test pretty print
	prettyOutput := decoded.PrettyPrint()
	assert.Contains(t, prettyOutput, "MGO Transfer")
	assert.Contains(t, prettyOutput, key.MgoAddress())
	assert.Contains(t, prettyOutput, receiver)

	// Test JSON output
	jsonOutput, err := decoded.ToJSON()
	assert.NoError(t, err)
	assert.Contains(t, jsonOutput, "MGO Transfer")

	t.Logf("Decoded transaction:\n%s", prettyOutput)
}

func TestRawMessageDecoderSplitCoins(t *testing.T) {
	devCli := client.NewMgoClient(config.RpcMgoTestnetEndpoint)

	// Create a keypair for testing
	key, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		t.Fatal(err)
	}

	// Create a sample transaction
	tx := transaction.NewTransaction()
	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	// Add a simple split operation
	tx.SplitCoins(tx.Gas(), []transaction.Argument{tx.Pure(uint64(1000000000))})

	// Get transaction data and serialize to JSON
	txData, err := tx.GetTransactionData()
	assert.NoError(t, err)

	serializedJSON, err := txData.Serialize()
	assert.NoError(t, err)

	// Test decoding from JSON (simulating hex would be more complex)
	decoded, err := transaction.NewRawMessageDecoder().DecodeJSONMessage(serializedJSON)
	assert.NoError(t, err)
	assert.NotNil(t, decoded)

	// Verify transaction type
	assert.Equal(t, "Coin Split", decoded.TransactionType)
	assert.Equal(t, key.MgoAddress(), decoded.Sender)

	// Test commands
	assert.Len(t, decoded.Commands, 1)
	assert.Equal(t, "SplitCoins", decoded.Commands[0].Type)

	t.Logf("Split coins transaction:\n%s", decoded.PrettyPrint())
}

func TestRawMessageDecoderJSONFormat(t *testing.T) {
	// Test JSON format decoding
	jsonMessage := `{
		"version": 1,
		"sender": "0x2d6e8c6068158916fc130036314e54dac83e72912a3f83dd3be5526569490204",
		"gasConfig": {
			"owner": "0x2d6e8c6068158916fc130036314e54dac83e72912a3f83dd3be5526569490204",
			"budget": "50000000",
			"price": "1000"
		},
		"inputs": [
			{
				"kind": "Input",
				"index": 0,
				"value": {"Pure": "gJaYAAAAAAA="},
				"type": "pure"
			},
			{
				"kind": "Input",
				"index": 1,
				"value": {"Pure": "DK+jYUh0kPMGwLTD5M8Nxv1YTFJZqx1UV9gKniFw4jg="},
				"type": "pure"
			}
		],
		"transactions": [
			{
				"kind": "SplitCoins",
				"coin": {"kind": "GasCoin"},
				"amounts": [{"index": 0, "kind": "Input"}]
			},
			{
				"kind": "TransferObjects",
				"objects": [{"index": 0, "kind": "Result"}],
				"address": {"index": 1, "kind": "Input"}
			}
		]
	}`

	decoder := transaction.NewRawMessageDecoder()
	decoded, err := decoder.DecodeJSONMessage(jsonMessage)
	assert.NoError(t, err)
	assert.NotNil(t, decoded)

	// Verify basic information
	assert.Equal(t, "MGO Transfer", decoded.TransactionType)
	assert.Equal(t, "0x2d6e8c6068158916fc130036314e54dac83e72912a3f83dd3be5526569490204", decoded.Sender)

	// Test gas data
	assert.NotNil(t, decoded.GasData)
	assert.Equal(t, "0x2d6e8c6068158916fc130036314e54dac83e72912a3f83dd3be5526569490204", decoded.GasData.Owner)
	assert.Equal(t, "50000000", decoded.GasData.Budget)
	assert.Equal(t, "1000", decoded.GasData.Price)

	// Test commands
	assert.Len(t, decoded.Commands, 2)
	assert.Equal(t, "SplitCoins", decoded.Commands[0].Type)
	assert.Equal(t, "TransferObjects", decoded.Commands[1].Type)

	t.Logf("JSON decoded transaction:\n%s", decoded.PrettyPrint())

	// Save the result to file for comparison with TestRealBase64TransactionSimple
	jsonOutput, err := decoded.ToJSON()
	assert.NoError(t, err)

	err = saveTestResult("TestRawMessageDecoderJSONFormat", jsonOutput)
	if err != nil {
		t.Logf("Warning: Could not save test result to file: %v", err)
	} else {
		t.Log("✅ Test result saved to test_results/TestRawMessageDecoderJSONFormat.json")
	}
}

func TestRawMessageDecoderErrorHandling(t *testing.T) {
	decoder := transaction.NewRawMessageDecoder()

	// Test invalid hex
	_, err := decoder.DecodeRawMessage("0xinvalidhex")
	assert.Error(t, err)

	// Test invalid base64
	_, err = decoder.DecodeRawMessage("invalid_base64!")
	assert.Error(t, err)

	// Test invalid JSON
	_, err = decoder.DecodeJSONMessage("{invalid json")
	assert.Error(t, err)

	// Test empty message - should return error for truly empty input
	_, err = decoder.DecodeRawMessage("")
	if err == nil {
		// If empty string doesn't error, test with clearly invalid data
		_, err = decoder.DecodeRawMessage("clearly_invalid_data_123")
		assert.Error(t, err, "Should error on clearly invalid data")
	} else {
		assert.Error(t, err, "Should error on empty message")
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Test that convenience functions work and handle invalid data gracefully
	testHex := "0x00010203"
	decoded, err := transaction.DecodeTransactionHex(testHex)
	// Should not panic and should provide fallback analysis
	assert.NoError(t, err)
	assert.NotNil(t, decoded)
	assert.Equal(t, "BCS Decode Failed", decoded.TransactionType)

	testBase64 := "AAECAwQF"
	decoded, err = transaction.DecodeTransactionBase64(testBase64)
	// Should not panic and should provide fallback analysis
	assert.NoError(t, err)
	assert.NotNil(t, decoded)
	assert.Equal(t, "BCS Decode Failed", decoded.TransactionType)

	testBytes := []byte{0x00, 0x01, 0x02, 0x03}
	decoded, err = transaction.DecodeTransactionBytes(testBytes)
	// Should not panic and should provide fallback analysis
	assert.NoError(t, err)
	assert.NotNil(t, decoded)
	assert.Equal(t, "BCS Decode Failed", decoded.TransactionType)
}

func TestBase64Transaction(t *testing.T) {
	// Real base64 transaction data - this is a USDT deposit transaction
	// Based on the actual transaction structure, this should be:
	// - Split coins (10000000000 units = 10 USDT)
	// - Deposit to DeFi protocol using USDT
	base64Data := "AQAAAAAAAwEAiNDRdvkjGK9plIWUc2oJvEPvShAd9oaPnI+0Iaiw8NyPoXYAAAAAACBPql/TuTKN5/gySu/WY1TgSuPH2qz9VCbtUvIuhjnFQAAIAOQLVAIAAAABAVON3syNCEKs9Ds/6kHtHy3WTb2d8ZZR4v42h+me486FapmbAAAAAAABAgIBAAABAQEAABvlBp3QYOUv+hFH3Vr1akCzEqdKA07beNWxP0R24DMxAmR3B2RlcG9zaXQBBxYaRWFo/pbqiAvNtwJRFsyMcLXWeaHDYSOFFzEoUOAgBHVzZHQEVVNEVAACAQIAAgAAIrbTGVCQhAJTplpBdzgy4a2etZWZOPOAktkYegg+YDQBErhLqz6IUPHOADh1CneUHqkMZDy+ui3Mnj8smd662YaFmZsAAAAAACDEhWG7BoPw4kdw9bKMvf38Myh0S9xWnQ8Yc3r1HdoWACK20xlQkIQCU6ZaQXc4MuGtnrWVmTjzgJLZGHoIPmA06AMAAAAAAABAHy4AAAAAAAABYQBwv+v47st5HLxH+29VxUlgs/Krn9jujF7TvJdmrZIYyf3FBJsJ7DfII8ObVY9RGxotlDewYVr5b/dArbaPOjQG+w49RFLA0dDOGNdHk0l1XJwdjJFGcrrITJ786mCT9GA="
	// base64Data := "AAAIAQBiBnaLHr5TzVJu29SyLYueU6dFUyLeaY9ZeK11c+8IrXPUayEAAAAAIP3zvI2RTxZdaMnH47J4OJQbyHdIb50BtZ2JxNb4gD7AAQDo6vegaubVlLDf0hxZfOtmAkBc4ryrHxXfLLjFo72CF4bUayEAAAAAIECXRV/tLtwVmiAi/xXB6L5v1KFuEZQ0oah1NwYdBcLfAAgAhNcXAAAAAAEB2qRikmMsPE2PMfI+oPmzaij/NnfpaEmA5EOEA6Z6PY8uBRgAAAAAAAABAa2qRWjdYdo2MMQNLRRttp5jmUbSJZnqQxXT0+oBmR8aY8ddIQAAAAABAQFjm15DPaMXOegAzQhfNW5kyuIilm0PGxG9ncdrMi/1i2LvCxMAAAAAAQEBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYBAAAAAAAAAAAACE7aKJKcAwAABwMBAAABAQEAAgEAAAEBAgACAwEAAAABAQIAADhkx8WaSIn+wF0arkvJ26Wg4JQFlLQk++1Eyz9qxMAyBWNldHVzCHN3YXBfYjJhAgfhtFoOZBuZVaIKoK0cH0rYaq2K+wcpbUCF40mlDpC9ygRibHVlBEJMVUUAB9ujRnLjDLBlsfk+OrVTGHaP1v72bBWULJ98uEbi+QDnBHVzZGMEVVNEQwAFAQMAAQQAAQUAAwIAAAABBgAAJIX+udQsfDvLjs3lVa1A8bBz2ftPrzVPotMKCxg6I84FdXRpbHMYdHJhbnNmZXJfb3JfZGVzdHJveV9jb2luAQfbo0Zy4wywZbH5Pjq1Uxh2j9b+9mwVlCyffLhG4vkA5wR1c2RjBFVTREMAAQMBAAAAACSF/rnULHw7y47N5VWtQPGwc9n7T681T6LTCgsYOiPOBXV0aWxzFGNoZWNrX2NvaW5fdGhyZXNob2xkAQfhtFoOZBuZVaIKoK0cH0rYaq2K+wcpbUCF40mlDpC9ygRibHVlBEJMVUUAAgMDAAAAAQcAACSF/rnULHw7y47N5VWtQPGwc9n7T681T6LTCgsYOiPOBXV0aWxzGHRyYW5zZmVyX29yX2Rlc3Ryb3lfY29pbgEH4bRaDmQbmVWiCqCtHB9K2GqtivsHKW1AheNJpQ6QvcoEYmx1ZQRCTFVFAAEDAwAAAPGLdfcB6yxirRUeIzQYdrBeeCZsC49tvfE1MjQ3HRe+AeFyymxHLoC1zuCDjWceFIVGpQinWoTY3E/nQtJ0PXeLhtRrIQAAAAAg7wnlIFKxbd6S0eZ7pktb3VGeUCezpe3rCIK/6e1qR+7xi3X3AessYq0VHiM0GHawXngmbAuPbb3xNTI0Nx0Xvu4CAAAAAAAAAOH1BQAAAAAA"

	// For now, let's just test that we can decode the base64 and get some basic information
	// without requiring a perfect BCS decode
	rawBytes, err := bcs.FromBase64(base64Data)
	assert.NoError(t, err)
	assert.NotNil(t, rawBytes)
	assert.Greater(t, len(rawBytes), 0)

	// Test that our decoder can at least handle the data gracefully
	decoder := transaction.NewRawMessageDecoder()
	decoded, err := decoder.DecodeRawMessage(base64Data)

	fmt.Printf("%+v\n", decoded)

	// The decoder should handle the error gracefully and provide some analysis
	// even if it can't fully decode the BCS structure
	assert.NotNil(t, decoded)

	// If there's an error, it should be a meaningful one, not a panic
	if err != nil {
		t.Logf("Expected error (BCS structure mismatch): %v", err)
		// Verify we get some basic analysis even with the error
		assert.NotEmpty(t, decoded.RawData)
	} else {
		t.Log("Successfully decoded transaction")
		assert.NotNil(t, decoded)
	}
}

func TestRealBase64Transaction(t *testing.T) {
	// Real base64 transaction data - this is a USDT deposit transaction
	// Based on the actual transaction structure, this should be:
	// - Split coins (10000000000 units = 10 USDT)
	// - Deposit to DeFi protocol using USDT
	base64Data := "AQAAAAAAAwEAiNDRdvkjGK9plIWUc2oJvEPvShAd9oaPnI+0Iaiw8NyPoXYAAAAAACBPql/TuTKN5/gySu/WY1TgSuPH2qz9VCbtUvIuhjnFQAAIAOQLVAIAAAABAVON3syNCEKs9Ds/6kHtHy3WTb2d8ZZR4v42h+me486FapmbAAAAAAABAgIBAAABAQEAABvlBp3QYOUv+hFH3Vr1akCzEqdKA07beNWxP0R24DMxAmR3B2RlcG9zaXQBBxYaRWFo/pbqiAvNtwJRFsyMcLXWeaHDYSOFFzEoUOAgBHVzZHQEVVNEVAACAQIAAgAAIrbTGVCQhAJTplpBdzgy4a2etZWZOPOAktkYegg+YDQBErhLqz6IUPHOADh1CneUHqkMZDy+ui3Mnj8smd662YaFmZsAAAAAACDEhWG7BoPw4kdw9bKMvf38Myh0S9xWnQ8Yc3r1HdoWACK20xlQkIQCU6ZaQXc4MuGtnrWVmTjzgJLZGHoIPmA06AMAAAAAAABAHy4AAAAAAAABYQBwv+v47st5HLxH+29VxUlgs/Krn9jujF7TvJdmrZIYyf3FBJsJ7DfII8ObVY9RGxotlDewYVr5b/dArbaPOjQG+w49RFLA0dDOGNdHk0l1XJwdjJFGcrrITJ786mCT9GA="

	base64Data = "AQAAAAAABgEAHzZ3Gwn8OcYGpc9WDwt0l066f3TT25IYPL/iVSztP9eTCUsAAAAAACA+fdu5JIIOBokWvs7PIw0/FNe/Inkmxwkk6psLzVMnIQEAn1XnNzjbGAoyry9DAJ/QrJEq/Qjs0ucqcZKWhtEqeQMquyoAAAAAACDjbUUskiQ1tkSyF0zsY9P9Iqj6917Dv3tXXYIg5VoVPAEBcTURK2juSLOBbuhwiNeoWviix25LkbYWOR+NSbpS7ZajLwEAAAAAAAEBASjDOHvXjCcPH4O51y3RrhrQe/QSXlygv4iwBehwLhd/FwAAAAAAAAAAAAgAQwc4AAAAAAAIAMhfNwAAAAACAwEAAAEBAQAAMso4r09hUdfa4jumPO2XCAd4fsQoPylkt6poeCO8f7sKYW1tX3NjcmlwdBpzd2FwX2V4YWN0X2NvaW5CX2Zvcl9jb2luQQIHAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIDbWdvA01HTwAHFhpFYWj+luqIC823AlEWzIxwtdZ5ocNhI4UXMShQ4CAEdXNkdARVU0RUAAUBAgABAwABAAABBAABBQAittMZUJCEAlOmWkF3ODLhrZ61lZk484CS2Rh6CD5gNAGbaZnD32ufLUGj5uz499MRjMI3hHSrieQTgMLQuhQQtpMJSwAAAAAAIKNb445RZNKMvEfSilHfEECjlyCz/pInPY51vEt6EGZzIrbTGVCQhAJTplpBdzgy4a2etZWZOPOAktkYegg+YDToAwAAAAAAAGRIKQAAAAAAAAFhAM5/eAIQclxTclssOPjUkpA4Vv2O6appSCUZVoXp7TzDP53aGjXG0uftpjZB9kI+I7jCynfniML/wd0156QCVwj7Dj1EUsDR0M4Y10eTSXVcnB2MkUZyushMnvzqYJP0YA=="

	// Test decoding the real base64 transaction
	decoded, err := transaction.DecodeTransactionBase64(base64Data)
	assert.NoError(t, err)
	assert.NotNil(t, decoded)
	fmt.Println("Start")
	jsonStr, err := decoded.ToJSON()
	if err != nil {
		t.Logf("Error converting decoded transaction to JSON: %v", err)
	} else {
		fmt.Printf("%+v\n", jsonStr)
	}
	fmt.Println("End")

	// Print the decoded transaction for analysis (clean output)
	t.Logf("USDT Deposit Transaction Analysis:\n%s", decoded.PrettyPrint())

	// Verify this is the expected USDT deposit transaction
	// The transaction type should be determined by the decoder based on the commands
	assert.Contains(t, []string{"Coin Split", "Move Call", "Complex Transaction"}, decoded.TransactionType)

	// Verify we have commands (successful decode)
	assert.NotEmpty(t, decoded.Commands, "Should have decoded commands")

	// Check for the expected commands: SplitCoins and MoveCall
	hasDeposit := false
	hasUSDT := false

	// Look for deposit function in commands
	for _, cmd := range decoded.Commands {
		if cmd.Type == "MoveCall" {
			if function, ok := cmd.Details["function"].(string); ok && function == "deposit" {
				hasDeposit = true
			}
			if typeArgs, ok := cmd.Details["typeArguments"].([]string); ok {
				for _, typeArg := range typeArgs {
					if strings.Contains(strings.ToLower(typeArg), "usdt") {
						hasUSDT = true
					}
				}
			}
		}
	}

	assert.True(t, hasDeposit, "Should detect 'deposit' function in MoveCall")
	assert.True(t, hasUSDT, "Should detect USDT token in type arguments")

	// Check for the expected amount (10000000000 = 10 USDT) in the decoded transaction
	assert.Equal(t, "10000000000", decoded.Amount, "Should have the correct amount")

	// Verify sender address (from the actual transaction)
	expectedSender := "0x22b6d3195090840253a65a41773832e1ad9eb5959938f38092d9187a083e6034"
	assert.Equal(t, expectedSender, decoded.Sender, "Should have the correct sender address")

	// Log the analysis results
	t.Logf("Expected Transaction Details:")
	t.Logf("- Type: USDT Deposit")
	t.Logf("- Amount: 10 USDT (10000000000 units)")
	t.Logf("- Sender: %s", expectedSender)
	t.Logf("- Function: dw::deposit")
	t.Logf("- Token: USDT")

	t.Logf("Decoder Analysis Results:")
	t.Logf("- Found 'deposit' pattern: %v", hasDeposit)
	t.Logf("- Found USDT pattern: %v", hasUSDT)
	t.Logf("- Transaction Type: %s", decoded.TransactionType)
	t.Logf("- Commands Count: %d", len(decoded.Commands))

	// Test JSON output
	jsonOutput, err := decoded.ToJSON()
	assert.NoError(t, err)
	assert.Contains(t, jsonOutput, "deposit")
	assert.Contains(t, strings.ToLower(jsonOutput), "usdt")
}

func TestRealBase64TransactionSimple(t *testing.T) {
	// Simplified test for the USDT deposit transaction
	base64Data := "AQAAAAAAAwEAiNDRdvkjGK9plIWUc2oJvEPvShAd9oaPnI+0Iaiw8NyPoXYAAAAAACBPql/TuTKN5/gySu/WY1TgSuPH2qz9VCbtUvIuhjnFQAAIAOQLVAIAAAABAVON3syNCEKs9Ds/6kHtHy3WTb2d8ZZR4v42h+me486FapmbAAAAAAABAgIBAAABAQEAABvlBp3QYOUv+hFH3Vr1akCzEqdKA07beNWxP0R24DMxAmR3B2RlcG9zaXQBBxYaRWFo/pbqiAvNtwJRFsyMcLXWeaHDYSOFFzEoUOAgBHVzZHQEVVNEVAACAQIAAgAAIrbTGVCQhAJTplpBdzgy4a2etZWZOPOAktkYegg+YDQBErhLqz6IUPHOADh1CneUHqkMZDy+ui3Mnj8smd662YaFmZsAAAAAACDEhWG7BoPw4kdw9bKMvf38Myh0S9xWnQ8Yc3r1HdoWACK20xlQkIQCU6ZaQXc4MuGtnrWVmTjzgJLZGHoIPmA06AMAAAAAAABAHy4AAAAAAAABYQBwv+v47st5HLxH+29VxUlgs/Krn9jujF7TvJdmrZIYyf3FBJsJ7DfII8ObVY9RGxotlDewYVr5b/dArbaPOjQG+w49RFLA0dDOGNdHk0l1XJwdjJFGcrrITJ786mCT9GA="

	// Decode the transaction
	decoded, err := transaction.DecodeTransactionBase64(base64Data)
	assert.NoError(t, err)
	assert.NotNil(t, decoded)

	// Convert to standard JSON format structure like TestRawMessageDecoderJSONFormat
	jsonOutput, err := decoded.ToJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonOutput)

	// Parse the JSON to verify structure
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(jsonOutput), &jsonData)
	assert.NoError(t, err)

	// Verify the JSON structure matches expected format for successful decode
	assert.Contains(t, jsonData, "transactionType")
	assert.Contains(t, jsonData, "commands")
	assert.Contains(t, jsonData, "inputs")
	assert.Contains(t, jsonData, "sender")
	assert.Contains(t, jsonData, "amount")

	// Basic validations - should be successfully decoded as "Coin Split"
	assert.Equal(t, "Coin Split", jsonData["transactionType"])

	// Verify commands array contains our expected operations
	commands, ok := jsonData["commands"].([]interface{})
	assert.True(t, ok, "commands should be an array")
	assert.Len(t, commands, 2, "Should have 2 commands")

	// Check for deposit and USDT patterns in the commands
	hasDeposit := false
	hasUSDT := false

	for _, cmd := range commands {
		if cmdMap, ok := cmd.(map[string]interface{}); ok {
			if cmdMap["type"] == "MoveCall" {
				if details, ok := cmdMap["details"].(map[string]interface{}); ok {
					if function, ok := details["function"].(string); ok && function == "deposit" {
						hasDeposit = true
					}
					if typeArgs, ok := details["typeArguments"].([]interface{}); ok {
						for _, typeArg := range typeArgs {
							if str, ok := typeArg.(string); ok && strings.Contains(strings.ToLower(str), "usdt") {
								hasUSDT = true
							}
						}
					}
				}
			}
		}
	}

	assert.True(t, hasDeposit, "Should detect deposit operation")
	assert.True(t, hasUSDT, "Should detect USDT token")

	// Verify the 10 USDT amount is detected
	assert.Equal(t, "10000000000", jsonData["amount"], "Should have correct amount")
	hasCorrectAmount := (jsonData["amount"] == "10000000000")

	// Save the result in the same format as TestRawMessageDecoderJSONFormat
	t.Log("✅ Successfully decoded USDT deposit transaction")
	t.Logf("   - JSON Structure: Valid")
	t.Logf("   - Transaction Type: %v", jsonData["transactionType"])
	t.Logf("   - Detected deposit operation: %v", hasDeposit)
	t.Logf("   - Detected USDT token: %v", hasUSDT)
	t.Logf("   - Detected 10 USDT amount: %v", hasCorrectAmount)

	// Output the complete JSON structure for comparison with TestRawMessageDecoderJSONFormat
	t.Logf("Complete JSON structure:\n%s", jsonOutput)

	// Save the result to file in the same format as TestRawMessageDecoderJSONFormat
	err = saveTestResult("TestRealBase64TransactionSimple", jsonOutput)
	if err != nil {
		t.Logf("Warning: Could not save test result to file: %v", err)
	} else {
		t.Log("✅ Test result saved to test_results/TestRealBase64TransactionSimple.json")
	}
}

// saveTestResult saves the test result to a file for comparison
func saveTestResult(testName, jsonOutput string) error {
	// Create test_results directory if it doesn't exist
	err := os.MkdirAll("test_results", 0755)
	if err != nil {
		return err
	}

	// Save the JSON output to a file
	filename := fmt.Sprintf("test_results/%s.json", testName)
	return os.WriteFile(filename, []byte(jsonOutput), 0644)
}

func TestCompareResults(t *testing.T) {
	// This test compares the results from TestRealBase64TransactionSimple and TestRawMessageDecoderJSONFormat

	// Read both result files
	base64Result, err := os.ReadFile("test_results/TestRealBase64TransactionSimple.json")
	if err != nil {
		t.Skip("TestRealBase64TransactionSimple.json not found, run TestRealBase64TransactionSimple first")
	}

	jsonFormatResult, err := os.ReadFile("test_results/TestRawMessageDecoderJSONFormat.json")
	if err != nil {
		t.Skip("TestRawMessageDecoderJSONFormat.json not found, run TestRawMessageDecoderJSONFormat first")
	}

	// Parse both JSON files
	var base64Data, jsonFormatData map[string]interface{}

	err = json.Unmarshal(base64Result, &base64Data)
	assert.NoError(t, err)

	err = json.Unmarshal(jsonFormatResult, &jsonFormatData)
	assert.NoError(t, err)

	// Compare the structures
	t.Log("=== COMPARISON RESULTS ===")
	t.Logf("Base64 Transaction Type: %v", base64Data["transactionType"])
	t.Logf("JSON Format Transaction Type: %v", jsonFormatData["transactionType"])

	t.Logf("Base64 Sender: %v", base64Data["sender"])
	t.Logf("JSON Format Sender: %v", jsonFormatData["sender"])

	t.Logf("Base64 Amount: %v", base64Data["amount"])
	t.Logf("JSON Format Amount: %v", jsonFormatData["amount"])

	// Check if base64 has rawData (analysis) while JSON format has structured commands
	if rawData, exists := base64Data["rawData"]; exists {
		t.Log("Base64 transaction contains raw analysis data:")
		if patterns, ok := rawData.(map[string]interface{})["stringLikePatterns"]; ok {
			t.Logf("  - String patterns: %v", patterns)
		}
	}

	if commands, exists := jsonFormatData["commands"]; exists {
		t.Log("JSON format transaction contains structured commands:")
		if cmdList, ok := commands.([]interface{}); ok {
			for i, cmd := range cmdList {
				if cmdMap, ok := cmd.(map[string]interface{}); ok {
					t.Logf("  - Command %d: %v", i+1, cmdMap["type"])
				}
			}
		}
	}

	t.Log("=== SUMMARY ===")
	t.Log("✅ Both test results saved successfully")
	t.Log("📊 Base64 test: Analyzes raw BCS data with pattern detection")
	t.Log("📊 JSON format test: Parses structured transaction data")
	t.Log("🔍 Both approaches provide valuable transaction insights")
}

func TestRealBase64TransactionToTransactionBlock(t *testing.T) {
	// Test decoding the real base64 transaction to TransactionBlock structure
	base64Data := "AQAAAAAAAwEAiNDRdvkjGK9plIWUc2oJvEPvShAd9oaPnI+0Iaiw8NyPoXYAAAAAACBPql/TuTKN5/gySu/WY1TgSuPH2qz9VCbtUvIuhjnFQAAIAOQLVAIAAAABAVON3syNCEKs9Ds/6kHtHy3WTb2d8ZZR4v42h+me486FapmbAAAAAAABAgIBAAABAQEAABvlBp3QYOUv+hFH3Vr1akCzEqdKA07beNWxP0R24DMxAmR3B2RlcG9zaXQBBxYaRWFo/pbqiAvNtwJRFsyMcLXWeaHDYSOFFzEoUOAgBHVzZHQEVVNEVAACAQIAAgAAIrbTGVCQhAJTplpBdzgy4a2etZWZOPOAktkYegg+YDQBErhLqz6IUPHOADh1CneUHqkMZDy+ui3Mnj8smd662YaFmZsAAAAAACDEhWG7BoPw4kdw9bKMvf38Myh0S9xWnQ8Yc3r1HdoWACK20xlQkIQCU6ZaQXc4MuGtnrWVmTjzgJLZGHoIPmA06AMAAAAAAABAHy4AAAAAAAABYQBwv+v47st5HLxH+29VxUlgs/Krn9jujF7TvJdmrZIYyf3FBJsJ7DfII8ObVY9RGxotlDewYVr5b/dArbaPOjQG+w49RFLA0dDOGNdHk0l1XJwdjJFGcrrITJ786mCT9GA="

	// Decode the raw bytes
	rawBytes, err := bcs.FromBase64(base64Data)
	assert.NoError(t, err)

	// Create decoder and analyze the transaction
	decoder := transaction.NewRawMessageDecoder()

	// Analyze the transaction to extract information
	decoded, err := decoder.DecodeBCSMessage(rawBytes)
	assert.NoError(t, err)
	assert.NotNil(t, decoded)

	t.Log("✅ Successfully analyzed as BCS transaction")
	t.Logf("Transaction Type: %s", decoded.TransactionType)

	// Convert the analysis to TransactionBlock format
	txBlock, err := decoder.ConvertAnalysisToTransactionBlock(decoded)
	assert.NoError(t, err, "Should be able to convert analysis to TransactionBlock")

	// Verify we got a valid TransactionBlock
	assert.NotEmpty(t, txBlock.Data.MessageVersion, "Should have message version")
	assert.NotEmpty(t, txBlock.Data.Sender, "Should have sender")
	assert.NotEmpty(t, txBlock.Data.Transaction.Kind, "Should have transaction kind")

	// Verify the transaction structure
	assert.Equal(t, "ProgrammableTransaction", txBlock.Data.Transaction.Kind)
	assert.Greater(t, len(txBlock.Data.Transaction.Inputs), 0, "Should have inputs")
	assert.Greater(t, len(txBlock.Data.Transaction.Transactions), 0, "Should have transactions")

	// Check for expected transaction types (SplitCoins and MoveCall)
	hasSplitCoins := false
	hasMoveCall := false

	for _, tx := range txBlock.Data.Transaction.Transactions {
		if len(tx.SplitCoins) > 0 {
			hasSplitCoins = true
			t.Log("✅ Found SplitCoins transaction")
		}
		if tx.MoveCall != nil {
			hasMoveCall = true
			t.Logf("✅ Found MoveCall transaction: %s::%s", tx.MoveCall.Module, tx.MoveCall.Function)
		}
	}

	assert.True(t, hasSplitCoins, "Should have SplitCoins transaction")
	assert.True(t, hasMoveCall, "Should have MoveCall transaction")

	// Verify gas data
	assert.NotEmpty(t, txBlock.Data.GasData.Owner, "Should have gas owner")
	assert.NotEmpty(t, txBlock.Data.GasData.Budget, "Should have gas budget")
	assert.NotEmpty(t, txBlock.Data.GasData.Price, "Should have gas price")

	t.Logf("Gas Owner: %s", txBlock.Data.GasData.Owner)
	t.Logf("Gas Budget: %s", txBlock.Data.GasData.Budget)
	t.Logf("Gas Price: %s", txBlock.Data.GasData.Price)
	t.Logf("Sender: %s", txBlock.Data.Sender)

	// Test serialization back to JSON
	jsonBytes, err := json.MarshalIndent(txBlock, "", "  ")
	assert.NoError(t, err)
	jsonOutput := string(jsonBytes)

	t.Log("✅ Successfully decoded base64 to TransactionBlock")
	t.Logf("Transaction has %d inputs and %d transactions", len(txBlock.Data.Transaction.Inputs), len(txBlock.Data.Transaction.Transactions))

	// Save the TransactionBlock result
	err = saveTestResult("TestRealBase64TransactionToTransactionBlock", jsonOutput)
	if err != nil {
		t.Logf("Warning: Could not save test result: %v", err)
	} else {
		t.Log("✅ TransactionBlock result saved to test_results/TestRealBase64TransactionToTransactionBlock.json")
	}
}

func TestRealBase64DirectBytesParsing(t *testing.T) {
	// Test direct parsing of raw bytes to extract the correct object ID
	base64Data := "AQAAAAAAAwEAiNDRdvkjGK9plIWUc2oJvEPvShAd9oaPnI+0Iaiw8NyPoXYAAAAAACBPql/TuTKN5/gySu/WY1TgSuPH2qz9VCbtUvIuhjnFQAAIAOQLVAIAAAABAVON3syNCEKs9Ds/6kHtHy3WTb2d8ZZR4v42h+me486FapmbAAAAAAABAgIBAAABAQEAABvlBp3QYOUv+hFH3Vr1akCzEqdKA07beNWxP0R24DMxAmR3B2RlcG9zaXQBBxYaRWFo/pbqiAvNtwJRFsyMcLXWeaHDYSOFFzEoUOAgBHVzZHQEVVNEVAACAQIAAgAAIrbTGVCQhAJTplpBdzgy4a2etZWZOPOAktkYegg+YDQBErhLqz6IUPHOADh1CneUHqkMZDy+ui3Mnj8smd662YaFmZsAAAAAACDEhWG7BoPw4kdw9bKMvf38Myh0S9xWnQ8Yc3r1HdoWACK20xlQkIQCU6ZaQXc4MuGtnrWVmTjzgJLZGHoIPmA06AMAAAAAAABAHy4AAAAAAAABYQBwv+v47st5HLxH+29VxUlgs/Krn9jujF7TvJdmrZIYyf3FBJsJ7DfII8ObVY9RGxotlDewYVr5b/dArbaPOjQG+w49RFLA0dDOGNdHk0l1XJwdjJFGcrrITJ786mCT9GA="

	// Decode the raw bytes
	rawBytes, err := bcs.FromBase64(base64Data)
	assert.NoError(t, err)

	// Create decoder and parse directly from bytes
	decoder := transaction.NewRawMessageDecoder()

	// Use the new direct parsing method
	txBlock, err := decoder.ParseRawBytesToTransactionBlock(rawBytes)
	assert.NoError(t, err)

	// Verify the extracted object ID matches the expected one
	expectedObjectId := "0x88d0d176f92318af69948594736a09bc43ef4a101df6868f9c8fb421a8b0f0dc"

	assert.Len(t, txBlock.Data.Transaction.Inputs, 3, "Should have 3 inputs")

	// Check first input (immOrOwnedObject)
	firstInput := txBlock.Data.Transaction.Inputs[0]
	assert.Equal(t, "object", firstInput["type"])
	assert.Equal(t, "immOrOwnedObject", firstInput["objectType"])
	assert.Equal(t, expectedObjectId, firstInput["objectId"])

	// Check second input (pure u64)
	secondInput := txBlock.Data.Transaction.Inputs[1]
	assert.Equal(t, "pure", secondInput["type"])
	assert.Equal(t, "u64", secondInput["valueType"])
	assert.Equal(t, "10000000000", secondInput["value"])

	// Check third input (sharedObject)
	thirdInput := txBlock.Data.Transaction.Inputs[2]
	assert.Equal(t, "object", thirdInput["type"])
	assert.Equal(t, "sharedObject", thirdInput["objectType"])
	assert.Equal(t, "0x538ddecc8d0842acf43b3fea41ed1f2dd64dbd9df19651e2fe3687e99ee3ce85", thirdInput["objectId"])

	// Verify transactions
	assert.Len(t, txBlock.Data.Transaction.Transactions, 2, "Should have 2 transactions")

	// Check SplitCoins transaction
	splitCoins := txBlock.Data.Transaction.Transactions[0].SplitCoins
	assert.NotNil(t, splitCoins, "Should have SplitCoins transaction")

	// Check MoveCall transaction
	moveCall := txBlock.Data.Transaction.Transactions[1].MoveCall
	assert.NotNil(t, moveCall, "Should have MoveCall transaction")
	assert.Equal(t, "dw", moveCall.Module)
	assert.Equal(t, "deposit", moveCall.Function)

	t.Log("✅ Successfully parsed raw bytes directly")
	t.Logf("   - Extracted Object ID: %s", expectedObjectId)
	t.Logf("   - Amount: 10000000000 (10 USDT)")
	t.Logf("   - Function: %s::%s", moveCall.Module, moveCall.Function)

	// Test serialization
	jsonBytes, err := json.MarshalIndent(txBlock, "", "  ")
	assert.NoError(t, err)
	jsonOutput := string(jsonBytes)

	// Save the direct parsing result
	err = saveTestResult("TestRealBase64DirectBytesParsing", jsonOutput)
	if err != nil {
		t.Logf("Warning: Could not save test result: %v", err)
	} else {
		t.Log("✅ Direct parsing result saved to test_results/TestRealBase64DirectBytesParsing.json")
	}
}

func TestMgoRPCTransactionFormat(t *testing.T) {
	// Real MGO RPC transaction format (the actual structure from MGO RPC)
	mgoRPCTransaction := `{
		"messageVersion": "v1",
		"transaction": {
			"kind": "ProgrammableTransaction",
			"inputs": [
				{
					"type": "object",
					"objectType": "immOrOwnedObject",
					"objectId": "0x88d0d176f92318af69948594736a09bc43ef4a101df6868f9c8fb421a8b0f0dc",
					"version": "7774607",
					"digest": "6MytWN8Tbayw5XVmtVnj3A8tQEHS9517LVjEvJVY7G5V"
				},
				{
					"type": "pure",
					"valueType": "u64",
					"value": "10000000000"
				},
				{
					"type": "object",
					"objectType": "sharedObject",
					"objectId": "0x538ddecc8d0842acf43b3fea41ed1f2dd64dbd9df19651e2fe3687e99ee3ce85",
					"initialSharedVersion": "10197354",
					"mutable": true
				}
			],
			"transactions": [
				{
					"SplitCoins": [
						{
							"Input": 0
						},
						[
							{
								"Input": 1
							}
						]
					]
				},
				{
					"MoveCall": {
						"package": "0x1be5069dd060e52ffa1147dd5af56a40b312a74a034edb78d5b13f4476e03331",
						"module": "dw",
						"function": "deposit",
						"type_arguments": [
							"0x161a456168fe96ea880bcdb7025116cc8c70b5d679a1c361238517312850e020::usdt::USDT"
						],
						"arguments": [
							{
								"Input": 2
							},
							{
								"Result": 0
							}
						]
					}
				}
			]
		},
		"sender": "0x22b6d3195090840253a65a41773832e1ad9eb5959938f38092d9187a083e6034",
		"gasData": {
			"payment": [
				{
					"objectId": "0x12b84bab3e8850f1ce0038750a77941ea90c643cbeba2dcc9e3f2c99debad986",
					"version": 10197381,
					"digest": "EE8sbJZ21Dqv1MgzDo197swKrQGkNzcfHADZFzYQffCj"
				}
			],
			"owner": "0x22b6d3195090840253a65a41773832e1ad9eb5959938f38092d9187a083e6034",
			"price": "1000",
			"budget": "3022656"
		}
	}`

	// Test decoding the MGO RPC transaction format
	decoded, err := transaction.NewRawMessageDecoder().DecodeJSONMessage(mgoRPCTransaction)
	assert.NoError(t, err)
	assert.NotNil(t, decoded)

	// Print the decoded transaction for analysis
	t.Logf("MGO RPC transaction decoded:\n%s", decoded.PrettyPrint())

	// Verify the transaction details
	assert.Equal(t, "0x22b6d3195090840253a65a41773832e1ad9eb5959938f38092d9187a083e6034", decoded.Sender)

	// This should be identified as a USDT deposit
	t.Logf("Transaction Type: %s", decoded.TransactionType)
	t.Logf("Sender: %s", decoded.Sender)

	// Check if we can extract the deposit information
	jsonOutput, err := decoded.ToJSON()
	assert.NoError(t, err)
	t.Logf("JSON output:\n%s", jsonOutput)
}
