package client

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/MangoNet-Labs/mgo-go-sdk/client/httpconn"
	"github.com/MangoNet-Labs/mgo-go-sdk/client/wsconn"
	"github.com/MangoNet-Labs/mgo-go-sdk/model/request"
	"github.com/MangoNet-Labs/mgo-go-sdk/model/response"

	"github.com/tidwall/gjson"
)

type ISubscribeAPI interface {
	SubscribeEvent(ctx context.Context, req request.MgoSubscribeEventsRequest, msgCh chan response.MgoEventResponse) error
	SubscribeTransaction(ctx context.Context, req request.MgoSubscribeTransactionsRequest, msgCh chan response.MgoEffects) error
}
type mgoSubscribeImpl struct {
	conn *wsconn.WsConn
}

// MgoXGetDynamicFieldObject implements the method `mgox_getDynamicFieldObject`, gets the dynamic field object information for a specified object.
func (c *Client) MgoXGetDynamicFieldObject(ctx context.Context, req request.MgoXGetDynamicFieldObjectRequest) (response.MgoObjectResponse, error) {
	var rsp response.MgoObjectResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_getDynamicFieldObject",
		Params: []interface{}{
			req.ObjectId,
			req.DynamicFieldName,
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

// MgoXGetDynamicField implements the method `mgpx_getDynamicFields`, gets the list of dynamic field objects owned by an object.
func (c *Client) MgoXGetDynamicFields(ctx context.Context, req request.MgoXGetDynamicFieldsRequest) (response.PaginatedDynamicFieldInfoResponse, error) {
	var rsp response.PaginatedDynamicFieldInfoResponse
	if err := validate.ValidateStruct(req); err != nil {
		return rsp, err
	}
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_getDynamicFields",
		Params: []interface{}{
			req.ObjectId,
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

// MgoXGetOwnedObjects implements the method `mgox_getOwnedObjects`, gets the list of objects owned by an address.
func (c *Client) MgoXGetOwnedObjects(ctx context.Context, req request.MgoXGetOwnedObjectsRequest) (response.PaginatedObjectsResponse, error) {
	var rsp response.PaginatedObjectsResponse
	if err := validate.ValidateStruct(req); err != nil {
		return rsp, err
	}
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_getOwnedObjects",
		Params: []interface{}{
			req.Address,
			req.Query,
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

// MgoXQueryEvents implements the method `mgox_queryEvents`, gets list of events for a specified query criteria.
func (c *Client) MgoXQueryEvents(ctx context.Context, req request.MgoXQueryEventsRequest) (response.PaginatedEventsResponse, error) {
	var rsp response.PaginatedEventsResponse
	if err := validate.ValidateStruct(req); err != nil {
		return rsp, err
	}
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_queryEvents",
		Params: []interface{}{
			req.MgoEventFilter,
			req.Cursor,
			req.Limit,
			req.DescendingOrder,
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

// MgoXQueryTransactionBlocks implements the method `mgox_queryTransactionBlocks`, gets list of transactions for a specified query criteria.
func (c *Client) MgoXQueryTransactionBlocks(ctx context.Context, req request.MgoXQueryTransactionBlocksRequest) (response.MgoXQueryTransactionBlocksResponse, error) {
	var rsp response.MgoXQueryTransactionBlocksResponse
	if err := validate.ValidateStruct(req); err != nil {
		return rsp, err
	}
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_queryTransactionBlocks",
		Params: []interface{}{
			req.MgoTransactionBlockResponseQuery,
			req.Cursor,
			req.Limit,
			req.DescendingOrder,
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

// MgoXResolveNameServiceAddress implements the method `mgox_resolveNameServiceAddress`, get the resolved address given resolver and name.
func (c *Client) MgoXResolveNameServiceAddress(ctx context.Context, req request.MgoXResolveNameServiceAddressRequest) (string, error) {
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_resolveNameServiceAddress",
		Params: []interface{}{
			req.Name,
		},
	})
	if err != nil {
		return "", err
	}
	if gjson.ParseBytes(respBytes).Get("error").Exists() {
		return "", errors.New(gjson.ParseBytes(respBytes).Get("error").String())
	}

	return gjson.ParseBytes(respBytes).Get("result").String(), nil
}

// MgoXResolveNameServiceNames implements the method `mgox_resolveNameServiceNames`, return the resolved names given address, if multiple names are resolved, the first one is the primary name.
func (c *Client) MgoXResolveNameServiceNames(ctx context.Context, req request.MgoXResolveNameServiceNamesRequest) (response.MgoXResolveNameServiceNamesResponse, error) {
	var rsp response.MgoXResolveNameServiceNamesResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgox_resolveNameServiceNames",
		Params: []interface{}{
			req.Address,
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

// SubscribeEvent implements the method `mgox_subscribeEvent`, subscribe to a stream of Mgo event.
func (s *mgoSubscribeImpl) SubscribeEvent(ctx context.Context, req request.MgoSubscribeEventsRequest, msgCh chan response.MgoEventResponse) error {
	rsp := make(chan []byte, 10)
	err := s.conn.Call(ctx, wsconn.CallOp{
		Method: "mgox_subscribeEvent",
		Params: []interface{}{
			req.MgoEventFilter,
		},
	}, rsp)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case messageData := <-rsp:
				var result response.MgoEventResponse
				if gjson.ParseBytes(messageData).Get("error").Exists() {
					log.Fatal(gjson.ParseBytes(messageData).Get("error").String())
				}

				err := json.Unmarshal([]byte(gjson.ParseBytes(messageData).Get("params.result").String()), &result)
				if err != nil {
					log.Fatal(err)
				}

				msgCh <- result
			}
		}
	}()

	return nil
}

// SubscribeTransaction implements the method `mgox_subscribeTransaction`, subscribe to a stream of Mgo transaction effects.
func (s *mgoSubscribeImpl) SubscribeTransaction(ctx context.Context, req request.MgoSubscribeTransactionsRequest, msgCh chan response.MgoEffects) error {
	rsp := make(chan []byte, 10)
	err := s.conn.Call(ctx, wsconn.CallOp{
		Method: "mgox_subscribeTransaction",
		Params: []interface{}{
			req.TransactionFilter,
		},
	}, rsp)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case messageData := <-rsp:
				var result response.MgoEffects
				if gjson.ParseBytes(messageData).Get("error").Exists() {
					log.Fatal(gjson.ParseBytes(messageData).Get("error").String())
				}

				err := json.Unmarshal([]byte(gjson.ParseBytes(messageData).Get("params.result").String()), &result)
				if err != nil {
					log.Fatal(err)
				}

				msgCh <- result
			}
		}
	}()

	return nil
}
