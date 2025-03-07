package keypair

import (
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/MangoNet-Labs/mgo-go-sdk/account/signer"
	"github.com/MangoNet-Labs/mgo-go-sdk/account/signer/ed25519"
	"github.com/MangoNet-Labs/mgo-go-sdk/account/signer/secp256k1"
	"github.com/MangoNet-Labs/mgo-go-sdk/bcs"
	"github.com/MangoNet-Labs/mgo-go-sdk/config"
	"github.com/MangoNet-Labs/mgo-go-sdk/model"
	"github.com/MangoNet-Labs/mgo-go-sdk/utils"

	"github.com/tyler-smith/go-bip39"
)

type Keypair struct {
	signer.Signer
	Scheme config.Scheme
}

// NewKeypair creates a new Keypair given a signature scheme flag. The function
// returns a Keypair instance initialized with the specified scheme and a
// randomly generated private key, or an error if the generation process fails.
func NewKeypair(scheme config.Scheme) (*Keypair, error) {
	_, exists := config.SIGNATURE_FLAG_TO_SCHEME[scheme]
	if !exists {
		return nil, errors.New("invalid signature scheme flag")
	}
	switch scheme {
	case config.Secp256k1Flag:
		signer, err := secp256k1.NewSigner()
		if err != nil {
			return nil, err
		}
		return &Keypair{Scheme: scheme, Signer: signer}, nil
	case config.Ed25519Flag:
		signer, err := ed25519.NewSigner()
		if err != nil {
			return nil, err
		}
		return &Keypair{Scheme: scheme, Signer: signer}, nil
	default:
		return nil, errors.New("invalid signature scheme flag")
	}
}

// NewKeypairWithPrivateKey creates a new Keypair from a private key string,
// given a signature scheme flag. The private key string should be a
// hexadecimal string of the raw private key bytes. The function returns a
// Keypair instance initialized with the specified scheme and private key,
// or an error if the decoding process fails.
func NewKeypairWithPrivateKey(scheme config.Scheme, privateKey string) (*Keypair, error) {
	switch scheme {
	case config.Secp256k1Flag:
		signer, err := secp256k1.NewSignerByHex(privateKey)
		if err != nil {
			return nil, err
		}
		return &Keypair{Scheme: scheme, Signer: signer}, nil
	case config.Ed25519Flag:
		signer, err := ed25519.NewSignerByHex(privateKey)
		if err != nil {
			return nil, err
		}
		return &Keypair{Scheme: scheme, Signer: signer}, nil
	default:
		return nil, errors.New("invalid signature scheme flag")
	}
}

// NewKeypairWithMgoPrivateKey creates a new Keypair from a Mango (MGO) private key.
// The MGO private key is decoded to extract the signature scheme and the raw private key.
// It returns a Keypair instance initialized with the specified scheme and private key,
// or an error if the decoding process fails.

func NewKeypairWithMgoPrivateKey(mgoPrivateKey string) (*Keypair, error) {
	scheme, privateKey, err := DecodeMgoPrivateKey(mgoPrivateKey)
	if err != nil {
		return nil, err
	}
	return NewKeypairWithPrivateKey(scheme, privateKey)
}

// NewKeypairWithMnemonic returns a new Keypair derived from the given mnemonic and
// signature scheme. The mnemonic must be a valid BIP39 mnemonic string, and the
// signature scheme must be a valid signature scheme flag. The derivation path
// used is the one specified in the config.DERIVATION_PATH map.
func NewKeypairWithMnemonic(mnemonic string, keytype config.Scheme) (*Keypair, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	derivation, ok := config.DERIVATION_PATH[keytype]
	if !ok {
		return nil, errors.New("invalid signature scheme flag")
	}
	var privateKey []byte
	switch keytype {
	case config.Secp256k1Flag:
		key, err := secp256k1.DeriveForPath(derivation, seed)
		if err != nil {
			return nil, err
		}
		privateKey = key.Key
	case config.Ed25519Flag:
		key, err := ed25519.DeriveForPath(derivation, seed)
		if err != nil {
			return nil, err
		}
		privateKey = key.Key
	default:
		return nil, errors.New("invalid signature scheme flag")
	}

	signer, err := NewKeypairWithPrivateKey(keytype, hex.EncodeToString(privateKey))
	return signer, err
}

