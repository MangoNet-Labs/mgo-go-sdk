package request

type GetMoveFunctionArgTypesRequest struct {
	Package  string
	Module   string
	Function string
}

type GetNormalizedMoveModulesByPackageRequest struct {
	Package string `json:"package"`
}

type GetNormalizedMoveModuleRequest struct {
	Package    string `json:"package"`
	ModuleName string `json:"moduleName"`
}

type GetNormalizedMoveStructRequest struct {
	Package    string `json:"package"`
	ModuleName string `json:"moduleName"`
	StructName string `json:"structName"`
}

type GetNormalizedMoveFunctionRequest struct {
	Package      string `json:"package"`
	ModuleName   string `json:"moduleName"`
	FunctionName string `json:"functionName"`
}
