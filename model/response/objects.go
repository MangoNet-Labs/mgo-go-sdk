package response

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type Owner struct {
	AddressOwner string `json:"addressOwner,omitempty"`
	ObjectOwner  string `json:"objectOwner,omitempty"`
}

type MgoObjectInfo struct {
	MgoObjectRef
	Type string `json:"type_"`
	Owner
	PreviousTransaction string `json:"previousTransaction"`
}

type MgoMoveObject struct {
	Type              string                 `json:"type"`
	Fields            map[string]interface{} `json:"fields"`
	HasPublicTransfer bool                   `json:"hasPublicTransfer"`
}

type MgoMovePackage struct {
	Disassembled interface{} `json:"disassembled"`
}

type MgoMoveModuleId struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type MgoMoveNormalizedModule struct {
	FileFormatVersion uint64
	Address           string
	Name              string
	Friends           []MgoMoveModuleId
}

type MgoRawMovePackage struct {
	Id        string            `json:"id,omitempty"`
	ModuleMap map[string]string `json:"moduleMap,omitempty"`
}

type MgoRawMoveObject struct {
	Type              string `json:"type"`
	HasPublicTransfer bool   `json:"hasPublicTransfer"`
	Version           int    `json:"version"`
	BcsBytes          string `json:"bcsBytes"`
}

type MgoRawData struct {
	DataType string `json:"dataType"`
	MgoRawMoveObject
	MgoRawMovePackage
}

type DynamicFieldName struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

func (v DynamicFieldName) Field(field string) gjson.Result {
	return gjson.GetBytes(v.Value, field)
}

type DisplayFieldsResponse struct {
	Data  any                     `json:"data"`
	Error *MgoObjectResponseError `json:"error"`
}

type DynamicFieldInfo struct {
	Name       DynamicFieldName `json:"name"`
	BcsName    string           `json:"bcsName"`
	Type       string           `json:"type"`
	ObjectType string           `json:"objectType"`
	ObjectId   string           `json:"objectId"`
	Version    int              `json:"version"`
	Digest     string           `json:"digest"`
}

type PaginatedDynamicFieldInfoResponse struct {
	Data        []DynamicFieldInfo `json:"data"`
	NextCursor  string             `json:"nextCursor"`
	HasNextPage bool               `json:"hasNextPage"`
}
