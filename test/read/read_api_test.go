package read

import (
	"context"
	"testing"

	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/utils"
)

var ctx = context.Background()
var devCli = client.NewMgoClient(config.RpcMgoTestnetEndpoint)

func TestGetChainIdentifier(t *testing.T) {
	identifier, err := devCli.MgoGetChainIdentifier(ctx)
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(identifier)
}

func TestGetCheckPoint(t *testing.T) {
	checkpoint, err := devCli.MgoGetCheckpoint(ctx, request.MgoGetCheckpointRequest{
		CheckpointID: "GJjR94qoW5fbDN4pJJMDBLiAZnyYKAKoDhzyYYBW8GPJ",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(checkpoint)
}

func TestGetCheckPoints(t *testing.T) {
	checkpoints, err := devCli.MgoGetCheckpoints(ctx, request.MgoGetCheckpointsRequest{
		Cursor:          "18847201",
		Limit:           3,
		DescendingOrder: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(checkpoints)
}

func TestGetLatestCheckpointSequenceNumber(t *testing.T) {
	number, err := devCli.MgoGetLatestCheckpointSequenceNumber(ctx)
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(number)
}

func TestGetLoadedChildObjects(t *testing.T) {
	ob, err := devCli.MgoGetLoadedChildObjects(ctx, request.MgoGetLoadedChildObjectsRequest{
		Digest: "GxRMedP6tgLbwMY3fpUDgXMaRqZgw9ELhghFNVZSXLaU",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(ob)
}

func TestGetObject(t *testing.T) {
	ob, err := devCli.MgoGetObject(ctx, request.MgoGetObjectRequest{
		ObjectId: "0x11ac113ffd2befec14988aa242635b3a59e2675bf11d95c07d055513bcbf6484",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(ob)
}

func TestGetProtocolConfig(t *testing.T) {
	config, err := devCli.MgoGetProtocolConfig(ctx, request.MgoGetProtocolConfigRequest{
		Version: "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(config)
}

func TestGetTotalTransactionBlocks(t *testing.T) {
	blocks, err := devCli.MgoGetTotalTransactionBlocks(ctx)
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(blocks)
}

func TestGetTransactionBlock(t *testing.T) {
	block, err := devCli.MgoGetTransactionBlock(ctx, request.MgoGetTransactionBlockRequest{
		Digest: "Ed7ZdSUJUbKZDnSQ1uGuEhZ85sR9Mh9xoAK89CyLh8CB",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(block)
}

func TestMultiGetObjects(t *testing.T) {
	obs, err := devCli.MgoMultiGetObjects(ctx, request.MgoMultiGetObjectsRequest{
		ObjectIds: []string{"0x11ac113ffd2befec14988aa242635b3a59e2675bf11d95c07d055513bcbf6484", "0x229b6eb9bf8c0cf365da2d05dbbf8b1cea40168f8fc18c6f4356de9bc21da253"},
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(obs)
}

func TestMutiGetTransactionBlocks(t *testing.T) {
	block, err := devCli.MgoMultiGetTransactionBlocks(ctx, request.MgoMultiGetTransactionBlocksRequest{
		Digests: []string{"Ed7ZdSUJUbKZDnSQ1uGuEhZ85sR9Mh9xoAK89CyLh8CB", "CFTxHp2M7JumzzpUPkXfXthsFSxLmWvYZjGPMQ4fUUEU"},
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(block)
}

func TestTryGetPastObject(t *testing.T) {
	object, err := devCli.MgoTryGetPastObject(ctx, request.MgoTryGetPastObjectRequest{
		ObjectId: "0x229b6eb9bf8c0cf365da2d05dbbf8b1cea40168f8fc18c6f4356de9bc21da253",
		Options: request.MgoObjectDataOptions{
			ShowType:                true,
			ShowOwner:               true,
			ShowPreviousTransaction: true,
			ShowDisplay:             false,
			ShowContent:             true,
			ShowBcs:                 false,
			ShowStorageRebate:       true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(object)
}

func TestTryMultiGetPastObjects(t *testing.T) {
	pastObject1 := request.PastObject{
		ObjectId: "0x11ac113ffd2befec14988aa242635b3a59e2675bf11d95c07d055513bcbf6484",
		Version:  "10251041",
	}
	pastObject2 := request.PastObject{
		ObjectId: "0x229b6eb9bf8c0cf365da2d05dbbf8b1cea40168f8fc18c6f4356de9bc21da253",
		Version:  "106811",
	}
	object, err := devCli.MgoTryMultiGetPastObjects(ctx, request.MgoTryMultiGetPastObjectsRequest{
		MultiGetPastObjects: []*request.PastObject{
			&pastObject1,
			&pastObject2,
		},
		Options: request.MgoObjectDataOptions{
			ShowType:                true,
			ShowOwner:               true,
			ShowPreviousTransaction: true,
			ShowDisplay:             false,
			ShowContent:             true,
			ShowBcs:                 false,
			ShowStorageRebate:       true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(object)
}
