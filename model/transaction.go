package model

type TransactionBlock struct {
	Data         TransactionBlockData `json:"data"`
	TxSignatures []string             `json:"txSignatures"`
}

type TransactionBlockData struct {
	MessageVersion string               `json:"messageVersion"`
	Transaction    TransactionBlockKind `json:"transaction"`
	Sender         string               `json:"sender"`
	GasData        GasData              `json:"gasData"`
}

type TransactionBlockKind struct {
	Kind         string        `json:"kind"`
	Inputs       []CallArg     `json:"inputs"`
	Transactions []Transaction `json:"transactions"`
}

type CallArg map[string]interface{}

type Transaction struct {
	MakeMoveVec     []interface{}        `json:"MakeMoveVec,omitempty"`
	MergeCoins      []interface{}        `json:"MergeCoins,omitempty"`
	SplitCoins      []interface{}        `json:"SplitCoins,omitempty"`
	TransferObjects []interface{}        `json:"TransferObjects,omitempty"`
	Publish         []interface{}        `json:"Publish,omitempty"`
	Upgrade         []interface{}        `json:"Upgrade,omitempty"`
	MoveCall        *MoveCallTransaction `json:"MoveCall,omitempty"`
}

type MoveCallTransaction struct {
	Package       string   `json:"package"`
	Module        string   `json:"module"`
	Function      string   `json:"function"`
	TypeArguments []string `json:"type_arguments"`
	Arguments     []struct {
		Input  int `json:"Input"`
		Result int `json:"Result"`
	} `json:"arguments"`
}

type GasData struct {
	Payment []ObjectRef `json:"payment"`
	Owner   string      `json:"owner"`
	Price   string      `json:"price"`
	Budget  string      `json:"budget"`
}

type TxnMetaData struct {
	Gas          []MgoObjectRef `json:"gas"`
	InputObjects []interface{}  `json:"inputObjects"`
	TxBytes      string         `json:"txBytes"`
}