// GetMgoAddress returns the MGO address of the Keypair instance as a hexadecimal string prefixed with "0x".
// It is derived from the public key of the signer by taking the first 64 characters of the Keccak-256 hash
// of the flag byte and the public key.
func (k *Keypair) MgoAddress() string {
	inputBytes := append([]byte{byte(k.Scheme)}, k.PublicKeyBytes()...)
	return "0x" + hex.EncodeToString(utils.Keccak256(inputBytes))[:config.MGO_ADDRESS_LENGTH]
}

// MgoPrivateKey returns the MGO private key of the Keypair instance as a bech32-encoded string with the prefix "mgoprivkey".
// The private key is derived from the raw private key of the signer by appending the scheme as the first byte,
// converting the resulting byte slice to a base32 format, and encoding it using bech32.
func (k *Keypair) MgoPrivateKey() string {
	mgoPrivateKey, _ := EncodeMgoPrivateKey(k.Scheme, k.PrivateKeyHex())
	return mgoPrivateKey
}

type SignedTransactionSerializedSig struct {
	TxBytes   string `json:"tx_bytes"`
	Signature string `json:"signature"`
}

// SignPersonalMessage signs a personal message using the Keypair's private key.
// The method prefixes the message with its length using ULEB encoding, hashes
// the result with a personal message intent, and then signs the digest.
// It returns the signature concatenated with the public key and a signature
// scheme flag byte.
func (k *Keypair) SignPersonalMessage(message []byte) []byte {
	message = append(bcs.ULEBEncode(uint64(len(message))), message...)
	data := dataWithIntent(message, config.PersonalMessage)
	digest := digestData(data)
	sigBytes := k.Sign(digest[:])
	publicKey := k.PublicKeyBytes()
	signData := append(sigBytes, publicKey...)
	signData = append([]byte{byte(config.Ed25519Flag)}, signData...)
	return signData
}

// SignTransactionBlock signs a transaction block with the Keypair's private key.
// The method hashes the transaction block with a transaction intent, signs the digest,
// and then returns the signature concatenated with the public key and a signature
// scheme flag byte.
func (k *Keypair) SignTransactionBlock(txn *model.TxnMetaData) *SignedTransactionSerializedSig {
	txBytes, _ := base64.StdEncoding.DecodeString(txn.TxBytes)
	data := dataWithIntent(txBytes, config.TransactionData)
	digest := digestData(data)

	sigBytes := k.Sign(digest[:])

	return &SignedTransactionSerializedSig{
		TxBytes:   txn.TxBytes,
		Signature: k.toSerializedSignature(sigBytes),
	}
}

// dataWithIntent adds a header to the given data that marks it with the given intent.
// The header consists of the intent byte followed by two zero bytes, and the
// result is a new byte slice that contains the header followed by the original
// data.
func dataWithIntent(data []byte, intent config.Signtype) []byte {
	header := []byte{byte(intent), 0, 0}
	markData := make([]byte, len(header)+len(data))
	copy(markData, header)
	copy(markData[len(header):], data)
	return markData
}

// digestData takes a byte slice of data and returns its Keccak-256 hash.
func digestData(data []byte) []byte {
	return utils.Keccak256(data)
}

// toSerializedSignature converts a given signature into a serialized format.
// The serialized signature is composed of a scheme byte, followed by the
// signature bytes, and then the public key bytes. The entire serialized
// signature is then base64 encoded and returned as a string. This format is
// used to ensure the signature, scheme, and public key can be easily
// transmitted and stored as a single string.
func (k *Keypair) toSerializedSignature(signature []byte) string {
	signatureLen := len(signature)
	pubKeyLen := len(k.PublicKeyBytes())
	serializedSignature := make([]byte, 1+signatureLen+pubKeyLen)
	serializedSignature[0] = byte(k.Scheme)
	copy(serializedSignature[1:], signature)
	copy(serializedSignature[1+signatureLen:], k.PublicKeyBytes())
	return base64.StdEncoding.EncodeToString(serializedSignature)
}

type SignatureInfo struct {
	SerializedSignature []byte
	SignatureScheme     string
	Signature           []byte
	PublicKey           []byte
	Bytes               []byte
}

