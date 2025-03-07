package request

type MgoXGetStakesRequest struct {
	Owner string `json:"owner"`
}
type MgoXGetStakesByIdsRequest struct {
	StakedMgoIds []string `json:"stakedMgoIds"`
}

type MgoXGetCommitteeInfoRequest struct {
	Epoch string `json:"epoch"`
}
