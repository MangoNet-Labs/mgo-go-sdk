package client

import (
	"time"

	"github.com/MangoNet-Labs/mgo-go-sdk/client/wsconn"
)

// IMgoWebsocketAPI defines the subscription API related interface, and then implement it by the WebsocketClient.
type IMgoWebsocketAPI interface {
	ISubscribeAPI
}

// WebsocketClient implements MgoWebsocketAPI related interfaces.
type WebsocketClient struct {
	ISubscribeAPI
}

// NewMgoWebsocketClient instantiates the WebsocketClient to call the methods of each module.
func NewMgoWebsocketClient(rpcUrl string) IMgoWebsocketAPI {
	conn := wsconn.NewWsConn(rpcUrl)
	return &WebsocketClient{
		ISubscribeAPI: &mgoSubscribeImpl{
			conn: conn,
		},
	}
}

// NewMgoWebsocketClientWithDuration instantiates the WebsocketClient to call the methods of each module, parameter d is for sending heartbeat kit
func NewMgoWebsocketClientWithDuration(rpcUrl string, d time.Duration) IMgoWebsocketAPI {
	conn := wsconn.NewWsConnWithDuration(rpcUrl, d)
	return &WebsocketClient{
		ISubscribeAPI: &mgoSubscribeImpl{
			conn: conn,
		},
	}
}
