package response

type BatchTransactionResponse struct {
	Gas          []MgoObjectRef `json:"gas"`
	InputObjects []interface{}  `json:"inputObjects"`
	TxBytes      string         `json:"txBytes"`
}
