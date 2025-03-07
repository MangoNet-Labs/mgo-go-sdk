package response

type CoinLockedBalance struct {
	EpochId int `json:"epochId"`
	Number  int `json:"number"`
}

type CoinBalanceResponse struct {
	CoinType        string            `json:"coinType"`
	CoinObjectCount int               `json:"coinObjectCount"`
	TotalBalance    string            `json:"totalBalance"`
	LockedBalance   CoinLockedBalance `json:"lockedBalance"`
}

type CoinAllBalanceResponse []CoinBalanceResponse

type PaginatedCoinsResponse struct {
	Data        []CoinData `json:"data"`
	NextCursor  string     `json:"nextCursor"`
	HasNextPage bool       `json:"hasNextPage"`
}

type CoinData struct {
	CoinType            string `json:"coinType"`
	CoinObjectId        string `json:"coinObjectId"`
	Version             string `json:"version"`
	Digest              string `json:"digest"`
	Balance             string `json:"balance"`
	LockedUntilEpoch    uint64 `json:"lockedUntilEpoch"`
	PreviousTransaction string `json:"previousTransaction"`
}

type CoinMetadataResponse struct {
	Id          string `json:"id"`
	Decimals    int    `json:"decimals"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	IconUrl     string `json:"iconUrl"`
}

type TotalSupplyResponse struct {
	Value string `json:"value"`
}
