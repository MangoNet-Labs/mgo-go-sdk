package ptb

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/transaction"
)

var ctx = context.Background()
var devCli = client.NewMgoClient(config.RpcMgoDevnetEndpoint)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}
}

func TestSimpleTransaction(t *testing.T) {
	key, err := keypair.NewKeypairWithMgoPrivateKey(os.Getenv("KEY_ONE"))
	if err != nil {
		return
	}

	receiver := "0x1313b921d90b9b70aa73132d7ae69a5f5914acece6a5bacb16ce082b4ce116e4"
	gasCoinObjectId := "0xdc9b8d1b0a44e0eda3e77ddc16470616584dff25ca971c073defac8c67bc1804"

	gasCoinObj, err := devCli.MgoGetObject(ctx, request.MgoGetObjectRequest{ObjectId: gasCoinObjectId})
	if err != nil {
		panic(err)
	}
	gasCoin, err := transaction.NewMgoObjectRef(
		model.MgoAddress(gasCoinObjectId),
		gasCoinObj.Data.Version,
		model.ObjectDigest(gasCoinObj.Data.Digest),
	)
	if err != nil {
		panic(err)
	}

	tx := transaction.NewTransaction()

	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasPayment([]transaction.MgoObjectRef{*gasCoin}).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	splitCoin := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.01)),
	})
	tx.TransferObjects([]transaction.Argument{splitCoin}, tx.Pure(receiver))

	resp, err := tx.Execute(
		ctx,
		request.MgoTransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
			ShowEvents:   true,
		},
		"WaitForLocalExecution",
	)
	if err != nil {
		panic(err)
	}

	t.Log(resp.Digest, resp.Effects)
}

func TestMoveCallTransaction(t *testing.T) {
	key, err := keypair.NewKeypairWithMgoPrivateKey(os.Getenv("KEY_ONE"))
	if err != nil {
		return
	}
	gasCoinObjectId := "0xdc9b8d1b0a44e0eda3e77ddc16470616584dff25ca971c073defac8c67bc1804"

	gasCoinObj, err := devCli.MgoGetObject(ctx, request.MgoGetObjectRequest{ObjectId: gasCoinObjectId})
	if err != nil {
		panic(err)
	}
	gasCoin, err := transaction.NewMgoObjectRef(
		model.MgoAddress(gasCoinObjectId),
		gasCoinObj.Data.Version,
		model.ObjectDigest(gasCoinObj.Data.Digest),
	)
	if err != nil {
		panic(err)
	}

	tx := transaction.NewTransaction()

	tx.SetMgoClient(devCli).
		SetSigner(key).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasPayment([]transaction.MgoObjectRef{*gasCoin}).
		SetGasOwner(model.MgoAddress(key.MgoAddress()))

	addressBytes, err := transaction.ConvertMgoAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		panic(err)
	}

	tx.MoveCall(
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"pay",
		"split",
		[]transaction.TypeTag{
			{
				Struct: &transaction.StructTag{
					Address: *addressBytes,
					Module:  "mgo",
					Name:    "MGO",
				},
			},
		},
		[]transaction.Argument{
			tx.Gas(),
			tx.Pure(uint64(1000000000 * 0.01)),
		},
	)

	resp, err := tx.Execute(
		ctx,
		request.MgoTransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
			ShowEvents:   true,
		},
		"WaitForLocalExecution",
	)
	if err != nil {
		panic(err)
	}

	t.Log(resp.Digest, resp.Effects)
}

func TestSponsoredTransaction(t *testing.T) {

	key, err := keypair.NewKeypairWithMgoPrivateKey(os.Getenv("KEY_ONE"))
	if err != nil {
		return
	}
	sponsoredkey, err := keypair.NewKeypairWithMgoPrivateKey(os.Getenv("KEY_TWO"))
	if err != nil {
		return
	}

	receiver := sponsoredkey.MgoAddress()
	transferCoinObjectId := "0xc8236b82442db40953f32d0c301ec20718b4552efa4df93584be8ab7ecd3fd76"
	sponsoredGasCoinObjectId := "0x0012efd22131c65179be63887600aa82171aa172a5fe29a11a145dc23059f6b3"

	// Raw transaction
	tx := transaction.NewTransaction().SetMgoClient(devCli)

	obj, err := devCli.MgoGetObject(ctx, request.MgoGetObjectRequest{ObjectId: transferCoinObjectId})
	if err != nil {
		panic(err)
	}
	ref, err := transaction.NewMgoObjectRef(
		model.MgoAddress(obj.Data.ObjectId),
		obj.Data.Version,
		model.ObjectDigest(obj.Data.Digest),
	)
	if err != nil {
		panic(err)
	}

	tx.TransferObjects(
		[]transaction.Argument{
			tx.Object(
				transaction.CallArg{
					Object: &transaction.ObjectArg{
						ImmOrOwnedObject: ref,
					},
				},
			)},
		tx.Pure(receiver),
	)

	// Sponsored transaction
	newTx, err := tx.NewTransactionFromKind()
	if err != nil {
		panic(err)
	}
	newTx.SetMgoClient(devCli)

	gasCoinObj, err := devCli.MgoGetObject(ctx, request.MgoGetObjectRequest{ObjectId: sponsoredGasCoinObjectId})
	if err != nil {
		panic(err)
	}
	gasCoin, err := transaction.NewMgoObjectRef(
		model.MgoAddress(sponsoredGasCoinObjectId),
		gasCoinObj.Data.Version,
		model.ObjectDigest(gasCoinObj.Data.Digest),
	)
	if err != nil {
		panic(err)
	}

	newTx.SetSigner(key).
		SetSponsoredSigner(sponsoredkey).
		SetSender(model.MgoAddress(key.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasPayment([]transaction.MgoObjectRef{*gasCoin}).
		SetGasOwner(model.MgoAddress(sponsoredkey.MgoAddress()))

	resp, err := newTx.Execute(
		ctx,
		request.MgoTransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
			ShowEvents:   true,
		},
		"WaitForLocalExecution",
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Digest, resp.Effects)
}
