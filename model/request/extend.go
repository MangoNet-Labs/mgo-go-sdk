package request

type MgoXResolveNameServiceAddressRequest struct {
	Name string `json:"name"`
}

type MgoXResolveNameServiceNamesRequest struct {
	Address string `json:"address"`
	// optional paging cursor
	Cursor interface{} `json:"cursor"`
	// maximum number of items per page
	Limit uint64 `json:"limit" validate:"lte=50"`
}

type DynamicFieldObjectName struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type MgoXGetDynamicFieldObjectRequest struct {
	ObjectId         string                 `json:"objectId"`
	DynamicFieldName DynamicFieldObjectName `json:"dynamicFieldName"`
}

type MgoXGetDynamicFieldsRequest struct {
	ObjectId string `json:"objectId"`
	// optional paging cursor
	Cursor interface{} `json:"cursor"`
	// maximum number of items per page
	Limit uint64 `json:"limit" validate:"lte=50"`
}

type MgoObjectDataOptions struct {
	ShowType                bool `json:"showType"`
	ShowContent             bool `json:"showContent"`
	ShowBcs                 bool `json:"showBcs"`
	ShowOwner               bool `json:"showOwner"`
	ShowPreviousTransaction bool `json:"showPreviousTransaction"`
	ShowStorageRebate       bool `json:"showStorageRebate"`
	ShowDisplay             bool `json:"showDisplay"`
}

type MgoObjectResponseQuery struct {
	// If None, no filter will be applied
	Filter interface{} `json:"filter"`
	// config which fields to include in the response, by default only digest is included
	Options MgoObjectDataOptions `json:"options"`
}

type MgoXGetOwnedObjectsRequest struct {
	// the owner's Mgo address
	Address string `json:"address" validate:"checkAddress"`
	// the objects query criteria
	Query MgoObjectResponseQuery
	// optional paging cursor
	Cursor interface{} `json:"cursor"`
	// maximum number of items per page
	Limit uint64 `json:"limit" validate:"lte=50"`
}

type TransactionFilter map[string]interface{}

type MgoTransactionBlockOptions struct {
	ShowInput          bool `json:"showInput,omitempty"`
	ShowRawInput       bool `json:"showRawInput,omitempty"`
	ShowEffects        bool `json:"showEffects,omitempty"`
	ShowEvents         bool `json:"showEvents,omitempty"`
	ShowObjectChanges  bool `json:"showObjectChanges,omitempty"`
	ShowBalanceChanges bool `json:"showBalanceChanges,omitempty"`
}

type MgoTransactionBlockResponseQuery struct {
	TransactionFilter TransactionFilter          `json:"filter"`
	Options           MgoTransactionBlockOptions `json:"options"`
}

type MgoXQueryTransactionBlocksRequest struct {
	MgoTransactionBlockResponseQuery MgoTransactionBlockResponseQuery
	// optional paging cursor
	Cursor interface{} `json:"cursor"`
	// maximum number of items per page
	Limit uint64 `json:"limit" validate:"lte=50"`
	// query result ordering, default to false (ascending order), oldest record first
	DescendingOrder bool `json:"descendingOrder"`
}

type EventFilterByPackage struct {
	Package string `json:"Package"`
}

type MoveModule struct {
	Package string `json:"package"`
	Module  string `json:"module"`
}
type EventFilterByMoveModule struct {
	MoveModule MoveModule `json:"MoveModule"`
}

type EventFilterByMoveEventType struct {
	MoveEventType string `json:"MoveEventType"`
}

type EventFilterByMoveEvent struct {
	MoveEvent string `json:"MoveEvent"`
}

type MoveEventModule struct {
	Package string `json:"package"`
	Module  string `json:"module"`
	Event   string `json:"event"`
}

type EventFilterByMoveEventModule struct {
	MoveEventModule MoveEventModule `json:"MoveEventModule"`
}

type MoveEventField struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

type EventFilterByMoveEventField struct {
	MoveEventField MoveEventField `json:"MoveEventField"`
}

type EventFilterBySenderAddress struct {
	SenderAddress string `json:"SenderAddress"`
}

type EventFilterBySender struct {
	Sender string `json:"Sender"`
}

type Recipient struct {
	AddressOwner string `json:"AddressOwner"`
}
type EventFilterByRecipient struct {
	Recipient Recipient `json:"Recipient"`
}
type EventFilterByObject struct {
	Object string `json:"Object"`
}

type EventFilterByEventType struct {
	EventType string `json:"EventType"`
}

type EventFilterByTransaction struct {
	Transaction string `json:"Transaction"`
}
type TimeRange struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
type EventFilterByTimeRange struct {
	TimeRange TimeRange `json:"TimeRange"`
}
