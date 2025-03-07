package request

type RPCTransactionRequestParams struct {
	MoveCallRequestParams       *MoveCallRequest       `json:"moveCallRequestParams,omitempty"`
	TransferObjectRequestParams *TransferObjectRequest `json:"transferObjectRequestParams,omitempty"`
}

type BatchTransactionRequest struct {
	// the transaction signer's Mgo address
	Signer string `json:"signer"`
	// list of transaction request parameters
	RPCTransactionRequestParams []RPCTransactionRequestParams `json:"RPCTransactionRequestParams"`
	// gas object to be used in this transaction, node will pick one from the signer's possession if not provided
	Gas string `json:"gas"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget"`
	// Whether this is a regular transaction or a Dev Inspect Transaction
	// The optional enumeration values are: `DevInspect`, or `Commit`
	MgoTransactionBlockBuilderMode string `json:"mgoTransactionBlockBuilderMode"`
}

type PayRequest struct {
	// the transaction signer's Mgo address
	Signer      string   `json:"signer"`
	MgoObjectId []string `json:"mgoObjectId"`
	Recipient   []string `json:"recipient"`
	Amount      []string `json:"amount"`
	// gas object to be used in this transaction, node will pick one from the signer's possession if not provided
	Gas string `json:"gas"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget"`
}

type PayAllMgoRequest struct {
	// the transaction signer's Mgo address
	Signer      string   `json:"signer"`
	MgoObjectId []string `json:"mgoObjectId"`
	Recipient   string   `json:"recipient"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget"`
}

type PayMgoRequest struct {
	// the transaction signer's Mgo address
	Signer      string   `json:"signer"`
	MgoObjectId []string `json:"mgoObjectId"`
	Recipient   []string `json:"recipient"`
	Amount      []string `json:"amount"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget"`
}

type PublishRequest struct {
	// the transaction signer's Mgo address
	Sender          string   `json:"sender"`
	CompiledModules []string `json:"compiled_modules"`
	Dependencies    []string `json:"dependencies"`
	// gas object to be used in this transaction, node will pick one from the signer's possession if not provided
	Gas string `json:"gas"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget"`
}

type AddStakeRequest struct {
	// the transaction signer's Mgo address
	Signer string `json:"signer"`
	// Coin<MGO> object to stake
	Coins []string `json:"coins"`
	// stake amount
	Amount string `json:"amount"`
	// the validator's Mgo address
	Validator string `json:"validator"`
	// gas object to be used in this transaction, node will pick one from the signer's possession if not provided
	Gas string `json:"gas"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget"`
}

type WithdrawStakeRequest struct {
	// the transaction signer's Mgo address
	Signer string `json:"signer"`
	// StakedMgo object ID
	StakedObjectId string `json:"stakedObjectId"`
	// gas object to be used in this transaction, node will pick one from the signer's possession if not provided
	Gas string `json:"gas"`
	// the gas budget, the transaction will fail if the gas cost exceed the budget
	GasBudget string `json:"gasBudget"`
}
