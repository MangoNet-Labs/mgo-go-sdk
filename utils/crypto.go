package utils

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/crypto/sha3"
)

func Keccak256(input []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(input)
	return hash.Sum(nil)
}

//------------------test utils ------------------

func EncodeBase64(value []byte) string {
	return base64.StdEncoding.EncodeToString(value)
}

func DecodeBase64(value string) []byte {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil
	}
	return data
}

func HexStringToByteArray(hexString string) ([]byte, error) {
	return hex.DecodeString(hexString)
}

func ByteArrayToHexString(byteArray []byte) string {
	return hex.EncodeToString(byteArray)
}

func ByteArrayToBase64String(byteArray []byte) string {
	return base64.StdEncoding.EncodeToString(byteArray)
}

func Base64StringToByteArray(base64String string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64String)
}
func JsonPrint(v any) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("JSON encoding failed: %s", err)
	}
	fmt.Println(string(jsonData))
}
