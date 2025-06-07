package transaction

import "errors"

var (
	ErrSignerNotSet         = errors.New("signer not set")
	ErrSenderNotSet         = errors.New("sender not set")
	ErrMgoClientNotSet      = errors.New("mgo client not set")
	ErrGasDataNotAllSet     = errors.New("gas data not all set")
	ErrInvalidMgoAddress    = errors.New("invalid mgo address")
	ErrInvalidObjectId      = errors.New("invalid object id")
	ErrObjectNotSupportType = errors.New("object not support type")
)
