package write

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/MangoNet-Labs/mgo-go-sdk/account/keypair"
	"github.com/MangoNet-Labs/mgo-go-sdk/client"
	"github.com/MangoNet-Labs/mgo-go-sdk/config"
	"github.com/MangoNet-Labs/mgo-go-sdk/model/request"
	"github.com/MangoNet-Labs/mgo-go-sdk/utils"
)

var ctx = context.Background()
var devCli = client.NewMgoClient(config.RpcMgoTestnetEndpoint)

func getSigner() (*keypair.Keypair, error) {
	bytes, err := os.ReadFile("../../private_keys.json")
	if err != nil {
		return nil, err
	}
	store := []string{}
	err = json.Unmarshal(bytes, &store)
	if err != nil {
		return nil, err
	}

	key, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		return nil, err
	}
	return key, nil
}

func TestDryRunTransactionBlock(t *testing.T) {
	mergeCoins, err := devCli.MergeCoins(ctx, request.MergeCoinsRequest{
		Signer:      "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		PrimaryCoin: "0x991c66f0d0308fbf990fd70cb188b343f5a8721078fd801f5cf33498bb9881ee",
		CoinToMerge: "0x822f6705df64d073cbfeb2b2ef088f281aa2d486ea9b5c7fbb0ded58171d7f84",
		Gas:         "0x7509c59b1e64d881770296ff596b0b442afefb88ea0469efd1cfa354bcca57ee",
		GasBudget:   "10000000",
	})
	if err != nil {
		t.Fatal(err)
	}
	block, err := devCli.MgoDryRunTransactionBlock(ctx, request.MgoDryRunTransactionBlockRequest{
		TxBytes: mergeCoins.TxBytes,
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(block)
}

func TestExecuteTransactionBlock(t *testing.T) {
	mergeCoins, err := devCli.MergeCoins(ctx, request.MergeCoinsRequest{
		Signer:      "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		PrimaryCoin: "0x991c66f0d0308fbf990fd70cb188b343f5a8721078fd801f5cf33498bb9881ee",
		CoinToMerge: "0x822f6705df64d073cbfeb2b2ef088f281aa2d486ea9b5c7fbb0ded58171d7f84",
		Gas:         "0x7509c59b1e64d881770296ff596b0b442afefb88ea0469efd1cfa354bcca57ee",
		GasBudget:   "10000000",
	})
	if err != nil {
		t.Fatal(err)
	}
	ed25519Signer, err := getSigner()
	if err != nil {
		t.Fatal(err)
	}
	signedTxn := ed25519Signer.SignTransactionBlock(&mergeCoins)
	block, err := devCli.MgoExecuteTransactionBlock(ctx, request.MgoExecuteTransactionBlockRequest{
		TxBytes:   mergeCoins.TxBytes,
		Signature: []string{signedTxn.Signature},
		Options: request.MgoTransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
		},
		RequestType: "WaitForLocalExecution",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(block)
}
