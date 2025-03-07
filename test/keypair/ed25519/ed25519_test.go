package main

import (
	"testing"

	"github.com/MangoNet-Labs/mgo-go-sdk/account/signer/ed25519"
)

func TestEd25519(t *testing.T) {

	sig, err := ed25519.NewSigner()
	if err != nil {
		panic(err)
	}
	t.Log(sig.PrivateKeyHex())
	t.Log(sig.PublicKeyHex())

	sig, err = ed25519.NewSignerByHex(sig.PrivateKeyHex())
	if err != nil {
		panic(err)
	}
	t.Log(sig.PrivateKeyHex())
	t.Log(sig.PublicKeyHex())

	data := sig.Sign([]byte("hello world"))

	t.Log(data)
	t.Log(ed25519.Verify(sig.PublicKeyBytes(), []byte("hello world"), data))
}
