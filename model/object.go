package model

type ObjectOwner struct {
	AddressOwner string      `json:"AddressOwner"`
	ObjectOwner  string      `json:"ObjectOwner"`
	Shared       ObjectShare `json:"Shared"`
}

type ObjectShare struct {
	InitialSharedVersion int `json:"initial_shared_version"`
}

type ObjectChange struct {
	Type            string      `json:"type"`
	Sender          string      `json:"sender"`
	Owner           ObjectOwner `json:"owner"`
	ObjectType      string      `json:"objectType"`
	ObjectId        string      `json:"objectId"`
	PackageId       string      `json:"packageId"`
	Modules         []string    `json:"modules"`
	Version         string      `json:"version"`
	PreviousVersion string      `json:"previousVersion,omitempty"`
	Digest          string      `json:"digest"`
}

type BalanceChanges struct {
	Owner    ObjectOwner `json:"owner"`
	CoinType string      `json:"coinType"`
	Amount   string      `json:"amount"`
}
