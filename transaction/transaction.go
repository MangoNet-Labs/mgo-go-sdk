package transaction

import (
	"bytes"
	"context"
	"math"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/bcs"
	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/model"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/model/response"
	"github.com/mangonet-labs/mgo-go-sdk/utils"
	"github.com/samber/lo"
)

type Transaction struct {
	Data            TransactionData
	Signer          *keypair.Keypair
	SponsoredSigner *keypair.Keypair
	MgoClient       *client.Client
}

func (tx *Transaction) GetTransactionData() (*TransactionData, error) {
	return &tx.Data, nil
}

func NewTransaction() *Transaction {
	data := TransactionData{
		V1: &TransactionDataV1{},
	}
	data.V1.Kind = &TransactionKind{
		ProgrammableTransaction: &ProgrammableTransaction{},
	}
	data.V1.GasData = &GasData{}

	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) SetSigner(signer *keypair.Keypair) *Transaction {
	tx.Signer = signer

	return tx
}

func (tx *Transaction) SetSponsoredSigner(signer *keypair.Keypair) *Transaction {
	tx.SponsoredSigner = signer

	return tx
}

func (tx *Transaction) SetMgoClient(client *client.Client) *Transaction {
	tx.MgoClient = client

	return tx
}

func (tx *Transaction) SetTransactionData(data *TransactionData) *Transaction {
	tx.Data = *data

	return tx
}

func (tx *Transaction) SetSender(sender model.MgoAddress) *Transaction {
	address := utils.NormalizeMgoAddress(string(sender))
	addressBytes, err := ConvertMgoAddressStringToBytes(address)
	if err != nil {
		panic(err)
	}
	tx.Data.V1.Sender = addressBytes

	return tx
}

func (tx *Transaction) SetSenderIfNotSet(sender model.MgoAddress) *Transaction {
	if tx.Data.V1.Sender == nil {
		tx.SetSender(sender)
	}

	return tx
}

func (tx *Transaction) SetExpiration(expiration TransactionExpiration) *Transaction {
	tx.Data.V1.Expiration = &expiration

	return tx
}

func (tx *Transaction) SetGasPayment(payment []MgoObjectRef) *Transaction {
	tx.Data.V1.GasData.Payment = &payment

	return tx
}

func (tx *Transaction) SetGasOwner(owner model.MgoAddress) *Transaction {
	addressBytes, err := ConvertMgoAddressStringToBytes(owner)
	if err != nil {
		panic(err)
	}
	tx.Data.V1.GasData.Owner = addressBytes

	return tx
}

func (tx *Transaction) SetGasPrice(price uint64) *Transaction {
	tx.Data.V1.GasData.Price = &price

	return tx
}

func (tx *Transaction) SetGasBudget(budget uint64) *Transaction {
	tx.Data.V1.GasData.Budget = &budget

	return tx
}

func (tx *Transaction) SetGasBudgetIfNotSet(budget uint64) *Transaction {
	if tx.Data.V1.GasData.Budget == nil {
		tx.Data.V1.GasData.Budget = &budget
	}

	return tx
}

func (tx *Transaction) Gas() Argument {
	return Argument{
		GasCoin: struct{}{},
	}
}

func (tx *Transaction) Add(command Command) Argument {
	index := tx.Data.V1.AddCommand(command)

	return createTransactionResult(index, nil)
}

func (tx *Transaction) SplitCoins(coin Argument, amount []Argument) Argument {
	return tx.Add(splitCoins(SplitCoins{
		Coin:   &coin,
		Amount: convertArgumentsToArgumentPtrs(amount),
	}))
}

func (tx *Transaction) MergeCoins(destination Argument, sources []Argument) Argument {
	return tx.Add(mergeCoins(MergeCoins{
		Destination: &destination,
		Sources:     convertArgumentsToArgumentPtrs(sources),
	}))
}

func (tx *Transaction) Publish(modules []model.MgoAddress, dependencies []model.MgoAddress) Argument {
	moduleAddress := make([]model.MgoAddressBytes, len(modules))
	for i, module := range modules {
		v, err := ConvertMgoAddressStringToBytes(module)
		if err != nil {
			panic(err)
		}
		moduleAddress[i] = *v
	}

	dependenciesAddress := make([]model.MgoAddressBytes, len(dependencies))
	for i, dependency := range dependencies {
		v, err := ConvertMgoAddressStringToBytes(dependency)
		if err != nil {
			panic(err)
		}
		dependenciesAddress[i] = *v
	}

	return tx.Add(publish(Publish{
		Modules:      moduleAddress,
		Dependencies: dependenciesAddress,
	}))
}

