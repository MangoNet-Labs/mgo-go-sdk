package transaction

import (
	"encoding/hex"

	"github.com/mangonet-labs/mgo-go-sdk/model"
	"github.com/mangonet-labs/mgo-go-sdk/utils"
	"github.com/mr-tron/base58"
)

func ConvertMgoAddressStringToBytes(address model.MgoAddress) (*model.MgoAddressBytes, error) {
	normalized := utils.NormalizeMgoAddress(string(address))
	decoded, err := hex.DecodeString(string(normalized[2:]))
	if err != nil {
		return nil, err
	}
	if len(decoded) != 32 {
		return nil, ErrInvalidMgoAddress
	}

	var fixedBytes [32]byte
	copy(fixedBytes[:], decoded)

	return (*model.MgoAddressBytes)(&fixedBytes), nil
}

func ConvertMgoAddressBytesToString(addr model.MgoAddressBytes) string {
	return "0x" + hex.EncodeToString(addr[:])
}

func ConvertObjectDigestStringToBytes(digest model.ObjectDigest) (*model.ObjectDigestBytes, error) {
	decoded, err := base58.Decode(string(digest))
	if err != nil {
		return nil, err
	}
	if len(decoded) != 32 {
		return nil, ErrInvalidObjectId
	}

	return (*model.ObjectDigestBytes)(&decoded), nil
}

func ConvertObjectDigestBytesToString(digest model.ObjectDigestBytes) model.ObjectDigest {
	return model.ObjectDigest(base58.Encode(digest))
}

// ConvertBytesToMgoAddressBytes converts a byte slice to MgoAddressBytes
// The input must be exactly 32 bytes long
func ConvertBytesToMgoAddressBytes(bytes []byte) (*model.MgoAddressBytes, error) {
	if len(bytes) != 32 {
		return nil, ErrInvalidMgoAddress
	}

	var fixedBytes [32]byte
	copy(fixedBytes[:], bytes)

	return (*model.MgoAddressBytes)(&fixedBytes), nil
}
