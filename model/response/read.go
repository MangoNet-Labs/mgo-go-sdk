package response

type EpochRollingGasCostSummary struct {
	ComputationCost         string `json:"computationCost"`
	StorageCost             string `json:"storageCost"`
	StorageRebate           string `json:"storageRebate"`
	NonRefundableStorageFee string `json:"nonRefundableStorageFee"`
}

type CheckpointResponse struct {
	Epoch                      string                     `json:"epoch"`
	SequenceNumber             string                     `json:"sequenceNumber"`
	Digest                     string                     `json:"digest"`
	NetworkTotalTransactions   string                     `json:"networkTotalTransactions"`
	PreviousDigest             string                     `json:"previousDigest"`
	EpochRollingGasCostSummary EpochRollingGasCostSummary `json:"epochRollingGasCostSummary"`
	TimestampMs                string                     `json:"timestampMs"`
	Transactions               []string                   `json:"transactions"`
	CheckpointCommitments      []interface{}              `json:"checkpointCommitments"`
	ValidatorSignature         string                     `json:"validatorSignature"`
}

type PaginatedCheckpointsResponse struct {
	Data        []CheckpointResponse `json:"data"`
	NextCursor  string               `json:"nextCursor"`
	HasNextPage bool                 `json:"hasNextPage"`
}

type MgoLoadedChildObject struct {
	ObjectID       string `json:"objectId"`
	SequenceNumber string `json:"sequenceNumber"`
}

type ChildObjectsResponse struct {
	LoadedChildObjects []*MgoLoadedChildObject `json:"loadedChildObjects"`
}

type ProtocolConfigResponse struct {
	MinSupportedProtocolVersion string                       `json:"minSupportedProtocolVersion"`
	MaxSupportedProtocolVersion string                       `json:"maxSupportedProtocolVersion"`
	ProtocolVersion             string                       `json:"protocolVersion"`
	FeatureFlags                map[string]bool              `json:"featureFlags"`
	Attributes                  map[string]map[string]string `json:"attributes"`
}

type MgoMultiGetTransactionBlocksResponse []*MgoTransactionBlockResponse

type PastObjectResponse struct {
	Status  string      `json:"status"`
	Details interface{} `json:"details"`
}
