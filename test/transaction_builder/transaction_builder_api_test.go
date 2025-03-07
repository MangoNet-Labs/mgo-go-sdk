package transaction_builder

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
	"github.com/mangonet-labs/mgo-go-sdk/utils"
)

var ctx = context.Background()
var devCli = client.NewMgoClient(config.RpcMgoTestnetEndpoint)

func TestBatchTransaction(t *testing.T) {
	rsp, err := devCli.BatchTransaction(ctx, request.BatchTransactionRequest{
		Signer: "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		RPCTransactionRequestParams: []request.RPCTransactionRequestParams{
			{
				MoveCallRequestParams: &request.MoveCallRequest{
					PackageObjectId: "0x0000000000000000000000000000000000000000000000000000000000000002",
					Module:          "mgo",
					Function:        "transfer",
					TypeArguments:   []interface{}{},
					Arguments: []interface{}{
						"0xe586e913e413de2df45e3eca0f0adf342a1f5d8d71e61805e14fd4872529f727",
						"0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238",
					},
				},
			},
			{
				TransferObjectRequestParams: &request.TransferObjectRequest{
					ObjectId:  "0x11ac113ffd2befec14988aa242635b3a59e2675bf11d95c07d055513bcbf6484",
					Recipient: "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238",
				},
			},
		},
		Gas:                            "0xc7f7956990f4d210024f34d9c5692e0b42d01aa67c94532521e8be886a0dbaa7",
		GasBudget:                      "100000",
		MgoTransactionBlockBuilderMode: "DevInspect",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(rsp)
}

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

func TestPay(t *testing.T) {
	pay, err := devCli.Pay(ctx, request.PayRequest{
		Signer:      "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		MgoObjectId: []string{"0x05678c9529d3354a291fc3235f445dc480ebd476fc281654e4731d2739a5e542", "0x11ac113ffd2befec14988aa242635b3a59e2675bf11d95c07d055513bcbf6484"},
		Recipient:   []string{"0xbb3888e6c078a8ccedde58394873584ba39878984f1f8da4cba870de7eb5c3d2", "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"},
		Amount:      []string{"1000", "1000"},
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
		TxnMetaData: pay,
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
	fmt.Println(executeRes)

}

func TestPayAllMgo(t *testing.T) {
	pay, err := devCli.PayAllMgo(ctx, request.PayAllMgoRequest{
		Signer:      "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		MgoObjectId: []string{"0xe586e913e413de2df45e3eca0f0adf342a1f5d8d71e61805e14fd4872529f727", "0xc7f7956990f4d210024f34d9c5692e0b42d01aa67c94532521e8be886a0dbaa7"},
		Recipient:   "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238",
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
		TxnMetaData: pay,
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
	fmt.Println(executeRes)

}

func TestPayMgo(t *testing.T) {
	pay, err := devCli.PayMgo(ctx, request.PayMgoRequest{
		Signer:      "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		MgoObjectId: []string{"0xc7ead1d1ffd8957deaca7de892d26ef3c7f842db9a32ca5a937eee48bfd9a150", "0xafa150aa24319f7ba0b543b4923497b94aeac15c372883d5219f2f1ba0c1869a"},
		Recipient:   []string{"0xbb3888e6c078a8ccedde58394873584ba39878984f1f8da4cba870de7eb5c3d2", "0x0cafa361487490f306c0b4c3e4cf0dc6fd584c5259ab1d5457d80a9e2170e238"},
		Amount:      []string{"1000", "1000"},
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
		TxnMetaData: pay,
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
	fmt.Println(executeRes)

}

func TestPublish(t *testing.T) {
	publish, err := devCli.Publish(ctx, request.PublishRequest{
		Sender:          "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		CompiledModules: []string{"oRzrCwYAAAALAQAUAhQoAzxrBKcBHAXDAagBB+sC+gII5QVgBsUGJArpBisMlAevAw3DCgYAEwIKAg0CEgIcAh0CHwImAicBKQAEDAAAAAwAAAUDAAADAwABAQQBAAECAgwBAAEFBgIABggEAAgHAgAAGQABAAAOAgEAACoDAQAADwQBAAEQCwEBAAEbGBYBAAElDwsBAAEoGRYBAAErAQsBAAIVEBEBAAIaEQsBAAIoFRYBAAMRDQEBAwQgFxYABh4ABQAHIQ0BAQwHIgkBAQwIIwYHAAkJGgEBABAICAoPDAwOBgoJChASCwoKCgUKBwoEChIUDBsBBwgIAAMHCAEHCAAPBQcIAQcIAAMFBwgICAcIAAsFAQgGDw8PCgIKAgcICAEIBwEGCAgBBQEIAQIJAAUBCAYBCwQBCQABCAABCQABCAMCBwsEAQkAAwILBAEJAAcICAELBQEJAAELBQEIBgcLBAEIBgoCAwMLBAEIBgUPAQIBBgsFAQkAAQMCAwICBwsEAQkACwQBCQABBgsEAQkAAgcKCQAKCQABCAIIQWRtaW5DYXAHQmFsYW5jZQRDb2luDkNyb3NzQ2hhaW5DYWxsB0RlcG9zaXQORGVwb3NpdEFuZENhbGwDTUdPCVR4Q29udGV4dANVSUQGYXBwZW5kB2JhbGFuY2UJY2FsbF9kYXRhCmNhbGxfdmFsdWUEY29pbg5jcm9zc0NoYWluQ2FsbA5kZXBvc2l0QW5kQ2FsbAxkZXN0cm95X3plcm8EZW1pdAVldmVudANldm0MZXZtX3NlcXVlbmNlDGZyb21fYmFsYW5jZQlnYXNfbGltaXQJZ2FzX3ByaWNlAmlkBGluaXQMaW50b19iYWxhbmNlBGpvaW4EbWF0aANtZ28DbmV3Bm9iamVjdANwb3cTcHVibGljX3NoYXJlX29iamVjdA9wdWJsaWNfdHJhbnNmZXIGc2VuZGVyCHNlcXVlbmNlBXNwbGl0CHRyYW5zZmVyCnR4X2NvbnRleHQFdmFsdWUGdmVjdG9yCHdpdGhkcmF3BHplcm8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAQkCARIDCAEAAAAAAAAAAwgCAAAAAAAAAAMIAwAAAAAAAAAAAgQYCAckDxQPCgsEAQgGAQIBGAgHAgIGJA8jBQwPFw8WDwsKAgMCARQPAAAAAAEPCgARDhIBCgAuERE4AAsAEQ5KAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABKAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA4ARIAOAICAQEEAAEVCgIKARAAFEoBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABYhBAkFDQsBAQcEJwoCCwEPABULAhIDOAMCAgEEAAEJCwEPAQsCOAQLBDgFCwM4BgIDAQQAE2IOBUEUBioAAAAAAAAAIQQGBQwLAAELBwEHAicOATgHDAoKAwoEGAoCFgYKAAAAAAAAAAcBBwAXEQ1NGjQMCwoLCwolBCIFKAsAAQsHAQcDJwoHLhERDA0LATgIDAgNCAsLOAQMDAoADwELDDgJAQ4IOAoGAAAAAAAAAAAkBEMLCAsHOAUKDTgGBUcLBwELCDgLQBQAAAAAAAAAAAwJDQkLBTgMDQkLBjgMCgAQAhRKAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWDA4KDgsADwIVCw4LDQsCCwMLBAsJEgI4DQIAAgADAAEA"},
		Dependencies:    []string{"0x0000000000000000000000000000000000000000000000000000000000000002", "0x0000000000000000000000000000000000000000000000000000000000000001"},
		Gas:             "0x70bb8f6c182333aecbc475b97780098fa3ee8de61781f84edfa87c64bfe84ca2",
		GasBudget:       "100000000",
	})
	if err != nil {
		t.Fatal(err)
	}
	ed25519Signer, err := getSigner()
	if err != nil {
		t.Fatal(err)
	}

	executeRes, err := devCli.SignAndExecuteTransactionBlock(ctx, request.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: publish,
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
	fmt.Println(executeRes)

}

func TestRequestAddStake(t *testing.T) {
	publish, err := devCli.RequestAddStake(ctx, request.AddStakeRequest{
		Signer:    "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		Coins:     []string{"0x05678c9529d3354a291fc3235f445dc480ebd476fc281654e4731d2739a5e542"},
		Amount:    "1000000000",
		Gas:       "0x577d7a07b9c8fb2e605264ebf88b932a2a3924de82f2c73f86cd008ce2d59c51",
		GasBudget: "100000000",
		Validator: "0x8520c27a20d69a275cd9cc8877f850b7a3bfe8ee4bd84a8c7749a43d76a8a380",
	})
	if err != nil {
		t.Fatal(err)
	}
	ed25519Signer, err := getSigner()
	if err != nil {
		t.Fatal(err)
	}

	executeRes, err := devCli.SignAndExecuteTransactionBlock(ctx, request.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: publish,
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
	fmt.Println(executeRes)

}

func TestRequestWithdrawStake(t *testing.T) {
	publish, err := devCli.RequestWithdrawStake(ctx, request.WithdrawStakeRequest{
		Signer:         "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		StakedObjectId: "0x70a3040054dede54d0e99be74ca80e22be5cd5710c57a725d55c2c7640b0028b",
		Gas:            "0x822f6705df64d073cbfeb2b2ef088f281aa2d486ea9b5c7fbb0ded58171d7f84",
		GasBudget:      "10000000",
	})
	if err != nil {
		t.Fatal(err)
	}
	ed25519Signer, err := getSigner()
	if err != nil {
		t.Fatal(err)
	}

	executeRes, err := devCli.SignAndExecuteTransactionBlock(ctx, request.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: publish,
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
	fmt.Println(executeRes)

}

func TestMoveCall(t *testing.T) {
	newsinger, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, "0xa1fbf2c281a52d8655a2c793376490bc4f4bef6a1e89346e5d9a255ba4972236")
	if err != nil {
		t.Fail()
	}
	t.Log(newsinger.PublicKeyHex())
	t.Log(newsinger.MgoAddress())
	var gas = "0xdf23fe8f7c87295d8992ab26422a45a11b5d08490c5f07b12602c841c22bb193"
	publish, err := devCli.MoveCall(ctx, request.MoveCallRequest{
		Signer:          newsinger.MgoAddress(),
		PackageObjectId: "0x2",
		Module:          "mgo",
		Function:        "transfer",
		TypeArguments:   []interface{}{},
		Arguments: []interface{}{
			"0xc6fe345768482c9d36229e5b5fd9a3136d102992e94a7833d9817f9fdd36b69f",
			"0xa242d62575ab9de9da98dfe9e1a9d52c5dd239bd7efd89517640fc54ded38c04"},
		Gas:       &gas,
		GasBudget: "100000000",
	})
	fmt.Println(publish, err)
	utils.JsonPrint(publish)

	executeRes, err := devCli.SignAndExecuteTransactionBlock(ctx, request.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: publish,
		Keypair:     newsinger,
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
	fmt.Println(executeRes)

}
