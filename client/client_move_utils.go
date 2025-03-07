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

// MgoGetMoveFunctionArgTypes implements method `mgo_getMoveFunctionArgTypes`, return the argument types of a Move function based on normalized type.
func (c *Client) MgoGetMoveFunctionArgTypes(ctx context.Context, req request.GetMoveFunctionArgTypesRequest) (response.GetMoveFunctionArgTypesResponse, error) {
	var rsp response.GetMoveFunctionArgTypesResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getMoveFunctionArgTypes",
		Params: []interface{}{
			req.Package,
			req.Module,
			req.Function,
		},
	})
	if err != nil {
		return rsp, nil
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

// MgoGetNormalizedMoveFunction implements method `mgo_getNormalizedMoveFunction`, return a structured representation of a Move function.
func (c *Client) MgoGetNormalizedMoveFunction(ctx context.Context, req request.GetNormalizedMoveFunctionRequest) (response.GetNormalizedMoveFunctionResponse, error) {
	var rsp response.GetNormalizedMoveFunctionResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getNormalizedMoveFunction",
		Params: []interface{}{
			req.Package,
			req.ModuleName,
			req.FunctionName,
		},
	})
	if err != nil {
		return rsp, nil
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

// MgoGetNormalizedMoveModule implements method `mgo_getNormalizedMoveModule`, return a structured representation of a Move module.
func (c *Client) MgoGetNormalizedMoveModule(ctx context.Context, req request.GetNormalizedMoveModuleRequest) (response.GetNormalizedMoveModuleResponse, error) {
	var rsp response.GetNormalizedMoveModuleResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getNormalizedMoveModule",
		Params: []interface{}{
			req.Package,
			req.ModuleName,
		},
	})
	if err != nil {
		return rsp, nil
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

// MgoGetNormalizedMoveModulesByPackage implements method `mgo_getNormalizedMoveModulesByPackage`, return the structured representations of all modules in the given package.
func (c *Client) MgoGetNormalizedMoveModulesByPackage(ctx context.Context, req request.GetNormalizedMoveModulesByPackageRequest) (response.GetNormalizedMoveModulesByPackageResponse, error) {
	var rsp response.GetNormalizedMoveModulesByPackageResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getNormalizedMoveModulesByPackage",
		Params: []interface{}{
			req.Package,
		},
	})
	if err != nil {
		return rsp, nil
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

// MgoGetNormalizedMoveStruct implements method `mgo_getNormalizedMoveStruct`, return a structured representation of a Move struct.
func (c *Client) MgoGetNormalizedMoveStruct(ctx context.Context, req request.GetNormalizedMoveStructRequest) (response.GetNormalizedMoveStructResponse, error) {
	var rsp response.GetNormalizedMoveStructResponse
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_getNormalizedMoveStruct",
		Params: []interface{}{
			req.Package,
			req.ModuleName,
			req.StructName,
		},
	})
	if err != nil {
		return rsp, nil
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
