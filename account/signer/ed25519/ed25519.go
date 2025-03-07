package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

type SignerEd25519 struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

func NewSigner() (*SignerEd25519, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return newSigner(privateKey)
}

func NewSignerByHex(privatekey string) (signerEd25519 *SignerEd25519, err error) {
	if strings.HasPrefix(privatekey, "0x") || strings.HasPrefix(privatekey, "0X") {
		privatekey = privatekey[2:]
	}
	seed, err := hex.DecodeString(privatekey)
	if err != nil {
		return nil, err
	}
	return NewSignerBySeed(seed)
}

func NewSignerBySeed(seed []byte) (signerEd25519 *SignerEd25519, err error) {
	defer func() {
		if r := recover(); r != nil {
			signerEd25519 = nil
			err = fmt.Errorf("recovered from panic: %s", r)
		}
	}()
	return newSigner(ed25519.NewKeyFromSeed(seed))
}

func newSigner(privateKey ed25519.PrivateKey) (*SignerEd25519, error) {
	publicKey := privateKey.Public().(ed25519.PublicKey)
	return &SignerEd25519{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

func (s SignerEd25519) String() string {
	return fmt.Sprintf("PrivateKey: %s\nPublicKey: %s\n",
		hex.EncodeToString(s.PrivateKey),
		hex.EncodeToString(s.PublicKey))
}

func (s *SignerEd25519) SecretKeyHex() string {
	return "0x" + hex.EncodeToString(s.PrivateKey)
}

func (s *SignerEd25519) SecretKeyBytes() []byte {
	prv := make([]byte, len(s.PrivateKey))
	copy(prv, s.PrivateKey)
	return prv
}

func (s *SignerEd25519) PrivateKeyHex() string {
	return "0x" + hex.EncodeToString(s.PrivateKey[:32])
}

func (s *SignerEd25519) PrivateKeyBytes() []byte {
	prv := make([]byte, 32)
	copy(prv, s.PrivateKey[:32])
	return prv
}

func (s *SignerEd25519) PublicKeyHex() string {
	return "0x" + hex.EncodeToString(s.PublicKey)
}

func (s *SignerEd25519) PublicKeyBytes() []byte {
	pk := make([]byte, len(s.PublicKey))
	copy(pk, s.PublicKey)
	return pk
}

func (s *SignerEd25519) PublicBase64Key() string {
	return base64.StdEncoding.EncodeToString(s.PublicKey)
}

func (s *SignerEd25519) Sign(message []byte) []byte {
	return ed25519.Sign(s.PrivateKey, message)
}

func Verify(publicKey []byte, message []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, message, signature)
}
