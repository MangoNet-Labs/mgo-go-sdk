package wsconn

type SubscriptionResp struct {
	Jsonrpc string `json:"jsonrpc" yaml:"jsonrpc"`
	Result  int64  `json:"result"  yaml:"result"`
	Id      int64  `json:"id"      yaml:"id"`
}
