package response

type MgoXResolveNameServiceNamesResponse struct {
	Data        []string `json:"data"`
	NextCursor  string   `json:"nextCursor"`
	HasNextPage bool     `json:"hasNextPage"`
}

type MgoParsedData struct {
	DataType string `json:"dataType"`
	MgoMoveObject
	MgoMovePackage
}

type MgoObjectData struct {
	ObjectId            string                `json:"objectId"`
	Version             string                `json:"version"`
	Digest              string                `json:"digest"`
	Type                string                `json:"type"`
	Owner               interface{}           `json:"owner"`
	PreviousTransaction string                `json:"previousTransaction,omitempty"`
	Display             DisplayFieldsResponse `json:"display"`
	Content             *MgoParsedData        `json:"content,omitempty"`
	Bcs                 *MgoRawData           `json:"bcs,omitempty"`
}

type MgoObjectResponseError struct {
	Code     string `json:"code"`
	Error    string `json:"error"`
	ObjectId string `json:"object_id"`
	Version  int    `json:"version"`
	Digest   string `json:"digest"`
}

type MgoObjectResponse struct {
	Data  *MgoObjectData          `json:"data,omitempty"`
	Error *MgoObjectResponseError `json:"error,omitempty"`
}

type PaginatedObjectsResponse struct {
	Data        []MgoObjectResponse `json:"data"`
	NextCursor  string              `json:"nextCursor"`
	HasNextPage bool                `json:"hasNextPage"`
}

type MgoXQueryTransactionBlocksResponse struct {
	Data        []MgoTransactionBlockResponse `json:"data"`
	NextCursor  string                        `json:"nextCursor"`
	HasNextPage bool                          `json:"hasNextPage"`
}
