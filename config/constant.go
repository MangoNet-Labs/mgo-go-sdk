package config

import "math"

var (
	MGO_PRIVATE_KEY_PREFIX = "mgoprivkey"
	PRIVATE_KEY_SIZE       = 32
	MGO_ADDRESS_LENGTH     = 64
)

var (
	SIGNATURE_FLAG_TO_SCHEME = map[Scheme]string{
		0x00: "ED25519",
		0x01: "Secp256k1",
	}
	SIGNATURE_SCHEME_TO_FLAG = map[string]Scheme{
		"ED25519":   0x00,
		"Secp256k1": 0x01,
	}
	SIGNATURE_SCHEME_TO_SIZE = map[string]int{
		"ED25519":   32,
		"Secp256k1": 33,
	}

	DERIVATION_PATH = map[Scheme]string{
		0x00: `m/44'/938'/0'/0'/0'`,
		0x01: `m/54'/938'/0'/0/0`,
	}
)

type Scheme byte
type Keytype byte
type Signtype byte

const (
	Ed25519Flag   Scheme = 0
	Secp256k1Flag Scheme = 1
	ErrorFlag     byte   = math.MaxUint8
)

const (
	ErrKey      Keytype = 0
	HexKey      Keytype = 1
	MgoKey      Keytype = 2
	B64Key      Keytype = 3
	GenerateKey Keytype = 99
)
const (
	TransactionData    Signtype = 0
	TransactionEffects Signtype = 1
	CheckpointSummary  Signtype = 2
	PersonalMessage    Signtype = 3
)
const (
	Ed25519PublicKeyLength   = 32
	Secp256k1PublicKeyLength = 33
)

const (
	DefaultAccountAddressLength = 16
	AccountAddress20Length      = 20
	AccountAddress32Length      = 32
)

const (
	RpcMgoTestnetEndpoint       = "https://fullnode.testnet2.mangonetwork.io/"
	RpcMgoMirrorTestnetEndpoint = "https://fullnode.testnet.mangonetwork.io/"
	RpcMgoDevnetEndpoint        = "https://fullnode.devnet.mangonetwork.io/"

	WssMgoTestnetEndpoint       = "wss://fullnode.testnet2.mangonetwork.io"
	WssMgoMirrorTestnetEndpoint = "wss://fullnode.testnet.mangonetwork.io"
	WssMgoDevnetEndpoint        = "wss://fullnode.devnet.mangonetwork.io"

	FaucetTestnetEndpoint       = "https://faucet.testnet2.mangonetwork.io/gas"
	FaucetMirrorTestnetEndpoint = "https://faucet.testnet.mangonetwork.io/gas"
	FaucetDevnetEndpoint        = "https://faucet.devnet.mangonetwork.io/gas"
)
