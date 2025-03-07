package client

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mangonet-labs/mgo-go-sdk/client/httpconn"
	"github.com/mangonet-labs/mgo-go-sdk/model"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/model/response"

	"github.com/tidwall/gjson"
)

// BatchTransaction implements the method `unsafe_batchTransaction`, creates an unsigned batched transaction.
func (c *Client) BatchTransaction(ctx context.Context, req request.BatchTransactionRequest) (response.BatchTransactionResponse, error) {
	var rsp response.BatchTransactionResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_batchTransaction",
		Params: []interface{}{
			req.Signer,
			req.RPCTransactionRequestParams,
			req.Gas,
			req.GasBudget,
			req.MgoTransactionBlockBuilderMode,
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

// MergeCoins implements the method `unsafe_mergeCoins`, creates an unsigned transaction to merge multiple coins into one coin.
func (c *Client) MergeCoins(ctx context.Context, req request.MergeCoinsRequest) (model.TxnMetaData, error) {
	var rsp model.TxnMetaData
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_mergeCoins",
		Params: []interface{}{
			req.Signer,
			req.PrimaryCoin,
			req.CoinToMerge,
			req.Gas,
			req.GasBudget,
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

func (c *Client) MoveCall(ctx context.Context, req request.MoveCallRequest) (model.TxnMetaData, error) {
	var rsp model.TxnMetaData
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_moveCall",
		Params: []interface{}{
			req.Signer,
			req.PackageObjectId,
			req.Module,
			req.Function,
			req.TypeArguments,
			req.Arguments,
			req.Gas,
			req.GasBudget,
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

// Pay implements the method `unsafe_pay`, send `Coin<T>` to a list of addresses, where `T` can be any coin type, following a list of amounts.
// The object specified in the `gas` field will be used to pay the gas fee for the transaction.
// The gas object can not appear in `input_coins`. If the gas object is not specified, the RPC server will auto-select one.
func (c *Client) Pay(ctx context.Context, req request.PayRequest) (model.TxnMetaData, error) {
	var rsp model.TxnMetaData
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_pay",
		Params: []interface{}{
			req.Signer,
			req.MgoObjectId,
			req.Recipient,
			req.Amount,
			req.Gas,
			req.GasBudget,
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

// PayAllMgo implements the method `unsafe_payAllMgo`, send all MGO coins to one recipient.
// This is for MGO coin only and does not require a separate gas coin object.
// Specifically, what pay_all_mgo does are:
// 1. accumulate all MGO from input coins and deposit all MGO to the first input coin.
// 2. transfer the updated first coin to the recipient and also use this first coin as gas coin object.
// 3. the balance of the first input coin after tx is sum(input_coins) - actual_gas_cost.
// 4. all other input coins other than the first are deleted.
func (c *Client) PayAllMgo(ctx context.Context, req request.PayAllMgoRequest) (model.TxnMetaData, error) {
	var rsp model.TxnMetaData
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_payAllMgo",
		Params: []interface{}{
			req.Signer,
			req.MgoObjectId,
			req.Recipient,
			req.GasBudget,
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

// PayMgo implements the method `unsafe_payMgo`, send MGO coins to a list of addresses, following a list of amounts.
// This is for MGO coin only and does not require a separate gas coin object.
// Specifically, what pay_mgo does are:
// 1. debit each input_coin to create new coin following the order of amounts and assign it to the corresponding recipient.
// 2. accumulate all residual MGO from input coins left and deposit all MGO to the first input coin, then use the first input coin as the gas coin object.
// 3. the balance of the first input coin after tx is sum(input_coins) - sum(amounts) - actual_gas_cost
// 4. all other input coints other than the first one are deleted.
func (c *Client) PayMgo(ctx context.Context, req request.PayMgoRequest) (model.TxnMetaData, error) {
	var rsp model.TxnMetaData
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_payMgo",
		Params: []interface{}{
			req.Signer,
			req.MgoObjectId,
			req.Recipient,
			req.Amount,
			req.GasBudget,
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

// Publish implements the method `unsafe_publish`, creates an unsigned transaction to publish a Move package.
func (c *Client) Publish(ctx context.Context, req request.PublishRequest) (model.TxnMetaData, error) {
	var rsp model.TxnMetaData
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_publish",
		Params: []interface{}{
			req.Sender,
			req.CompiledModules,
			req.Dependencies,
			req.Gas,
			req.GasBudget,
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

// RequestAddStake implements the method `unsafe_requestAddStake`, add stake to a validator's staking pool using multiple coins and amount.
func (c *Client) RequestAddStake(ctx context.Context, req request.AddStakeRequest) (model.TxnMetaData, error) {
	var rsp model.TxnMetaData
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_requestAddStake",
		Params: []interface{}{
			req.Signer,
			req.Coins,
			req.Amount,
			req.Validator,
			req.Gas,
			req.GasBudget,
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

// RequestWithdrawStake implements the method `unsafe_requestWithdrawStake`, withdraw stake from a validator's staking pool.
func (c *Client) RequestWithdrawStake(ctx context.Context, req request.WithdrawStakeRequest) (model.TxnMetaData, error) {
	var rsp model.TxnMetaData
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_requestWithdrawStake",
		Params: []interface{}{
			req.Signer,
			req.StakedObjectId,
			req.Gas,
			req.GasBudget,
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

// SplitCoin implements the method `unsafe_splitCoin`, creates an unsigned transaction to split a coin object into multiple coins.
func (c *Client) SplitCoin(ctx context.Context, req request.SplitCoinRequest) (model.TxnMetaData, error) {
	var rsp model.TxnMetaData
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_splitCoin",
		Params: []interface{}{
			req.Signer,
			req.CoinObjectId,
			req.SplitAmounts,
			req.Gas,
			req.GasBudget,
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

// SplitCoinEqual implements the method `unsafe_splitCoinEqual`, creates an unsigned transaction to split a coin object into multiple equal-size coins.
func (c *Client) SplitCoinEqual(ctx context.Context, req request.SplitCoinEqualRequest) (model.TxnMetaData, error) {
	var rsp model.TxnMetaData
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_splitCoinEqual",
		Params: []interface{}{
			req.Signer,
			req.CoinObjectId,
			req.SplitCount,
			req.Gas,
			req.GasBudget,
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

// TransferObject implements the method `unsafe_transferObject`, creates an unsigned transaction to transfer a specified object from a signer to a recipient.
func (c *Client) TransferObject(ctx context.Context, req request.TransferObjectRequest) (model.TxnMetaData, error) {
	var rsp model.TxnMetaData
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_transferObject",
		Params: []interface{}{
			req.Signer,
			req.ObjectId,
			req.Gas,
			req.GasBudget,
			req.Recipient,
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

// TransferMgo implements the method `unsafe_transferMgo`, creates an unsigned transaction to transfer a specified amount of MGO from a signer to a recipient.
func (c *Client) TransferMgo(ctx context.Context, req request.TransferMgoRequest) (model.TxnMetaData, error) {
	var rsp model.TxnMetaData
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "unsafe_transferMgo",
		Params: []interface{}{
			req.Signer,
			req.MgoObjectId,
			req.GasBudget,
			req.Recipient,
			req.Amount,
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
