package response

type MgoEventResponse struct {
	Id                EventId                `json:"id"`
	PackageId         string                 `json:"packageId"`
	TransactionModule string                 `json:"transactionModule"`
	Sender            string                 `json:"sender"`
	Type              string                 `json:"type"`
	ParsedJson        map[string]interface{} `json:"parsedJson"`
	Bcs               string                 `json:"bcs"`
	TimestampMs       string                 `json:"timestampMs"`
}

type GetEventsResponse []*MgoEventResponse

type PaginatedEventsResponse struct {
	Data        []MgoEventResponse `json:"data"`
	NextCursor  EventId            `json:"nextCursor"`
	HasNextPage bool               `json:"hasNextPage"`
}
