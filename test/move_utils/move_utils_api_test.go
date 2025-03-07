package move_utils

import (
	"context"
	"testing"

	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/utils"
)

var ctx = context.Background()
var devCli = client.NewMgoClient(config.RpcMgoTestnetEndpoint)

func TestGetMoveFunctionArgTypes(t *testing.T) {
	types, err := devCli.MgoGetMoveFunctionArgTypes(ctx, request.GetMoveFunctionArgTypesRequest{
		Package:  "0x0000000000000000000000000000000000000000000000000000000000000003",
		Module:   "mgo_system",
		Function: "report_validator",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(types)
}

func TestGetNormalizedMoveFunction(t *testing.T) {
	types, err := devCli.MgoGetNormalizedMoveFunction(ctx, request.GetNormalizedMoveFunctionRequest{
		Package:      "0x0000000000000000000000000000000000000000000000000000000000000003",
		ModuleName:   "mgo_system",
		FunctionName: "report_validator",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(types)
}

func TestGetNormalizedMoveModule(t *testing.T) {
	types, err := devCli.MgoGetNormalizedMoveModule(ctx, request.GetNormalizedMoveModuleRequest{
		Package:    "0x0000000000000000000000000000000000000000000000000000000000000003",
		ModuleName: "mgo_system",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(types)
}

func TestGetNormalizedMoveModulesByPackage(t *testing.T) {
	types, err := devCli.MgoGetNormalizedMoveModulesByPackage(ctx, request.GetNormalizedMoveModulesByPackageRequest{
		Package: "0x0000000000000000000000000000000000000000000000000000000000000003",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(types)
}

func TestGetNormalizedMoveStruct(t *testing.T) {
	types, err := devCli.MgoGetNormalizedMoveStruct(ctx, request.GetNormalizedMoveStructRequest{
		Package:    "0x0000000000000000000000000000000000000000000000000000000000000003",
		ModuleName: "mgo_system",
		StructName: "MgoSystemState",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(types)
}
