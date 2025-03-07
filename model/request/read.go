package request

type MgoGetCheckpointRequest struct {
	CheckpointID string `json:"id"`
}

type MgoGetCheckpointsRequest struct {
	// optional paging cursor
	Cursor interface{} `json:"cursor"`
	// maximum number of items per page
	Limit uint64 `json:"limit" validate:"lte=50"`
	// query result ordering, default to false (ascending order), oldest record first
	DescendingOrder bool `json:"descendingOrder"`
}

type MgoGetLoadedChildObjectsRequest struct {
	Digest string `json:"digest"`
}

type MgoGetObjectRequest struct {
	// the ID of the queried object
	ObjectId string `json:"ObjectId"`
	// config which fields to include in the response, by default only digest is included
	Options MgoObjectDataOptions `json:"options"`
}

type MgoGetProtocolConfigRequest struct {
	Version string `json:"version"`
}

type MgoGetTransactionBlockRequest struct {
	Digest  string                     `json:"digest"`
	Options MgoTransactionBlockOptions `json:"options"`
}

type MgoMultiGetObjectsRequest struct {
	ObjectIds []string             `json:"objectIds"`
	Options   MgoObjectDataOptions `json:"options"`
}

type MgoMultiGetTransactionBlocksRequest struct {
	Digests []string                   `json:"digests"`
	Options MgoTransactionBlockOptions `json:"options"`
}

type MgoTryGetPastObjectRequest struct {
	// the ID of the queried object
	ObjectId string `json:"objectId"`
	// the version of the queried object
	Version uint64               `json:"version"`
	Options MgoObjectDataOptions `json:"options"`
}

type PastObject struct {
	// the ID of the queried object
	ObjectId string `json:"objectId"`
	// the version of the queried object
	Version string `json:"version"`
}

type MgoTryMultiGetPastObjectsRequest struct {
	// a vector of object and versions to be queried
	MultiGetPastObjects []*PastObject
	// options for specifying the content to be returned
	Options MgoObjectDataOptions
}
