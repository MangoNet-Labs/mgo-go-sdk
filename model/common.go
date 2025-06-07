package model

import (
	"reflect"

	"github.com/mangonet-labs/mgo-go-sdk/bcs"
)

type MgoAddress string
type MgoAddressBytes [32]byte
type TransactionDigest string
type ObjectDigest string
type ObjectDigestBytes []byte

func init() {
	var mgoAddressBytes MgoAddressBytes
	if reflect.ValueOf(mgoAddressBytes).Type().Name() != bcs.MgoAddressBytesName {
		panic("MgoAddressBytes type name not match")
	}
}

func (s MgoAddressBytes) IsEqual(other MgoAddressBytes) bool {
	for i, b := range s {
		if b != other[i] {
			return false
		}
	}
	return true
}

func (o ObjectDigestBytes) IsEqual(other ObjectDigestBytes) bool {
	if len(o) != len(other) {
		return false
	}

	for i, b := range o {
		if b != other[i] {
			return false
		}
	}
	return true
}
