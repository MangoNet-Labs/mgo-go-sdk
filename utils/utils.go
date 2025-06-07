package utils

import (
	"strings"

	"github.com/mangonet-labs/mgo-go-sdk/model"
)

func NormalizeMgoAddress(input string) model.MgoAddress {
	addr := strings.ToLower(string(input))
	if strings.HasPrefix(addr, "0x") {
		addr = addr[2:]
	}

	addr = strings.Repeat("0", 64-len(addr)) + addr
	return model.MgoAddress("0x" + addr)
}

func IsValidMgoAddress(addr model.MgoAddress) bool {
	addr = NormalizeMgoAddress(string(addr))
	return len(addr) == 66 && strings.HasPrefix(string(addr), "0x")
}
