package extend

import (
	"context"
	"fmt"
	"testing"

	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/utils"
)

var ctx = context.Background()
var devCli = client.NewMgoClient(config.RpcMgoTestnetEndpoint)

func TestResolveNameServiceAddress(t *testing.T) {
	address, err := devCli.MgoXResolveNameServiceAddress(ctx, request.MgoXResolveNameServiceAddressRequest{
		Name: "example.mgo",
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	t.Log(address)
}

func TestResolveNameServiceNames(t *testing.T) {
	address, err := devCli.MgoXResolveNameServiceNames(ctx, request.MgoXResolveNameServiceNamesRequest{
		Address: "0x0000000000000000000000000000000000000000000000000000000000000002",
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	t.Log(address)
}

func TestGetDynamicFieldObject(t *testing.T) {
	object, err := devCli.MgoXGetDynamicFieldObject(ctx, request.MgoXGetDynamicFieldObjectRequest{
		ObjectId: "0x11ac113ffd2befec14988aa242635b3a59e2675bf11d95c07d055513bcbf6484",
		DynamicFieldName: request.DynamicFieldObjectName{
			Type:  "0x2::mgp::MGO",
			Value: "",
		},
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	t.Log(object)
}

func TestGetDynamicFields(t *testing.T) {
	object, err := devCli.MgoXGetDynamicFields(ctx, request.MgoXGetDynamicFieldsRequest{
		ObjectId: "0x11ac113ffd2befec14988aa242635b3a59e2675bf11d95c07d055513bcbf6484",
		Cursor:   "0xa9334aeacc435c70ab9635e47a277d8f8dd9d87765d1aadec2db8cc24c312542",
		Limit:    3,
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	t.Log(object)
}

func TestOwnedObjects(t *testing.T) {
	object, err := devCli.MgoXGetOwnedObjects(ctx, request.MgoXGetOwnedObjectsRequest{
		Address: "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		Limit:   10,
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	utils.JsonPrint(object)
}

func TestQueryTransactionBlocks(t *testing.T) {
	transactionBlocks, err := devCli.MgoXQueryTransactionBlocks(ctx, request.MgoXQueryTransactionBlocksRequest{
		//MgoTransactionBlockResponseQuery: request.MgoTransactionBlockResponseQuery{
		//	TransactionFilter: request.TransactionFilter{
		//		"InputObject": "0x93633829fcba6d6e0ccb13d3dbfe7614b81ea76b255e5d435032cd8595f37eb8",
		//	},
		//},
		Cursor:          "5uKUUtqgd7aocMBrPUqu4yXjHhKooVHeNqQ1HyM8e6BC",
		Limit:           10,
		DescendingOrder: false,
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	utils.JsonPrint(transactionBlocks)
}

func TestQueryEventsByMoveEventType(t *testing.T) {
	rsp, err := devCli.MgoXQueryEvents(ctx, request.MgoXQueryEventsRequest{
		MgoEventFilter: request.EventFilterByMoveEventType{
			MoveEventType: "0x0000000000000000000000000000000000000000000000000000000000000003::validator::StakingRequestEvent",
		},
		Limit:           20,
		DescendingOrder: true,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	utils.PrettyPrint(rsp)
}

func TestQueryEventsByTimeRange(t *testing.T) {
	rsp, err := devCli.MgoXQueryEvents(ctx, request.MgoXQueryEventsRequest{
		MgoEventFilter: request.EventFilterByTimeRange{
			request.TimeRange{
				StartTime: "1723798821000",
				EndTime:   "1723800621000",
			},
		},
		Limit:           20,
		DescendingOrder: true,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	utils.PrettyPrint(rsp)
}

func TestQueryEventsByTimeRangeA(t *testing.T) {
	eventType := request.EventFilterByMoveEventType{
		MoveEventType: "0x0000000000000000000000000000000000000000000000000000000000000003::validator::StakingRequestEvent",
	}
	timeRange := request.TimeRange{
		StartTime: "1723798821000",
		EndTime:   "1723800621000",
	}
	m := map[string]interface{}{}
	m["And"] = []interface{}{timeRange, eventType}

	rsp, err := devCli.MgoXQueryEvents(ctx, request.MgoXQueryEventsRequest{
		MgoEventFilter:  m,
		Limit:           20,
		DescendingOrder: true,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	utils.PrettyPrint(rsp)
}
