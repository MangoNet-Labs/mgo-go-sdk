package websocket

import (
	"context"
	"testing"

	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/model/response"
	"github.com/mangonet-labs/mgo-go-sdk/utils"
)

var ctx = context.Background()
var devWsCli = client.NewMgoWebsocketClient(config.WssMgoDevnetEndpoint)

func TestSubscribeEventByMoveEventType(t *testing.T) {

	receiveMsgCh := make(chan response.MgoEventResponse, 10)
	err := devWsCli.SubscribeEvent(ctx, request.MgoSubscribeEventsRequest{
		MgoEventFilter: request.EventFilterByMoveEventType{
			MoveEventType: "0x0000000000000000000000000000000000000000000000000000000000000003::validator::StakingRequestEvent",
		},
	}, receiveMsgCh)
	if err != nil {
		panic(err)
	}

	for {
		select {
		// receive Mgo event
		case msg := <-receiveMsgCh:
			utils.PrettyPrint(msg)
		case <-ctx.Done():
			return
		}
	}

}

func TestSubscribeEventBySender(t *testing.T) {
	receiveMsgCh := make(chan response.MgoEventResponse, 10)
	err := devWsCli.SubscribeEvent(ctx, request.MgoSubscribeEventsRequest{
		MgoEventFilter: request.EventFilterBySender{
			Sender: "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		},
	}, receiveMsgCh)
	if err != nil {
		panic(err)
	}

	for {
		select {
		// receive Mgo event
		case msg := <-receiveMsgCh:
			utils.PrettyPrint(msg)
		case <-ctx.Done():
			return
		}
	}

}

func TestSubscribeEventByMoveEventTypeAndSender(t *testing.T) {
	byMoveEventType := request.EventFilterByMoveEventType{
		MoveEventType: "0x5ff2c7fb02e5eb9ed9175f47b2f7b6ea07099e43c7ef5d85b51de4f5372f38ef::amm_swap::SwapEvent",
	}
	bySender := request.EventFilterBySender{
		Sender: "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
	}
	eventFilterList := []interface{}{bySender, byMoveEventType}
	eventFilter := map[string]interface{}{}
	eventFilter["And"] = eventFilterList
	receiveMsgCh := make(chan response.MgoEventResponse, 10)
	err := devWsCli.SubscribeEvent(ctx, request.MgoSubscribeEventsRequest{
		MgoEventFilter: eventFilter,
	}, receiveMsgCh)
	if err != nil {
		panic(err)
	}

	for {
		select {
		// receive Mgo event
		case msg := <-receiveMsgCh:
			utils.PrettyPrint(msg)
		case <-ctx.Done():
			return
		}
	}

}

func TestSubscribeTransaction(t *testing.T) {
	receiveMsgCh := make(chan response.MgoEffects, 10)
	err := devWsCli.SubscribeTransaction(ctx, request.MgoSubscribeTransactionsRequest{
		TransactionFilter: request.TransactionFilterByFromAddress{
			FromAddress: "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
		},
	}, receiveMsgCh)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		// receive
		case msg := <-receiveMsgCh:
			utils.PrettyPrint(msg)
		case <-ctx.Done():
			return
		}
	}
}
