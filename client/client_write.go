package client

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mangonet-labs/mgo-go-sdk/client/httpconn"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/model/response"

	"github.com/tidwall/gjson"
)

// MgoDevInspectTransactionBlock implements the method `mgo_devInspectTransactionBlock`, returns the detailed information of a transaction block, without executing it. The transaction is signed using the specified sender, and gas price and epoch are used to calculate gas cost.
func (c *Client) MgoDevInspectTransactionBlock(ctx context.Context, req request.MgoDevInspectTransactionBlockRequest) (response.MgoTransactionBlockResponse, error) {
	var rsp response.MgoTransactionBlockResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_devInspectTransactionBlock",
		Params: []interface{}{
			req.Sender,
			req.TxBytes,
			req.GasPrice,
			req.Epoch,
		},
	})
	if err != nil {
		return rsp, err
	}
	if gjson.ParseBytes(respBytes).Get("error").Exists() {
		return rsp, errors.New(gjson.ParseBytes(respBytes).Get("error").String())
	}
	err = json.Unmarshal([]byte(gjson.ParseBytes(respBytes).Get("result").Raw), &rsp)
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

// MgoDryRunTransactionBlock implements the method `mgo_dryRunTransactionBlock`, gets transaction execution effects including the gas cost summary, while the effects are not committed to the chain.
func (c *Client) MgoDryRunTransactionBlock(ctx context.Context, req request.MgoDryRunTransactionBlockRequest) (response.MgoTransactionBlockResponse, error) {
	var rsp response.MgoTransactionBlockResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_dryRunTransactionBlock",
		Params: []interface{}{
			req.TxBytes,
		},
	})
	if err != nil {
		return rsp, err
	}
	if gjson.ParseBytes(respBytes).Get("error").Exists() {
		return rsp, errors.New(gjson.ParseBytes(respBytes).Get("error").String())
	}
	err = json.Unmarshal([]byte(gjson.ParseBytes(respBytes).Get("result").Raw), &rsp)
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

// MgoExecuteTransactionBlock implements the method `mgo_executeTransactionBlock`, executes a transaction using the transaction data and signature(s).
func (c *Client) MgoExecuteTransactionBlock(ctx context.Context, req request.MgoExecuteTransactionBlockRequest) (response.MgoTransactionBlockResponse, error) {
	var rsp response.MgoTransactionBlockResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_executeTransactionBlock",
		Params: []interface{}{
			req.TxBytes,
			req.Signature,
			req.Options,
			req.RequestType,
		},
	})
	if err != nil {
		return rsp, err
	}
	if gjson.ParseBytes(respBytes).Get("error").Exists() {
		return rsp, errors.New(gjson.ParseBytes(respBytes).Get("error").String())
	}
	err = json.Unmarshal([]byte(gjson.ParseBytes(respBytes).Get("result").String()), &rsp)
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}
