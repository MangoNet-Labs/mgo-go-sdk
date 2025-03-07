package secp256k1

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

type SignerSecp256k1 struct {
	PrivateKey secp256k1.PrivateKey
	PublicKey  secp256k1.PublicKey
}

func NewSigner() (*SignerSecp256k1, error) {
	privateKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	return newSigner(privateKey)
}

func NewSignerByHex(privatekey string) (signerSecp256k1 *SignerSecp256k1, err error) {
	if strings.HasPrefix(privatekey, "0x") || strings.HasPrefix(privatekey, "0X") {
		privatekey = privatekey[2:]
	}
	seed, err := hex.DecodeString(privatekey)
	if err != nil {
		return nil, err
	}
	return NewSignerBySeed(seed)
}

func NewSignerBySeed(seed []byte) (signerSecp256k1 *SignerSecp256k1, err error) {
	return newSigner(secp256k1.PrivKeyFromBytes(seed))
}

func newSigner(privateKey *secp256k1.PrivateKey) (*SignerSecp256k1, error) {
	publicKey := privateKey.PubKey()
	return &SignerSecp256k1{
		PrivateKey: *privateKey,
		PublicKey:  *publicKey,
	}, nil
}

func (s SignerSecp256k1) String() string {
	return fmt.Sprintf("PrivateKey: %s\nPublicKey: %s\n",
		s.PrivateKeyHex(),
		s.PublicKeyHex())
}

func (s *SignerSecp256k1) PrivateKeyHex() string {
	return "0x" + hex.EncodeToString(s.PrivateKey.Serialize())
}
func (s *SignerSecp256k1) PrivateKeyBytes() []byte {
	return s.PrivateKey.Serialize()
}
func (s *SignerSecp256k1) PublicKeyHex() string {
	return "0x" + hex.EncodeToString(s.PublicKey.SerializeCompressed())
}
func (s *SignerSecp256k1) PublicKeyBytes() []byte {
	return s.PublicKey.SerializeCompressed()
}
func (s *SignerSecp256k1) PublicBase64Key() string {
	return base64.StdEncoding.EncodeToString(s.PublicKeyBytes())
}

func (s *SignerSecp256k1) Sign(message []byte) []byte {
	messageHash := sha256.Sum256(message)
	return parseDERSignature(ecdsa.Sign(&s.PrivateKey, messageHash[:]).Serialize())
}

func Verify(publicKey []byte, message []byte, signature []byte) bool {
	r := new(secp256k1.ModNScalar)
	s := new(secp256k1.ModNScalar)
	if r.SetByteSlice(signature[:32]) {
		return false
	}
	if s.SetByteSlice(signature[32:]) {
		return false
	}
	pk, err := secp256k1.ParsePubKey(publicKey)
	if err != nil {
		fmt.Println(err)
		return false
	}
	messageHash := sha256.Sum256(message)
	return ecdsa.NewSignature(r, s).Verify(messageHash[:], pk)
}

func parseDERSignature(sig []byte) []byte {
	const (
		asn1SequenceID = 0x30
		asn1IntegerID  = 0x02
		minSigLen      = 8
		maxSigLen      = 72
	)
	rLen := int(sig[3])
	rStart := 4
	rEnd := rStart + rLen
	rBytes := sig[rStart:rEnd]
	for len(rBytes) > 0 && rBytes[0] == 0x00 {
		rBytes = rBytes[1:]
	}
	sLen := int(sig[rEnd+1])
	sStart := rEnd + 2
	sEnd := sStart + sLen
	sBytes := sig[sStart:sEnd]
	for len(sBytes) > 0 && sBytes[0] == 0x00 {
		sBytes = sBytes[1:]
	}
	return append(rBytes, sBytes...)
}
