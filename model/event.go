package model

type EventId struct {
	TxDigest string `json:"txDigest"`
	EventSeq string `json:"eventSeq"`
}

type EventResponse struct {
	Id                EventId                `json:"id"`
	PackageId         string                 `json:"packageId"`
	TransactionModule string                 `json:"transactionModule"`
	Sender            string                 `json:"sender"`
	Type              string                 `json:"type"`
	ParsedJson        map[string]interface{} `json:"parsedJson"`
	Bcs               string                 `json:"bcs"`
	TimestampMs       string                 `json:"timestampMs"`
}
