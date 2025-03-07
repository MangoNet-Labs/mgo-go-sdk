package request

type MgoDevInspectTransactionBlockRequest struct {
	// the transaction signer's Mgo address
	Sender string `json:"sender"`
	// BCS encoded TransactionKind(as opposed to TransactionData, which include gasBudget and gasPrice)
	TxBytes string `json:"txBytes"`
	// Gas is not charged, but gas usage is still calculated. Default to use reference gas price
	GasPrice string `json:"gasPrice"`
	// The epoch to perform the call. Will be set from the system state object if not provided
	Epoch string `json:"epoch"`
}

type MgoDryRunTransactionBlockRequest struct {
	TxBytes string `json:"txBytes"`
}

type MgoExecuteTransactionBlockRequest struct {
	// BCS serialized transaction data bytes without its type tag, as base-64 encoded string.
	TxBytes string `json:"txBytes"`
	// A list of signatures (`flag || signature || pubkey` bytes, as base-64 encoded string).
	// Signature is committed to the intent message of the transaction data, as base-64 encoded string.
	Signature []string `json:"signature"`
	// Options for specifying the content to be returned
	Options MgoTransactionBlockOptions `json:"options"`
	// The request type, derived from `MgoTransactionBlockResponseOptions` if None.
	// The optional enumeration values are: `WaitForEffectsCert`, or `WaitForLocalExecution`
	RequestType string `json:"requestType"`
}
