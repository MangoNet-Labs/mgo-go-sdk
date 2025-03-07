package client

import (
	"context"
	"errors"

	"github.com/MangoNet-Labs/mgo-go-sdk/client/httpconn"

	"github.com/tidwall/gjson"
)

// MgoCall is a generic method that calls a JSON-RPC method on the server. Note that this method
// does not perform any error checking on the response, and just returns the raw response as a string.
// If the method call resulted in an error, the error is returned as a string.
func (c *Client) MgoCall(ctx context.Context, method string, params ...interface{}) (interface{}, error) {
	resp, err := c.conn.Request(ctx, httpconn.Operation{
		Method: method,
		Params: params,
	})
	if err != nil {
		return nil, err
	}
	if gjson.ParseBytes(resp).Get("error").Exists() {
		return nil, errors.New(gjson.ParseBytes(resp).Get("error").String())
	}
	return gjson.ParseBytes(resp).String(), nil
}
