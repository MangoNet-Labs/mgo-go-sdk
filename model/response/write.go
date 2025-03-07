package response

import "github.com/MangoNet-Labs/mgo-go-sdk/model"

type MgoTransactionBlockResponse struct {
	Digest                  string                 `json:"digest"                            yaml:"digest"`
	Transaction             model.TransactionBlock `json:"transaction,omitempty"             yaml:"transaction"`
	RawTransaction          string                 `json:"rawTransaction,omitempty"          yaml:"rawTransaction"`
	Effects                 model.Effects          `json:"effects,omitempty"                 yaml:"effects"`
	Events                  []model.EventResponse  `json:"events,omitempty"                  yaml:"events"`
	ObjectChanges           []model.ObjectChange   `json:"objectChanges,omitempty"           yaml:"objectChanges"`
	BalanceChanges          []model.BalanceChanges `json:"balanceChanges,omitempty"          yaml:"balanceChanges"`
	TimestampMs             string                 `json:"timestampMs,omitempty"             yaml:"timestampMs"`
	Checkpoint              string                 `json:"checkpoint,omitempty"              yaml:"checkpoint"`
	ConfirmedLocalExecution bool                   `json:"confirmedLocalExecution,omitempty" yaml:"confirmedLocalExecution"`
}

type TransactionFilter map[string]interface{}

type MgoTransactionBlockOptions struct {
	ShowInput          bool `json:"showInput,omitempty"          yaml:"showInput"`
	ShowRawInput       bool `json:"showRawInput,omitempty"       yaml:"showRawInput"`
	ShowEffects        bool `json:"showEffects,omitempty"        yaml:"showEffects"`
	ShowEvents         bool `json:"showEvents,omitempty"         yaml:"showEvents"`
	ShowObjectChanges  bool `json:"showObjectChanges,omitempty"  yaml:"showObjectChanges"`
	ShowBalanceChanges bool `json:"showBalanceChanges,omitempty" yaml:"showBalanceChanges"`
}

type MgoTransactionBlockResponseQuery struct {
	TransactionFilter TransactionFilter          `json:"filter"  yaml:"transactionFilter"`
	Options           MgoTransactionBlockOptions `json:"options" yaml:"options"`
}
