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

// MgoXGetAllBalance implements the method `mgox_getAllBalances`, gets all Coin balances owned by an address.
func (c *Client) MgoXGetAllBalance(ctx context.Context, req request.MgoXGetAllBalanceRequest) (response.CoinAllBalanceResponse, error) {
	var rsp response.CoinAllBalanceResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_getAllBalances",
		Params: []interface{}{
			req.Owner,
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

// MgoXGetAllCoins implements the method `mgox_getAllCoins`, gets all Coin objects owned by an address.
func (c *Client) MgoXGetAllCoins(ctx context.Context, req request.MgoXGetAllCoinsRequest) (response.PaginatedCoinsResponse, error) {
	var rsp response.PaginatedCoinsResponse
	if err := validate.ValidateStruct(req); err != nil {
		return rsp, err
	}
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_getAllCoins",
		Params: []interface{}{
			req.Owner,
			req.Cursor,
			req.Limit,
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

// MgoXGetBalance implements the method `mgox_getBalance`, gets the total Coin balance for each coin type owned by an address.
func (c *Client) MgoXGetBalance(ctx context.Context, req request.MgoXGetBalanceRequest) (response.CoinBalanceResponse, error) {
	var rsp response.CoinBalanceResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_getBalance",
		Params: []interface{}{
			req.Owner,
			req.CoinType,
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

// MgoXGetCoinMetadata implements the method `mgox_getCoinMetadata`, gets metadata(e.g., symbol, decimals) for a coin.
func (c *Client) MgoXGetCoinMetadata(ctx context.Context, req request.MgoXGetCoinMetadataRequest) (response.CoinMetadataResponse, error) {
	var rsp response.CoinMetadataResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_getCoinMetadata",
		Params: []interface{}{
			req.CoinType,
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

// MgoXGetCoins implements the method `mgox_getCoins`, gets a list of Coin objects by type owned by an address.
func (c *Client) MgoXGetCoins(ctx context.Context, req request.MgoXGetCoinsRequest) (response.PaginatedCoinsResponse, error) {
	var rsp response.PaginatedCoinsResponse
	if err := validate.ValidateStruct(req); err != nil {
		return rsp, err
	}
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_getCoins",
		Params: []interface{}{
			req.Owner,
			req.CoinType,
			req.Cursor,
			req.Limit,
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

// MgoXGetTotalSupply implements the method `mgox_getTotalSupply`, gets total supply for a coin
func (c *Client) MgoXGetTotalSupply(ctx context.Context, req request.MgoXGetTotalSupplyRequest) (response.TotalSupplyResponse, error) {
	var rsp response.TotalSupplyResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_getTotalSupply",
		Params: []interface{}{
			req.CoinType,
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
