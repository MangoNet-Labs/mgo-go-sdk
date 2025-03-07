package request

type MgoGetEventsRequest struct {
	Digest string `json:"digest"`
}
type MgoXQueryEventsRequest struct {
	// the event query criteria. See Event filter documentation[https://docs.mangonet.io/mango-api-ref/#mgox_queryevents] for examples.
	MgoEventFilter interface{} `json:"mgoEventFilter"`
	// optional paging cursor
	Cursor interface{} `json:"cursor"`
	// maximum number of items per page
	Limit uint64 `json:"limit" validate:"lte=50"`
	// query result ordering, default to false (ascending order), oldest record first
	DescendingOrder bool `json:"descendingOrder"`
}