// ExtractSignerMgoAddress takes a serialized signature and returns the MGO address of the signer as a hexadecimal string prefixed with "0x".
// The method extracts the signature scheme from the first byte of the signature, parses the signature into a SignatureInfo object,
// and then derives the MGO address from the public key of the signer by taking the first 64 characters of the Keccak-256 hash
// of the flag byte and the public key. If the signature is invalid, the method returns an error.
func ExtractSignerMgoAddress(sig []byte) (string, error) {
	schema := config.Scheme(sig[0])
	signatureScheme := GetSignatureScheme(sig)
	if sig == nil || signatureScheme == "" {
		return "", errors.New("invalid signature")
	}
	signatureInfo := ParseSignatureInfo(sig, signatureScheme)
	inputBytes := append([]byte{byte(schema)}, signatureInfo.PublicKey...)
	return "0x" + hex.EncodeToString(utils.Keccak256(inputBytes))[:config.MGO_ADDRESS_LENGTH], nil
}

// VerifyPersonalMessage verifies that the given signature was signed with the given public key over the given message.
// The method hashes the message with a personal message intent, verifies the signature, and returns true if the signature
// is valid, false otherwise.
func VerifyPersonalMessage(msg []byte, sig []byte) bool {
	signatureScheme := GetSignatureScheme(sig)
	if sig == nil || signatureScheme == "" {
		return false
	}
	signatureInfo := ParseSignatureInfo(sig, signatureScheme)
	publickey := signatureInfo.PublicKey
	msgReserialize := append(bcs.ULEBEncode(uint64(len(msg))), msg...)
	intentMessage := dataWithIntent(msgReserialize, config.PersonalMessage)
	digest := digestData(intentMessage)

	switch config.SIGNATURE_SCHEME_TO_FLAG[signatureScheme] {
	case config.Ed25519Flag:
		return ed25519.Verify(publickey, digest, signatureInfo.Signature)
	case config.Secp256k1Flag:
		return secp256k1.Verify(publickey, digest, signatureInfo.Signature)
	default:
		return false
	}

}

// VerifyTransactionBlock verifies that the given signature was signed with the given public key over the given transaction block.
// The method hashes the transaction block with a transaction intent, verifies the signature, and returns true if the signature
// is valid, false otherwise.
func VerifyTransactionBlock(txn []byte, sig []byte) bool {
	signatureScheme := GetSignatureScheme(sig)
	if sig == nil || signatureScheme == "" {
		return false
	}
	signatureInfo := ParseSignatureInfo(sig, signatureScheme)
	publickey := signatureInfo.PublicKey

	intentMessage := dataWithIntent(txn, config.TransactionData)
	digest := digestData(intentMessage)

	switch config.SIGNATURE_SCHEME_TO_FLAG[signatureScheme] {
	case config.Ed25519Flag:
		return ed25519.Verify(publickey, digest, signatureInfo.Signature)
	case config.Secp256k1Flag:
		return secp256k1.Verify(publickey, digest, signatureInfo.Signature)
	default:
		return false
	}
}

// GetSignatureScheme returns the signature scheme corresponding to the given byte slice.
// The scheme is determined from the first byte of the given byte slice, and is used to
// determine the size of the signature and public key in the given byte slice.
func GetSignatureScheme(bytes []byte) string {
	return config.SIGNATURE_FLAG_TO_SCHEME[config.Scheme(bytes[0])]
}

// ParseSignatureInfo takes a byte slice representing a serialized signature and a string representing the signature scheme
// as input, and returns a SignatureInfo object containing the signature, public key, and scheme.
// The method slices the input byte slice to extract the signature and public key, and returns a SignatureInfo object
// containing the serialized signature, the signature, the public key, and the scheme.
func ParseSignatureInfo(bytes []byte, signatureScheme string) *SignatureInfo {
	size := config.SIGNATURE_SCHEME_TO_SIZE[signatureScheme]
	signature := bytes[1 : len(bytes)-size]
	publicKey := bytes[1+len(signature):]

	return &SignatureInfo{
		SerializedSignature: bytes,
		SignatureScheme:     signatureScheme,
		Signature:           signature,
		PublicKey:           publicKey,
		Bytes:               bytes,
	}
}
