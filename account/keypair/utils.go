package keypair

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/MangoNet-Labs/mgo-go-sdk/utils"

	"github.com/MangoNet-Labs/mgo-go-sdk/config"

	"github.com/btcsuite/btcd/btcutil/bech32"
)

// DecodeMgoPrivateKey decodes a base32-encoded extended private key and returns a
// *model.ParsedKeypair. The input string should be a bech32-encoded string
// starting with the prefix "mgoprivkey".
func DecodeMgoPrivateKey(key string) (scheme config.Scheme, privateKey string, err error) {
	prefix, words, err := bech32.Decode(key)
	if err != nil {
		return
	}
	if prefix != config.MGO_PRIVATE_KEY_PREFIX {
		err = errors.New("invalid private key prefix")
		return
	}
	extendedSecretKey, err := bech32.ConvertBits(words, 5, 8, false)
	if err != nil {
		return
	}
	secretKey := extendedSecretKey[1:]
	_, exists := config.SIGNATURE_FLAG_TO_SCHEME[config.Scheme(extendedSecretKey[0])]
	if !exists {
		err = errors.New("invalid signature scheme flag")
		return
	}
	return config.Scheme(extendedSecretKey[0]), hex.EncodeToString(secretKey), nil
}

// EncodeMgoPrivateKey encodes a private key into a bech32-encoded string with the prefix "mgoprivkey".
// The input value should be a byte slice of length equal to PRIVATE_KEY_SIZE, and the scheme should
// be a valid signature scheme. The encoded string is formed by appending the scheme as the first byte,
// converting the resulting byte slice to a base32 format, and encoding it using bech32.
func EncodeMgoPrivateKey(scheme config.Scheme, privatekey string) (string, error) {
	if strings.HasPrefix(privatekey, "0x") || strings.HasPrefix(privatekey, "0X") {
		privatekey = privatekey[2:]
	}
	key, err := hex.DecodeString(privatekey)
	if err != nil {
		return "", err
	}
	if len(key) != config.PRIVATE_KEY_SIZE {
		return "", errors.New("invalid bytes length")
	}
	privKeyBytes := append([]byte{byte(scheme)}, key...)
	words, err := bech32.ConvertBits(privKeyBytes, 8, 5, true)
	if err != nil {
		return "", err
	}
	return bech32.Encode(config.MGO_PRIVATE_KEY_PREFIX, words)
}

func DecodeBase64WithFlag(key string) (scheme config.Scheme, privateKey string, err error) {
	extendedSecretKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return
	}
	secretKey := extendedSecretKey[1:]
	_, exists := config.SIGNATURE_FLAG_TO_SCHEME[config.Scheme(extendedSecretKey[0])]
	if !exists {
		err = errors.New("invalid signature scheme flag")
		return
	}

	return config.Scheme(extendedSecretKey[0]), hex.EncodeToString(secretKey), nil
}

func EncodeBase64WithFlag(scheme config.Scheme, privateKey string) (string, error) {
	if len(privateKey) != config.PRIVATE_KEY_SIZE {
		return "", errors.New("invalid bytes length")
	}
	privKeyBytes := append([]byte{byte(scheme)}, []byte(privateKey)...)
	return base64.StdEncoding.EncodeToString(privKeyBytes), nil
}

// PublicKeyToMgoAddress takes a public key and a signature scheme, and returns the corresponding MGO address
// as a hexadecimal string prefixed with "0x". It is derived from the public key of the signer by taking the
// first 64 characters of the Keccak-256 hash of the flag byte and the public key. If the scheme is not Ed25519,
// it returns an error.
func PublicKeyToMgoAddress(publicKey []byte, schema config.Scheme) (string, error) {
	if schema != config.Ed25519Flag {
		return "", errors.New("invalid signature scheme flag")
	}
	inputBytes := append([]byte{byte(config.Ed25519Flag)}, publicKey...)
	return "0x" + hex.EncodeToString(utils.Keccak256(inputBytes))[:config.MGO_ADDRESS_LENGTH], nil
}
