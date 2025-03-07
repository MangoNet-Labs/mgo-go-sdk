package response

type MgoMoveNormalizedStruct struct {
	Abilities      interface{}   `json:"abilities"`
	TypeParameters []interface{} `json:"typeParameters"`
	Fields         []interface{} `json:"fields"`
}

type MgoMoveNormalizedFunction struct {
	Visibility interface{}   `json:"visibility"`
	IsEntry    bool          `json:"isEntry"`
	Parameters []interface{} `json:"parameters"`
	Return_    []interface{} `json:"return_"`
}

type GetMoveFunctionArgTypesResponse []interface{}

type GetNormalizedMoveModulesByPackageResponse map[string]MgoMoveNormalizedModule

type GetNormalizedMoveModuleResponse MgoMoveNormalizedModule

type GetNormalizedMoveStructResponse MgoMoveNormalizedStruct

type GetNormalizedMoveFunctionResponse MgoMoveNormalizedFunction
