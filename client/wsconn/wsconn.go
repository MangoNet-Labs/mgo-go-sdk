package wsconn

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mangonet-labs/mgo-go-sdk/model/request"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

type WsConn struct {
	Conn   *websocket.Conn
	wsUrl  string
	ticker *time.Ticker // For heartbeat, default 30s
}

type CallOp struct {
	Method string
	Params []interface{}
}

// NewWsConn returns a new WsConn given a websocket url. It dials the url and
// returns a new WsConn. The WsConn has a ticker that sends a ping every 30
// seconds to keep the connection alive.
func NewWsConn(wsUrl string) *WsConn {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsUrl, nil)

	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err, wsUrl)
	}

	return &WsConn{
		Conn:   conn,
		wsUrl:  wsUrl,
		ticker: time.NewTicker(30 * time.Second),
	}
}

// NewWsConnWithDuration returns a new WsConn given a websocket url and a custom duration.
// It dials the url and returns a new WsConn. The WsConn has a ticker that sends a ping
// at the specified duration to keep the connection alive.
func NewWsConnWithDuration(wsUrl string, d time.Duration) *WsConn {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsUrl, nil)

	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err, wsUrl)
	}

	return &WsConn{
		Conn:   conn,
		wsUrl:  wsUrl,
		ticker: time.NewTicker(d),
	}
}

// Call sends a JSON-RPC request over the websocket connection and listens for responses.
// It marshals the provided CallOp into a JSON-RPC request, writes it to the websocket,
// and then reads messages from the websocket connection. If a message contains an error,
// it returns the error. It starts a goroutine to periodically send ping messages to
// keep the connection alive. Another goroutine listens for incoming messages and sends
// them to the provided channel. The function returns an error if there are issues with
// sending or receiving messages.
//
// Parameters:
// - ctx: The context to control cancellation and timeout.
// - op: The CallOp containing method and parameters for the RPC call.
// - receiveMsgCh: A channel to receive messages from the websocket.
//
// Returns:
// - error: An error if the call fails at any point.
func (w *WsConn) Call(ctx context.Context, op CallOp, receiveMsgCh chan []byte) error {
	jsonRPCCall := request.JsonRPCRequest{
		JsonRPC: "2.0",
		ID:      time.Now().UnixMilli(),
		Method:  op.Method,
		Params:  op.Params,
	}

	callBytes, err := json.Marshal(jsonRPCCall)
	if err != nil {
		return err
	}

	err = w.Conn.WriteMessage(websocket.TextMessage, callBytes)
	if nil != err {
		return err
	}

	_, messageData, err := w.Conn.ReadMessage()
	if nil != err {
		return err
	}

	var rsp SubscriptionResp
	if gjson.ParseBytes(messageData).Get("error").Exists() {
		return fmt.Errorf(gjson.ParseBytes(messageData).Get("error").String())
	}

	err = json.Unmarshal([]byte(gjson.ParseBytes(messageData).String()), &rsp)
	if err != nil {
		return err
	}

	fmt.Printf("establish successfully, subscriptionID: %d, Waiting to accept data...\n", rsp.Result)
	go func() {
		for {
			select {
			case <-ctx.Done():
				w.ticker.Stop()
				return
			case <-w.ticker.C:
				if err := w.Conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
					log.Println("ping failed:", err)
					return
				}
			}
		}
	}()

	go func(conn *websocket.Conn) {
		for {
			messageType, messageData, err := conn.ReadMessage()
			if nil != err {
				log.Println(err)
				break
			}
			switch messageType {
			case websocket.TextMessage:
				receiveMsgCh <- messageData

			default:
				continue
			}
		}
	}(w.Conn)

	return nil
}