func (tx *Transaction) Upgrade(
	modules []model.MgoAddress,
	dependencies []model.MgoAddress,
	packageId model.MgoAddress,
	ticket Argument,
) Argument {
	moduleAddress := make([]model.MgoAddressBytes, len(modules))
	for i, module := range modules {
		v, err := ConvertMgoAddressStringToBytes(module)
		if err != nil {
			panic(err)
		}
		moduleAddress[i] = *v
	}

	dependenciesAddress := make([]model.MgoAddressBytes, len(dependencies))
	for i, dependency := range dependencies {
		v, err := ConvertMgoAddressStringToBytes(dependency)
		if err != nil {
			panic(err)
		}
		dependenciesAddress[i] = *v
	}

	packageIdBytes, err := ConvertMgoAddressStringToBytes(packageId)
	if err != nil {
		panic(err)
	}

	return tx.Add(upgrade(Upgrade{
		Modules:      moduleAddress,
		Dependencies: dependenciesAddress,
		Package:      *packageIdBytes,
		Ticket:       &ticket,
	}))
}

func (tx *Transaction) MoveCall(
	packageId model.MgoAddress,
	module string,
	function string,
	typeArguments []TypeTag,
	arguments []Argument,
) Argument {
	packageIdBytes, err := ConvertMgoAddressStringToBytes(packageId)
	if err != nil {
		panic(err)
	}

	return tx.Add(moveCall(ProgrammableMoveCall{
		Package:       *packageIdBytes,
		Module:        module,
		Function:      function,
		TypeArguments: convertTypeTagsToTypeTagPtrs(typeArguments),
		Arguments:     convertArgumentsToArgumentPtrs(arguments),
	}))
}

func (tx *Transaction) TransferObjects(objects []Argument, address Argument) Argument {
	return tx.Add(transferObjects(TransferObjects{
		Objects: convertArgumentsToArgumentPtrs(objects),
		Address: &address,
	}))
}

func (tx *Transaction) MakeMoveVec(typeValue *string, elements []Argument) Argument {
	return tx.Add(makeMoveVec(MakeMoveVec{
		Type:     typeValue,
		Elements: convertArgumentsToArgumentPtrs(elements),
	}))
}

func (tx *Transaction) Object(input any) Argument {

	if s, ok := input.(string); ok {
		if utils.IsValidMgoAddress(model.MgoAddress(s)) {
			address := utils.NormalizeMgoAddress(s)
			addressBytes, err := ConvertMgoAddressStringToBytes(address)
			if err != nil {
				panic(err)
			}

			arg := tx.Data.V1.AddInput(CallArg{
				UnresolvedObject: &UnresolvedObject{
					ObjectId: *addressBytes,
				},
			})

			return arg
		} else {
			panic(ErrObjectNotSupportType)
		}
	}

	if arg, ok := input.(Argument); ok {
		return arg
	}

	if v, ok := input.(CallArg); ok {
		isTypeSupported := false

		if v.Object.SharedObject != nil {

			address := ConvertMgoAddressBytesToString(v.Object.SharedObject.ObjectId)
			if index := tx.Data.V1.GetInputObjectIndex(model.MgoAddress(address)); index != nil {
				if v.Object.SharedObject.Mutable {
					newExistObject := tx.Data.V1.Kind.ProgrammableTransaction.Inputs[*index]
					if newExistObject.Object.SharedObject != nil {
						newExistObject.Object.SharedObject.Mutable = true
						tx.Data.V1.Kind.ProgrammableTransaction.Inputs[*index] = newExistObject
					}
				}

				return Argument{
					Input: index,
				}
			}

			isTypeSupported = true
		}

		if v.Object.ImmOrOwnedObject != nil {
			isTypeSupported = true
		}
		if v.Object.Receiving != nil {
			isTypeSupported = true
		}

		if isTypeSupported {
			arg := tx.Data.V1.AddInput(CallArg{
				Object: v.Object,
			})
			return arg
		}
	}

	panic(ErrObjectNotSupportType)
}

func (tx *Transaction) Pure(input any) Argument {
	var val []byte
	if s, ok := input.(string); ok && utils.IsValidMgoAddress(model.MgoAddress(s)) {
		fixedAddressBytes, err := ConvertMgoAddressStringToBytes(model.MgoAddress(s))
		if err != nil {
			panic(err)
		}
		addressBytes := fixedAddressBytes[:]
		val = addressBytes
	} else {
		bcsEncodedMsg := bytes.Buffer{}
		bcsEncoder := bcs.NewEncoder(&bcsEncodedMsg)
		err := bcsEncoder.Encode(input)
		if err != nil {
			panic(err)
		}
		val = bcsEncodedMsg.Bytes()
	}

	arg := tx.Data.V1.AddInput(CallArg{Pure: &Pure{
		Bytes: val,
	}})

	return arg
}

