package secp256k1

import (
	"testing"

	"github.com/MangoNet-Labs/mgo-go-sdk/account/signer/secp256k1"
)

func TestSecp256k1(t *testing.T) {
	sig, err := secp256k1.NewSignerByHex("0xa11b0a4e1a132305652ee7a8eb7848f6ad5ea381e3ce20a2c086a2e388230811")
	if err != nil {
		panic(err)
	}
	Signature := sig.Sign([]byte("hello world"))
	t.Log(Signature)
	t.Log(secp256k1.Verify(sig.PublicKeyBytes(), []byte("hello world"), Signature))
	t.Log(sig)
	t.Log(sig.PrivateKeyHex())
	t.Log(sig.PublicKeyHex())

	newsig, err := secp256k1.NewSignerByHex(sig.PrivateKeyHex())
	if err != nil {
		panic(err)
	}
	t.Log(newsig)
	t.Log(newsig.PrivateKeyHex())
	t.Log(newsig.PublicKeyHex())

}
