package transaction

import (
	"bytes"

	"github.com/mangonet-labs/mgo-go-sdk/bcs"
	"github.com/mangonet-labs/mgo-go-sdk/model"
)

type TransactionData struct {
	V1 *TransactionDataV1
}

func (*TransactionData) IsBcsEnum() {}

func (td *TransactionData) Marshal() ([]byte, error) {
	bcsEncodedMsg := bytes.Buffer{}
	bcsEncoder := bcs.NewEncoder(&bcsEncodedMsg)
	err := bcsEncoder.Encode(td)
	if err != nil {
		return nil, err
	}

	return bcsEncodedMsg.Bytes(), nil
}

type TransactionDataV1 struct {
	Kind       *TransactionKind
	Sender     *model.MgoAddressBytes
	GasData    *GasData
	Expiration *TransactionExpiration `bcs:"optional" json:"expiration"`
}

func (td *TransactionDataV1) AddCommand(command Command) (index uint16) {
	index = uint16(len(td.Kind.ProgrammableTransaction.Commands))
	td.Kind.ProgrammableTransaction.Commands = append(td.Kind.ProgrammableTransaction.Commands, &command)

	return index
}

func (td *TransactionDataV1) AddInput(input CallArg) Argument {
	index := uint16(len(td.Kind.ProgrammableTransaction.Inputs))
	td.Kind.ProgrammableTransaction.Inputs = append(td.Kind.ProgrammableTransaction.Inputs, &input)

	return Argument{
		Input: &index,
	}
}

func (td *TransactionDataV1) GetInputObjectIndex(address model.MgoAddress) *uint16 {
	addressBytes, err := ConvertMgoAddressStringToBytes(address)
	if err != nil {
		return nil
	}

	for i, input := range td.Kind.ProgrammableTransaction.Inputs {
		if input.Object == nil {
			continue
		}

		if input.Object.ImmOrOwnedObject != nil {
			objectId := input.Object.ImmOrOwnedObject.ObjectId
			if objectId.IsEqual(*addressBytes) {
				index := uint16(i)
				return &index
			}
		}
		if input.Object.SharedObject != nil {
			objectId := input.Object.SharedObject.ObjectId
			if objectId.IsEqual(*addressBytes) {
				index := uint16(i)
				return &index
			}
		}
		if input.Object.Receiving != nil {
			objectId := input.Object.Receiving.ObjectId
			if objectId.IsEqual(*addressBytes) {
				index := uint16(i)
				return &index
			}
		}
	}

	return nil
}

type GasData struct {
	Payment *[]MgoObjectRef
	Owner   *model.MgoAddressBytes
	Price   *uint64
	Budget  *uint64
}

func (gd *GasData) IsAllSet() bool {
	if gd.Payment == nil || gd.Owner == nil || gd.Price == nil || gd.Budget == nil {
		return false
	}

	return true
}

type TransactionExpiration struct {
	None  any
	Epoch *uint64
}

func (*TransactionExpiration) IsBcsEnum() {}

type ProgrammableTransaction struct {
	Inputs   []*CallArg
	Commands []*Command
}

type TransactionKind struct {
	ProgrammableTransaction *ProgrammableTransaction
	ChangeEpoch             any
	Genesis                 any
	ConsensusCommitPrologue any
}

func (*TransactionKind) IsBcsEnum() {}

func (tk *TransactionKind) Marshal() ([]byte, error) {
	bcsEncodedMsg := bytes.Buffer{}
	bcsEncoder := bcs.NewEncoder(&bcsEncodedMsg)
	err := bcsEncoder.Encode(tk)
	if err != nil {
		return nil, err
	}

	return bcsEncodedMsg.Bytes(), nil
}

type CallArg struct {
	Pure             *Pure
	Object           *ObjectArg
	UnresolvedPure   *UnresolvedPure
	UnresolvedObject *UnresolvedObject
}

func (*CallArg) IsBcsEnum() {}

type Pure struct {
	Bytes []byte
}

type UnresolvedPure struct {
	Value any
}

type UnresolvedObject struct {
	ObjectId model.MgoAddressBytes
}

type ObjectArg struct {
	ImmOrOwnedObject *MgoObjectRef
	SharedObject     *SharedObjectRef
	Receiving        *MgoObjectRef
}

func (ObjectArg) IsBcsEnum() {}

type Command struct {
	MoveCall        *ProgrammableMoveCall
	TransferObjects *TransferObjects
	SplitCoins      *SplitCoins
	MergeCoins      *MergeCoins
	Publish         *Publish
	MakeMoveVec     *MakeMoveVec
	Upgrade         *Upgrade
}

func (*Command) IsBcsEnum() {}

type ProgrammableMoveCall struct {
	Package       model.MgoAddressBytes
	Module        string
	Function      string
	TypeArguments []*TypeTag
	Arguments     []*Argument
}

type TransferObjects struct {
	Objects []*Argument
	Address *Argument
}

type SplitCoins struct {
	Coin   *Argument
	Amount []*Argument
}

type MergeCoins struct {
	Destination *Argument
	Sources     []*Argument
}

type Publish struct {
	Modules      []model.MgoAddressBytes
	Dependencies []model.MgoAddressBytes
}

type MakeMoveVec struct {
	Type     *string
	Elements []*Argument
}

type Upgrade struct {
	Modules      []model.MgoAddressBytes
	Dependencies []model.MgoAddressBytes
	Package      model.MgoAddressBytes
	Ticket       *Argument
}

type Argument struct {
	GasCoin      any
	Input        *uint16
	Result       *uint16
	NestedResult *NestedResult
}

func (*Argument) IsBcsEnum() {}

type NestedResult struct {
	Index       uint16
	ResultIndex uint16
}

type MgoObjectRef struct {
	ObjectId model.MgoAddressBytes
	Version  uint64
	Digest   model.ObjectDigestBytes
}

type SharedObjectRef struct {
	ObjectId             model.MgoAddressBytes
	InitialSharedVersion uint64
	Mutable              bool
}

type StructTag struct {
	Address    model.MgoAddressBytes
	Module     string
	Name       string
	TypeParams []*TypeTag
}

type TypeTag struct {
	Bool    *bool
	U8      *bool
	U128    *bool
	U256    *bool
	Address *bool
	Signer  *bool
	Vector  *TypeTag
	Struct  *StructTag
	U16     *bool
	U32     *bool
	U64     *bool
}

func (*TypeTag) IsBcsEnum() {}