func (tx *Transaction) Execute(
	ctx context.Context,
	options request.MgoTransactionBlockOptions,
	requestType string,
) (*response.MgoTransactionBlockResponse, error) {
	if tx.MgoClient == nil {
		return nil, ErrMgoClientNotSet
	}
	req, err := tx.ToMgoExecuteTransactionBlockRequest(ctx, options, requestType)
	if err != nil {
		return nil, err
	}
	rsp, err := tx.MgoClient.MgoExecuteTransactionBlock(ctx, *req)
	if err != nil {
		return nil, err
	}

	return &rsp, nil
}

func (tx *Transaction) ToMgoExecuteTransactionBlockRequest(
	ctx context.Context,
	options request.MgoTransactionBlockOptions,
	requestType string,
) (*request.MgoExecuteTransactionBlockRequest, error) {
	if tx.Signer == nil {
		return nil, ErrSignerNotSet
	}

	b64TxBytes, err := tx.BuildTransaction(ctx)
	if err != nil {
		return nil, err
	}
	var signatures []string
	if tx.SponsoredSigner != nil {
		sponsoredMessage, err := tx.SponsoredSigner.SignTransactionBlock(&model.TxnMetaData{
			TxBytes: b64TxBytes,
		})
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, sponsoredMessage.Signature)
	}
	message, err := tx.Signer.SignTransactionBlock(&model.TxnMetaData{
		TxBytes: b64TxBytes,
	})
	if err != nil {
		return nil, err
	}
	signatures = append(signatures, message.Signature)

	return &request.MgoExecuteTransactionBlockRequest{
		TxBytes:     b64TxBytes,
		Signature:   signatures,
		Options:     options,
		RequestType: requestType,
	}, nil
}

func (tx *Transaction) BuildTransaction(ctx context.Context) (string, error) {
	if tx.Signer == nil {
		return "", ErrSignerNotSet
	}

	if tx.Data.V1.GasData.Price == nil {
		if tx.MgoClient != nil {
			rsp, err := tx.MgoClient.MgoXGetReferenceGasPrice(ctx)
			if err != nil {
				return "", err
			}
			tx.SetGasPrice(rsp)
		}
	}
	tx.SetGasBudgetIfNotSet(defaultGasBudget)
	tx.SetSenderIfNotSet(model.MgoAddress(tx.Signer.MgoAddress()))

	return tx.Build(false)
}

func (tx *Transaction) Build(onlyTransactionKind bool) (string, error) {
	if onlyTransactionKind {
		bcsEncodedMsg, err := tx.Data.V1.Kind.Marshal()
		if err != nil {
			return "", err
		}
		bcsBase64 := bcs.ToBase64(bcsEncodedMsg)
		return bcsBase64, nil
	}

	if tx.Data.V1.Sender == nil {
		return "", ErrSenderNotSet
	}
	if tx.Data.V1.GasData.Owner == nil {
		tx.SetGasOwner(model.MgoAddress(tx.Signer.MgoAddress()))
	}
	if !tx.Data.V1.GasData.IsAllSet() {
		return "", ErrGasDataNotAllSet
	}

	bcsEncodedMsg, err := tx.Data.Marshal()
	if err != nil {
		return "", err
	}
	bcsBase64 := bcs.ToBase64(bcsEncodedMsg)

	return bcsBase64, nil
}

func (tx *Transaction) NewTransactionFromKind() (newTx *Transaction, err error) {
	newTx = NewTransaction()
	err = copier.CopyWithOption(&newTx.Data.V1.Kind, &tx.Data.V1.Kind, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, err
	}
	return newTx, nil
}

func NewMgoObjectRef(objectId model.MgoAddress, version string, digest model.ObjectDigest) (*MgoObjectRef, error) {
	objectIdBytes, err := ConvertMgoAddressStringToBytes(objectId)
	if err != nil {
		return nil, err
	}
	digestBytes, err := ConvertObjectDigestStringToBytes(digest)
	if err != nil {
		return nil, err
	}
	versionUint64, err := strconv.ParseUint(version, 10, 64)
	if err != nil {
		return nil, err
	}

	return &MgoObjectRef{
		ObjectId: *objectIdBytes,
		Version:  versionUint64,
		Digest:   *digestBytes,
	}, nil
}

func createTransactionResult(index uint16, length *uint16) Argument {
	if length == nil {
		length = lo.ToPtr(uint16(math.MaxUint16))
	}

	return Argument{
		Result: lo.ToPtr(index),
	}
}

func convertArgumentsToArgumentPtrs(args []Argument) []*Argument {
	argPtrs := make([]*Argument, len(args))
	for i, arg := range args {
		v := arg
		argPtrs[i] = &v
	}

	return argPtrs
}

func convertTypeTagsToTypeTagPtrs(tags []TypeTag) []*TypeTag {
	tagPtrs := make([]*TypeTag, len(tags))
	for i, tag := range tags {
		v := tag
		tagPtrs[i] = &v
	}

	return tagPtrs
}
