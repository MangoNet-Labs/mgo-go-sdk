package request

type MgoXGetAllBalanceRequest struct {
	// the owner's Mgo address
	Owner string `json:"owner"`
}

type MgoXGetAllCoinsRequest struct {
	// the owner's Mgo address
	Owner string `json:"owner"`
	// optional paging cursor
	Cursor interface{} `json:"cursor"`
	// maximum number of items per page
	Limit uint64 `json:"limit" validate:"lte=50"`
}

type MgoXGetCoinMetadataRequest struct {
	CoinType string `json:"coinType"`
}
type MgoXGetBalanceRequest struct {
	// the owner's Mgo address
	Owner string `json:"owner"`
	// optional type name for the coin (e.g., 0x168da5bf1f48dafc111b0a488fa454aca95e0b5e::usdc::USDC), default to 0x2::mgo::MGO if not specified.
	CoinType string `json:"coinType"`
}

type MgoXGetCoinsRequest struct {
	// the owner's Mgo address
	Owner string `json:"owner"`
	// optional type name for the coin (e.g., 0x168da5bf1f48dafc111b0a488fa454aca95e0b5e::usdc::USDC), default to 0x2::mgo::MGO if not specified.
	CoinType string `json:"coin_type"`
	// optional paging cursor
	Cursor interface{} `json:"cursor"`
	// maximum number of items per page
	Limit uint64 `json:"limit" validate:"lte=50"`
}

type MgoXGetTotalSupplyRequest struct {
	// type name for the coin (e.g., 0x168da5bf1f48dafc111b0a488fa454aca95e0b5e::usdc::USDC)
	CoinType string `json:"coinType"`
}
