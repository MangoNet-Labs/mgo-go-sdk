package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
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
func TestMergeCoin(t *testing.T) {
	mergeCoins, err := devCli.MergeCoins(ctx, request.MergeCoinsRequest{
		Signer:      "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		PrimaryCoin: "0x05678c9529d3354a291fc3235f445dc480ebd476fc281654e4731d2739a5e542",
		CoinToMerge: "0xb421a6f124cc4da9d12b4242a24eeb4be7d6e69871f53cf24ffe9deb35f66ccf",
		Gas:         "0x822f6705df64d073cbfeb2b2ef088f281aa2d486ea9b5c7fbb0ded58171d7f84",
		GasBudget:   "10000000",
	})
	if err != nil {
		t.Fatal(err)
	}
	ed25519Signer, err := getSigner()
	if err != nil {
		t.Fatal(err)
	}
	executeRes, err := devCli.SignAndExecuteTransactionBlock(ctx, request.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: mergeCoins,
		Keypair:     ed25519Signer,
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

func TestSplitCoin(t *testing.T) {
	splitCoins, err := devCli.SplitCoin(ctx, request.SplitCoinRequest{
		Signer:       "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		CoinObjectId: "0x91d2925ccb7be261e9db6f23daf9678a38945cb274014e1521b7441fbbc1a18d",
		SplitAmounts: []string{"1000", "1000"},
		Gas:          "0x9e9944e470b44c1363409505ef6d154562572a97cbca88dccfd0d972858b54a5",
		GasBudget:    "10000000",
	})
	if err != nil {
		t.Fatal(err)
	}
	ed25519Signer, err := getSigner()
	if err != nil {
		t.Fatal(err)
	}
	executeRes, err := devCli.SignAndExecuteTransactionBlock(ctx, request.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: splitCoins,
		Keypair:     ed25519Signer,
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

func TestSplitCoinEqual(t *testing.T) {
	splitCoins, err := devCli.SplitCoinEqual(ctx, request.SplitCoinEqualRequest{
		Signer:       "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		CoinObjectId: "0x91d2925ccb7be261e9db6f23daf9678a38945cb274014e1521b7441fbbc1a18d",
		SplitCount:   "3",
		Gas:          "0x9e9944e470b44c1363409505ef6d154562572a97cbca88dccfd0d972858b54a5",
		GasBudget:    "10000000",
	})
	if err != nil {
		t.Fatal(err)
	}
	ed25519Signer, err := getSigner()
	if err != nil {
		t.Fatal(err)
	}
	executeRes, err := devCli.SignAndExecuteTransactionBlock(ctx, request.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: splitCoins,
		Keypair:     ed25519Signer,
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
