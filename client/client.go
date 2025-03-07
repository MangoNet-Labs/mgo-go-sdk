package client

import (
	"github.com/mangonet-labs/mgo-go-sdk/client/httpconn"
)

type Client struct {
	conn *httpconn.HttpConn
}

// NewMgoClient instantiates a new Client with the given net identity.
//
// It is used to create a new client for calling the methods of each module.
//
// Note that the net identity is used to determine the RPC URL of the full node.
func NewMgoClient(rpcUrl string) *Client {
	conn := httpconn.NewHttpConn(rpcUrl)
	return newClient(conn)
}

// newClient creates a new Client with the given http connection and net identity.
//
// The connection is used to send the RPC requests to the full node, and the net
// identity is used to determine the RPC URL of the full node.
func newClient(conn *httpconn.HttpConn) *Client {
	return &Client{
		conn: conn,
	}
}
