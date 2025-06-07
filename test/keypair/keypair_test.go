package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/utils"
)

var ctx = context.Background()
var devCli = client.NewMgoClient(config.RpcMgoTestnetEndpoint)

func TestSignPersonalMessage(t *testing.T) {

	sig, err := keypair.NewKeypair(config.Ed25519Flag)
	if err != nil {
		panic(err)
	}
	signData := sig.SignPersonalMessage([]byte("hello world"))

	t.Log(signData)
	t.Log(utils.ByteArrayToBase64String(signData))

	result := keypair.VerifyPersonalMessage([]byte("hello world"), signData)
	t.Log(result)

	scheme := keypair.GetSignatureScheme(signData)
	signatureInfo := keypair.ParseSignatureInfo(signData, scheme)

	address, err := keypair.PublicKeyToMgoAddress(signatureInfo.PublicKey, config.Ed25519Flag)

	if sig.MgoAddress() != address {
		t.Fatal("address not match")
	}
}

func TestSignTransactionBlock(t *testing.T) {

	sig, err := keypair.NewKeypairWithPrivateKey(config.Secp256k1Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		panic(err)
	}
	fmt.Println(sig.PrivateKeyHex())
	fmt.Println(sig.MgoAddress())
	pay, err := devCli.Pay(ctx, request.PayRequest{
		Signer:      sig.MgoAddress(),
		MgoObjectId: []string{"0x50ad258e7a6bf92dc13053d997f77a978a407a5f6651702ba8b469a5cc4f2d71"},
		Recipient:   []string{"0xd66993d8fd4657d7b9126178fa208d6c2eae35cd4677604206bd0b14b27189fa"},
		Amount:      []string{"1000000000"},
		Gas:         "0x5d645be778070cc56d59d4bc140538707d660ff492ee24b346506041cc83bd0e",
		GasBudget:   "10000000",
	})
	if err != nil {
		t.Fatal(err)
	}

	utils.JsonPrint(pay)
	data, err := sig.SignTransactionBlock(&pay)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(data.TxBytes)
	t.Log(len(utils.DecodeBase64(data.Signature)))
	t.Log(keypair.VerifyTransactionBlock(utils.DecodeBase64(data.TxBytes), utils.DecodeBase64(data.Signature)))

	executeRes, err := devCli.SignAndExecuteTransactionBlock(ctx, request.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: pay,
		Keypair:     sig,
		Options: request.TransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
		},
		RequestType: "WaitForLocalExecution",
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	t.Log(executeRes)

}

func TestKeypair(t *testing.T) {
	sig, err := keypair.NewKeypair(config.Ed25519Flag)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(sig.PrivateKeyHex())
	t.Log(sig.PublicKeyHex())
	t.Log(sig.MgoPrivateKey())
	t.Log(sig.MgoAddress())
}
