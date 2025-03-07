package signer

type Signer interface {
	Sign([]byte) []byte
	PublicKeyHex() string
	PrivateKeyHex() string
	PublicKeyBytes() []byte
	PrivateKeyBytes() []byte
	PublicBase64Key() string
}
