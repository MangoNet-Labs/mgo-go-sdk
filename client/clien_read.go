package client

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/MangoNet-Labs/mgo-go-sdk/client/httpconn"
	"github.com/MangoNet-Labs/mgo-go-sdk/model/request"
	"github.com/MangoNet-Labs/mgo-go-sdk/model/response"
	"github.com/MangoNet-Labs/mgo-go-sdk/utils"

	"github.com/tidwall/gjson"
)

// MgoGetChainIdentifier implements the method `mgo_getChainIdentifier`, return the chain's identifier.
func (c *Client) MgoGetChainIdentifier(ctx context.Context) (string, error) {
	var rsp string
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getChainIdentifier",
		Params: []interface{}{},
	})
	if err != nil {
		return rsp, err
	}

	if gjson.ParseBytes(respBytes).Get("error").Exists() {
		return rsp, errors.New(gjson.ParseBytes(respBytes).Get("error").String())
	}
	rsp = gjson.ParseBytes(respBytes).Get("result").String()
	return rsp, nil
}

// MgoGetCheckpoint implements the method `mgo_getCheckpoint`, gets a checkpoint.
func (c *Client) MgoGetCheckpoint(ctx context.Context, req request.MgoGetCheckpointRequest) (response.CheckpointResponse, error) {
	var rsp response.CheckpointResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getCheckpoint",
		Params: []interface{}{
			req.CheckpointID,
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

// MgoGetCheckpoints implements the method `mgo_getCheckpoints`, gets paginated list of checkpoints.
func (c *Client) MgoGetCheckpoints(ctx context.Context, req request.MgoGetCheckpointsRequest) (response.PaginatedCheckpointsResponse, error) {
	var rsp response.PaginatedCheckpointsResponse
	if err := validate.ValidateStruct(req); err != nil {
		return rsp, err
	}
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getCheckpoints",
		Params: []interface{}{
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

// MgoGetEvents implements the method `mgo_getEvents`, gets transaction events.
func (c *Client) MgoGetEvents(ctx context.Context, req request.MgoGetEventsRequest) (response.GetEventsResponse, error) {
	var rsp response.GetEventsResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getEvents",
		Params: []interface{}{
			req.Digest,
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

// MgoGetLatestCheckpointSequenceNumber implements the method `mgo_getLatestCheckpointSequenceNumber`, gets the sequence number of the latest checkpoint that has been executed.
func (c *Client) MgoGetLatestCheckpointSequenceNumber(ctx context.Context) (uint64, error) {
	var rsp uint64
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getLatestCheckpointSequenceNumber",
		Params: []interface{}{},
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

// MgoGetLoadedChildObjects implements the method `mgo_getLoadedChildObjects`, return the object information for a specified digest.
func (c *Client) MgoGetLoadedChildObjects(ctx context.Context, req request.MgoGetLoadedChildObjectsRequest) (response.ChildObjectsResponse, error) {
	var rsp response.ChildObjectsResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getLoadedChildObjects",
		Params: []interface{}{
			req.Digest,
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

// MgoGetObject implements the method `mgo_getObject`, gets the object information for a specified object.
func (c *Client) MgoGetObject(ctx context.Context, req request.MgoGetObjectRequest) (response.MgoObjectResponse, error) {
	var rsp response.MgoObjectResponse

	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getObject",
		Params: []interface{}{
			req.ObjectId,
			req.Options,
		},
	})
	if err != nil {
		return rsp, err
	}

	err = json.Unmarshal([]byte(gjson.ParseBytes(respBytes).Get("result").String()), &rsp)
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

// MgoGetProtocolConfig implements the method `mgo_getProtocolConfig`, return the protocol config table for the given version number.
// If the version number is not specified, If none is specified, the node uses the version of the latest epoch it has processed.
func (c *Client) MgoGetProtocolConfig(ctx context.Context, req request.MgoGetProtocolConfigRequest) (response.ProtocolConfigResponse, error) {
	var rsp response.ProtocolConfigResponse
	params := make([]interface{}, 0)
	if utils.IsFieldNonEmpty(req, "Version") {
		params = append(params, req.Version)
	}
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getProtocolConfig",
		Params: params,
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

// MgoGetTotalTransactionBlocks implements the method `mgo_getTotalTransactionBlocks`, gets the total number of transactions known to the node.
func (c *Client) MgoGetTotalTransactionBlocks(ctx context.Context) (uint64, error) {
	var rsp uint64
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getTotalTransactionBlocks",
		Params: []interface{}{},
	})
	if err != nil {
		return rsp, err
	}
	rsp = gjson.ParseBytes(respBytes).Get("result").Uint()
	return rsp, nil
}

// MgoGetTransactionBlock implements the method `mgo_getTransactionBlock`, gets the transaction response object for a specified transaction digest.
func (c *Client) MgoGetTransactionBlock(ctx context.Context, req request.MgoGetTransactionBlockRequest) (response.MgoTransactionBlockResponse, error) {
	var rsp response.MgoTransactionBlockResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getTransactionBlock",
		Params: []interface{}{
			req.Digest,
			req.Options,
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

// MgoMultiGetObjects implements the method `mgo_multiGetObjects`, gets the object data for a list of objects.
func (c *Client) MgoMultiGetObjects(ctx context.Context, req request.MgoMultiGetObjectsRequest) ([]*response.MgoObjectResponse, error) {
	var rsp []*response.MgoObjectResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_multiGetObjects",
		Params: []interface{}{
			req.ObjectIds,
			req.Options,
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

// MgoMultiGetTransactionBlocks implements the method `mgo_multiGetTransactionBlocks`, gets an ordered list of transaction responses.
func (c *Client) MgoMultiGetTransactionBlocks(ctx context.Context, req request.MgoMultiGetTransactionBlocksRequest) (response.MgoMultiGetTransactionBlocksResponse, error) {
	var rsp response.MgoMultiGetTransactionBlocksResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_multiGetTransactionBlocks",
		Params: []interface{}{
			req.Digests,
			req.Options,
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

// MgoTryGetPastObject implements the method `mgo_tryGetPastObject`, gets the object information for a specified version.
// There is no guarantee that objects with past versions can be retrieved by this API. The result may vary across nodes depending on their pruning policies.
func (c *Client) MgoTryGetPastObject(ctx context.Context, req request.MgoTryGetPastObjectRequest) (response.PastObjectResponse, error) {
	var rsp response.PastObjectResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_tryGetPastObject",
		Params: []interface{}{
			req.ObjectId,
			req.Version,
			req.Options,
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

// MgoTryMultiGetPastObjects implements the method `mgo_tryMultiGetPastObjects`,
// note there is no software-level guarantee/SLA that objects with past versions can be retrieved by this API,
// even if the object and version exists/existed. The result may vary across nodes depending on their pruning policies.
// Return the object information for a specified version.
func (c *Client) MgoTryMultiGetPastObjects(ctx context.Context, req request.MgoTryMultiGetPastObjectsRequest) ([]*response.PastObjectResponse, error) {
	var rsp []*response.PastObjectResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_tryMultiGetPastObjects",
		Params: []interface{}{
			req.MultiGetPastObjects,
			req.Options,
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
