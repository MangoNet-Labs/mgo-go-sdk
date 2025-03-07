package secp256k1

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/tyler-smith/go-bip32"
)

var (
	ErrInvalidPath        = errors.New("invalid derivation path")
	ErrNoPublicDerivation = errors.New("no public derivation for secp256k1")
)

func DeriveForPath(path string, seed []byte) (*bip32.Key, error) {
	if !isValidBIP32Path(path) {
		return nil, ErrInvalidPath
	}

	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}
	segments := strings.Split(path, "/")
	for _, segment := range segments[1:] {
		i64, err := strconv.ParseUint(strings.TrimRight(segment, "'"), 10, 32)
		if err != nil {
			return nil, err
		}
		i := uint32(i64)
		if strings.HasSuffix(segment, "'") {
			i = i + bip32.FirstHardenedChild
		}
		key, err = key.NewChildKey(i)
		if err != nil {
			return nil, err
		}
	}
	return key, nil
}

func isValidBIP32Path(path string) bool {
	re := regexp.MustCompile(`^m/(54|74)'/938'/\d+'/(\d+/)*\d+$`)
	return re.MatchString(path)
}
