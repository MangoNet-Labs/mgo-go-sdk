package model

type MgoObjectRef struct {
	Digest   string `json:"digest"`
	ObjectId string `json:"objectId"`
	Version  uint64 `json:"version"`
}
