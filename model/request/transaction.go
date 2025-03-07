package request

import (
	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/model"
)

type TransferMgoRequest struct {
	// the transaction signer's Mgo address
	Signer      string `json:"signer"      yaml:"signer"`
	MgoObjectId string `json:"mgoObjectId" yaml:"mgoObjectId"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget" yaml:"gasBudget"`
	Recipient string `json:"recipient" yaml:"recipient"`
	Amount    string `json:"amount"    yaml:"amount"`
}

type TransferObjectRequest struct {
	// the transaction signer's Mgo address
	Signer   string `json:"signer"   yaml:"signer"`
	ObjectId string `json:"objectId" yaml:"objectId"`
	// gas object to be used in this transaction, node will pick one from the signer's possession if not provided
	Gas string `json:"gas" yaml:"gas"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget" yaml:"gasBudget"`
	Recipient string `json:"recipient" yaml:"recipient"`
}

type SignAndExecuteTransactionBlockRequest struct {
	TxnMetaData model.TxnMetaData
	// the address private key to sign the transaction
	Keypair *keypair.Keypair
	Options TransactionBlockOptions `json:"options" yaml:"options"`
	// The optional enumeration values are: `WaitForEffectsCert`, or `WaitForLocalExecution`
	RequestType string `json:"requestType" yaml:"requestType"`
}
type TransactionBlockOptions struct {
	ShowInput          bool `json:"showInput,omitempty"          yaml:"showInput"`
	ShowRawInput       bool `json:"showRawInput,omitempty"       yaml:"showRawInput"`
	ShowEffects        bool `json:"showEffects,omitempty"        yaml:"showEffects"`
	ShowEvents         bool `json:"showEvents,omitempty"         yaml:"showEvents"`
	ShowObjectChanges  bool `json:"showObjectChanges,omitempty"  yaml:"showObjectChanges"`
	ShowBalanceChanges bool `json:"showBalanceChanges,omitempty" yaml:"showBalanceChanges"`
}
type MoveCallRequest struct {
	// the transaction signer's Mgo address
	Signer string `json:"signer" yaml:"signer"`
	// the package containing the module and function
	PackageObjectId string `json:"packageObjectId" yaml:"packageObjectId"`
	// the specific module in the package containing the function
	Module string `json:"module" yaml:"module"`
	// the function to be called
	Function string `json:"function" yaml:"function"`
	// the type arguments to the function
	TypeArguments []interface{} `json:"typeArguments" yaml:"typeArguments"`
	// the arguments to the function
	Arguments []interface{} `json:"arguments" yaml:"arguments"`
	// gas object to be used in this transaction, node will pick one from the signer's possession if not provided
	Gas *string `json:"gas" yaml:"gas"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget" yaml:"gasBudget"`

	ExecutionMode string `json:"executionMode" yaml:"executionMode"`
}

type MgoSubscribeEventsRequest struct {
	// the event query criteria.
	MgoEventFilter interface{} `json:"mgoEventFilter" yaml:"mgoEventFilter"`
}
type MgoSubscribeTransactionsRequest struct {
	// the transaction query criteria.
	TransactionFilter interface{} `json:"filter" yaml:"transactionFilter"`
}

type MergeCoinsRequest struct {
	// the transaction signer's Mgo address
	Signer      string `json:"signer"      yaml:"signer"`
	PrimaryCoin string `json:"primaryCoin" yaml:"primaryCoin"`
	CoinToMerge string `json:"coinToMerge" yaml:"coinToMerge"`
	// gas object to be used in this transaction, node will pick one from the signer's possession if not provided
	Gas string `json:"gas" yaml:"gas"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget" yaml:"gasBudget"`
}

type SplitCoinRequest struct {
	// the transaction signer's Mgo address
	Signer       string   `json:"signer"       yaml:"signer"`
	CoinObjectId string   `json:"coinObjectId" yaml:"coinObjectId"`
	SplitAmounts []string `json:"splitAmounts" yaml:"splitAmounts"`
	// gas object to be used in this transaction, node will pick one from the signer's possession if not provided
	Gas string `json:"gas" yaml:"gas"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget" yaml:"gasBudget"`
}

type SplitCoinEqualRequest struct {
	// the transaction signer's Mgo address
	Signer       string `json:"signer"       yaml:"signer"`
	CoinObjectId string `json:"coinObjectId" yaml:"coinObjectId"`
	SplitCount   string `json:"splitCount"   yaml:"splitCount"`
	// gas object to be used in this transaction, node will pick one from the signer's possession if not provided
	Gas string `json:"gas" yaml:"gas"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget" yaml:"gasBudget"`
}

// TransactionFilterByFromAddress is a filter for from address
type TransactionFilterByFromAddress struct {
	FromAddress string `json:"FromAddress" yaml:"fromAddress"`
}

// TransactionFilterByToAddress is a filter for to address
type TransactionFilterByToAddress struct {
	ToAddress string `json:"ToAddress" yaml:"toAddress"`
}

// TransactionFilterByInputObject is a filter for input objects
type TransactionFilterByInputObject struct {
	// InputObject is the id of the object
	InputObject string `json:"InputObject" yaml:"inputObject"`
}

// TransactionFilterByChangedObjectFilter is a filter for changed objects
type TransactionFilterByChangedObjectFilter struct {
	// ChangedObject is a filter for changed objects
	ChangedObject string `json:"ChangedObject" yaml:"changedObject"`
}

type MoveFunction struct {
	Package  string  `json:"package"  yaml:"package"`
	Module   *string `json:"module"   yaml:"module"`
	Function *string `json:"function" yaml:"function"`
}

// TransactionFilterByMoveFunction is a filter for move functions
type TransactionFilterByMoveFunction struct {
	MoveFunction MoveFunction `json:"MoveFunction" yaml:"moveFunction"`
}
