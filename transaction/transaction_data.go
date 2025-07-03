package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mangonet-labs/mgo-go-sdk/bcs"
	"github.com/mangonet-labs/mgo-go-sdk/model"
)

type TransactionData struct {
	V0 *TransactionDataV0
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

type TransactionDataV0 struct {
	Kind       *TransactionKind
	Sender     *model.MgoAddressBytes
	GasData    *GasData
	Expiration *TransactionExpiration `bcs:"optional" json:"expiration"`
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

// SignedTransaction represents a signed transaction envelope
type SignedTransaction struct {
	TransactionData *TransactionData
	TxSignatures    [][]byte
}

// SerializedTransactionDataV1 represents the serialized format of transaction data
type SerializedTransactionDataV1 struct {
	Version      int                     `json:"version"`
	Sender       string                  `json:"sender,omitempty"`
	Expiration   map[string]interface{}  `json:"expiration,omitempty"`
	GasConfig    SerializedGasConfig     `json:"gasConfig"`
	Inputs       []SerializedInput       `json:"inputs"`
	Transactions []SerializedTransaction `json:"transactions"`
}

type SerializedGasConfig struct {
	Owner   string   `json:"owner,omitempty"`
	Budget  string   `json:"budget,omitempty"`
	Price   string   `json:"price,omitempty"`
	Payment []string `json:"payment,omitempty"`
}

type SerializedInput struct {
	Kind  string                 `json:"kind"`
	Index int                    `json:"index"`
	Value map[string]interface{} `json:"value"`
	Type  string                 `json:"type"`
}

type SerializedTransaction struct {
	Kind          string        `json:"kind"`
	Target        string        `json:"target,omitempty"`
	TypeArguments []string      `json:"typeArguments,omitempty"`
	Arguments     []interface{} `json:"arguments,omitempty"`
	Objects       []interface{} `json:"objects,omitempty"`
	Address       interface{}   `json:"address,omitempty"`
	Destination   interface{}   `json:"destination,omitempty"`
	Sources       []interface{} `json:"sources,omitempty"`
	Coin          interface{}   `json:"coin,omitempty"`
	Amounts       []interface{} `json:"amounts,omitempty"`
	Modules       [][]byte      `json:"modules,omitempty"`
	Dependencies  []string      `json:"dependencies,omitempty"`
	PackageId     string        `json:"packageId,omitempty"`
	Ticket        interface{}   `json:"ticket,omitempty"`
	Type          interface{}   `json:"type,omitempty"`
}

// SerializeV1 converts transaction data to SerializedTransactionDataV1 format
func (td *TransactionData) SerializeV1() (*SerializedTransactionDataV1, error) {
	if td.V1 == nil {
		return nil, fmt.Errorf("transaction data V1 is nil")
	}

	result := &SerializedTransactionDataV1{
		Version: 1,
	}

	// Set sender if available
	if td.V1.Sender != nil {
		result.Sender = ConvertMgoAddressBytesToString(*td.V1.Sender)
	}

	// Handle expiration
	if td.V1.Expiration != nil {
		if td.V1.Expiration.Epoch != nil {
			result.Expiration = map[string]interface{}{
				"Epoch": *td.V1.Expiration.Epoch,
			}
		} else {
			result.Expiration = map[string]interface{}{
				"None": true,
			}
		}
	}

	// Configure gas
	if td.V1.GasData != nil {
		gasConfig := SerializedGasConfig{}
		if td.V1.GasData.Owner != nil {
			gasConfig.Owner = ConvertMgoAddressBytesToString(*td.V1.GasData.Owner)
		}
		if td.V1.GasData.Budget != nil {
			gasConfig.Budget = fmt.Sprintf("%d", *td.V1.GasData.Budget)
		}
		if td.V1.GasData.Price != nil {
			gasConfig.Price = fmt.Sprintf("%d", *td.V1.GasData.Price)
		}
		if td.V1.GasData.Payment != nil {
			payment := []string{}
			for _, p := range *td.V1.GasData.Payment {
				payment = append(payment, ConvertMgoAddressBytesToString(p.ObjectId))
			}
			gasConfig.Payment = payment
		}
		result.GasConfig = gasConfig
	}

	// Process inputs
	if td.V1.Kind != nil && td.V1.Kind.ProgrammableTransaction != nil {
		inputs := []SerializedInput{}
		for i, input := range td.V1.Kind.ProgrammableTransaction.Inputs {
			serializedInput := SerializedInput{
				Kind:  "Input",
				Index: i,
				Value: map[string]interface{}{},
			}

			if input.Object != nil {
				serializedInput.Type = "object"
				if input.Object.ImmOrOwnedObject != nil {
					serializedInput.Value["Object"] = map[string]interface{}{
						"ImmOrOwned": map[string]interface{}{
							"objectId": ConvertMgoAddressBytesToString(input.Object.ImmOrOwnedObject.ObjectId),
							"version":  input.Object.ImmOrOwnedObject.Version,
							"digest":   ConvertObjectDigestBytesToString(input.Object.ImmOrOwnedObject.Digest),
						},
					}
				} else if input.Object.Receiving != nil {
					serializedInput.Value["Object"] = map[string]interface{}{
						"Receiving": map[string]interface{}{
							"objectId": ConvertMgoAddressBytesToString(input.Object.Receiving.ObjectId),
							"version":  input.Object.Receiving.Version,
							"digest":   ConvertObjectDigestBytesToString(input.Object.Receiving.Digest),
						},
					}
				} else if input.Object.SharedObject != nil {
					serializedInput.Value["Object"] = map[string]interface{}{
						"Shared": map[string]interface{}{
							"objectId":             ConvertMgoAddressBytesToString(input.Object.SharedObject.ObjectId),
							"initialSharedVersion": input.Object.SharedObject.InitialSharedVersion,
							"mutable":              input.Object.SharedObject.Mutable,
						},
					}
				}
			} else if input.Pure != nil {
				serializedInput.Type = "pure"
				serializedInput.Value["Pure"] = bcs.ToBase64(input.Pure.Bytes)
			} else if input.UnresolvedPure != nil {
				serializedInput.Type = "pure"
				serializedInput.Value = input.UnresolvedPure.Value.(map[string]interface{})
			} else if input.UnresolvedObject != nil {
				serializedInput.Type = "object"
				serializedInput.Value["UnresolvedObject"] = map[string]interface{}{
					"objectId": ConvertMgoAddressBytesToString(input.UnresolvedObject.ObjectId),
				}
			}

			inputs = append(inputs, serializedInput)
		}
		result.Inputs = inputs

		// Process transactions (commands)
		transactions := []SerializedTransaction{}
		for _, cmd := range td.V1.Kind.ProgrammableTransaction.Commands {
			if cmd.MakeMoveVec != nil {
				tx := SerializedTransaction{
					Kind: "MakeMoveVec",
				}
				if cmd.MakeMoveVec.Type != nil {
					tx.Type = map[string]interface{}{
						"Some": *cmd.MakeMoveVec.Type,
					}
				} else {
					tx.Type = map[string]interface{}{
						"None": true,
					}
				}
				tx.Objects = convertArgumentsToSerializedFormat(cmd.MakeMoveVec.Elements, result.Inputs)
				transactions = append(transactions, tx)
			} else if cmd.MergeCoins != nil {
				tx := SerializedTransaction{
					Kind: "MergeCoins",
				}
				tx.Destination = convertArgumentToSerializedFormat(cmd.MergeCoins.Destination, result.Inputs)
				tx.Sources = convertArgumentsToSerializedFormat(cmd.MergeCoins.Sources, result.Inputs)
				transactions = append(transactions, tx)
			} else if cmd.MoveCall != nil {
				tx := SerializedTransaction{
					Kind: "MoveCall",
					Target: fmt.Sprintf("%s::%s::%s",
						ConvertMgoAddressBytesToString(cmd.MoveCall.Package),
						cmd.MoveCall.Module,
						cmd.MoveCall.Function),
				}
				typeArgs := []string{}
				for _, typeArg := range cmd.MoveCall.TypeArguments {
					// Convert TypeTag to string representation
					typeStr := convertTypeTagToString(typeArg)
					typeArgs = append(typeArgs, typeStr)
				}
				tx.TypeArguments = typeArgs
				tx.Arguments = convertArgumentsToSerializedFormat(cmd.MoveCall.Arguments, result.Inputs)
				transactions = append(transactions, tx)
			} else if cmd.Publish != nil {
				tx := SerializedTransaction{
					Kind: "Publish",
				}
				modules := [][]byte{}
				for _, mod := range cmd.Publish.Modules {
					modules = append(modules, mod[:])
				}
				tx.Modules = modules

				dependencies := []string{}
				for _, dep := range cmd.Publish.Dependencies {
					dependencies = append(dependencies, ConvertMgoAddressBytesToString(dep))
				}
				tx.Dependencies = dependencies
				transactions = append(transactions, tx)
			} else if cmd.SplitCoins != nil {
				tx := SerializedTransaction{
					Kind: "SplitCoins",
				}
				tx.Coin = convertArgumentToSerializedFormat(cmd.SplitCoins.Coin, result.Inputs)
				tx.Amounts = convertArgumentsToSerializedFormat(cmd.SplitCoins.Amount, result.Inputs)
				transactions = append(transactions, tx)
			} else if cmd.TransferObjects != nil {
				tx := SerializedTransaction{
					Kind: "TransferObjects",
				}
				tx.Objects = convertArgumentsToSerializedFormat(cmd.TransferObjects.Objects, result.Inputs)
				tx.Address = convertArgumentToSerializedFormat(cmd.TransferObjects.Address, result.Inputs)
				transactions = append(transactions, tx)
			} else if cmd.Upgrade != nil {
				tx := SerializedTransaction{
					Kind: "Upgrade",
				}
				modules := [][]byte{}
				for _, mod := range cmd.Upgrade.Modules {
					modules = append(modules, mod[:])
				}
				tx.Modules = modules

				dependencies := []string{}
				for _, dep := range cmd.Upgrade.Dependencies {
					dependencies = append(dependencies, ConvertMgoAddressBytesToString(dep))
				}
				tx.Dependencies = dependencies
				tx.PackageId = ConvertMgoAddressBytesToString(cmd.Upgrade.Package)
				tx.Ticket = convertArgumentToSerializedFormat(cmd.Upgrade.Ticket, result.Inputs)
				transactions = append(transactions, tx)
			}
		}
		result.Transactions = transactions
	}

	return result, nil
}

// Helper functions for serialization
func convertArgumentToSerializedFormat(arg *Argument, inputs []SerializedInput) interface{} {
	if arg == nil {
		return nil
	}

	if arg.Input != nil {
		return map[string]interface{}{
			"kind":  "Input",
			"index": *arg.Input,
		}
	} else if arg.Result != nil {
		return map[string]interface{}{
			"kind":  "Result",
			"index": *arg.Result,
		}
	} else if arg.NestedResult != nil {
		return map[string]interface{}{
			"kind":        "NestedResult",
			"index":       arg.NestedResult.Index,
			"resultIndex": arg.NestedResult.ResultIndex,
		}
	} else {
		return map[string]interface{}{
			"kind": "GasCoin",
		}
	}
}

func convertArgumentsToSerializedFormat(args []*Argument, inputs []SerializedInput) []interface{} {
	result := []interface{}{}
	for _, arg := range args {
		result = append(result, convertArgumentToSerializedFormat(arg, inputs))
	}
	return result
}

func convertTypeTagToString(typeTag *TypeTag) string {
	if typeTag == nil {
		return ""
	}

	if typeTag.Bool != nil {
		return "bool"
	} else if typeTag.U8 != nil {
		return "u8"
	} else if typeTag.U16 != nil {
		return "u16"
	} else if typeTag.U32 != nil {
		return "u32"
	} else if typeTag.U64 != nil {
		return "u64"
	} else if typeTag.U128 != nil {
		return "u128"
	} else if typeTag.U256 != nil {
		return "u256"
	} else if typeTag.Address != nil {
		return "address"
	} else if typeTag.Signer != nil {
		return "signer"
	} else if typeTag.Vector != nil {
		return fmt.Sprintf("vector<%s>", convertTypeTagToString(typeTag.Vector))
	} else if typeTag.Struct != nil {
		typeParams := ""
		if len(typeTag.Struct.TypeParams) > 0 {
			params := []string{}
			for _, param := range typeTag.Struct.TypeParams {
				params = append(params, convertTypeTagToString(param))
			}
			typeParams = fmt.Sprintf("<%s>", strings.Join(params, ", "))
		}
		return fmt.Sprintf("%s::%s::%s%s",
			ConvertMgoAddressBytesToString(typeTag.Struct.Address),
			typeTag.Struct.Module,
			typeTag.Struct.Name,
			typeParams)
	}

	return ""
}

// Serialize returns a JSON string representation of the transaction data in V1 format
func (td *TransactionData) Serialize() (string, error) {
	serialized, err := td.SerializeV1()
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(serialized)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// DeserializeFromJSON creates a TransactionData object from a JSON string
func DeserializeFromJSON(jsonStr string) (*TransactionData, error) {
	var serialized SerializedTransactionDataV1
	err := json.Unmarshal([]byte(jsonStr), &serialized)
	if err != nil {
		return nil, err
	}

	// Create a new TransactionData object
	td := &TransactionData{
		V1: &TransactionDataV1{
			Kind: &TransactionKind{
				ProgrammableTransaction: &ProgrammableTransaction{
					Inputs:   []*CallArg{},
					Commands: []*Command{},
				},
			},
			GasData:    &GasData{},
			Expiration: nil, // Will be set only if present in JSON
		},
	}

	// Set sender if available
	if serialized.Sender != "" {
		senderBytes, err := ConvertMgoAddressStringToBytes(model.MgoAddress(serialized.Sender))
		if err != nil {
			return nil, err
		}
		td.V1.Sender = senderBytes
	}

	// Set expiration
	if serialized.Expiration != nil {
		td.V1.Expiration = &TransactionExpiration{}
		if epoch, ok := serialized.Expiration["Epoch"]; ok {
			epochVal := uint64(epoch.(float64))
			td.V1.Expiration.Epoch = &epochVal
		} else {
			td.V1.Expiration.None = struct{}{}
		}
	}

	// Set gas data
	if serialized.GasConfig.Owner != "" {
		ownerBytes, err := ConvertMgoAddressStringToBytes(model.MgoAddress(serialized.GasConfig.Owner))
		if err != nil {
			return nil, err
		}
		td.V1.GasData.Owner = ownerBytes
	}

	if serialized.GasConfig.Budget != "" {
		budget, err := strconv.ParseUint(serialized.GasConfig.Budget, 10, 64)
		if err != nil {
			return nil, err
		}
		td.V1.GasData.Budget = &budget
	}

	if serialized.GasConfig.Price != "" {
		price, err := strconv.ParseUint(serialized.GasConfig.Price, 10, 64)
		if err != nil {
			return nil, err
		}
		td.V1.GasData.Price = &price
	}

	if len(serialized.GasConfig.Payment) > 0 {
		payment := []MgoObjectRef{}
		for _, p := range serialized.GasConfig.Payment {
			objBytes, err := ConvertMgoAddressStringToBytes(model.MgoAddress(p))
			if err != nil {
				return nil, err
			}
			// Note: This is simplified - in a real implementation you would need to get the version and digest
			payment = append(payment, MgoObjectRef{
				ObjectId: *objBytes,
				Version:  0,                         // You would need to get the actual version
				Digest:   model.ObjectDigestBytes{}, // You would need to get the actual digest
			})
		}
		td.V1.GasData.Payment = &payment
	}

	// Process inputs
	for _, input := range serialized.Inputs {
		callArg := &CallArg{}

		if input.Type == "object" {
			if objVal, ok := input.Value["Object"].(map[string]interface{}); ok {
				objectArg := &ObjectArg{}
				if immOrOwned, ok := objVal["ImmOrOwned"].(map[string]interface{}); ok {
					objectId, err := ConvertMgoAddressStringToBytes(model.MgoAddress(immOrOwned["objectId"].(string)))
					if err != nil {
						return nil, err
					}

					version := uint64(immOrOwned["version"].(float64))

					// Note: This is simplified - in a real implementation you would need to parse the digest properly
					digest := model.ObjectDigestBytes{}

					objectArg.ImmOrOwnedObject = &MgoObjectRef{
						ObjectId: *objectId,
						Version:  version,
						Digest:   digest,
					}
				} else if receiving, ok := objVal["Receiving"].(map[string]interface{}); ok {
					objectId, err := ConvertMgoAddressStringToBytes(model.MgoAddress(receiving["objectId"].(string)))
					if err != nil {
						return nil, err
					}

					version := uint64(receiving["version"].(float64))

					// Note: This is simplified - in a real implementation you would need to parse the digest properly
					digest := model.ObjectDigestBytes{}

					objectArg.Receiving = &MgoObjectRef{
						ObjectId: *objectId,
						Version:  version,
						Digest:   digest,
					}
				} else if shared, ok := objVal["Shared"].(map[string]interface{}); ok {
					objectId, err := ConvertMgoAddressStringToBytes(model.MgoAddress(shared["objectId"].(string)))
					if err != nil {
						return nil, err
					}

					initialSharedVersion := uint64(shared["initialSharedVersion"].(float64))
					mutable := shared["mutable"].(bool)

					objectArg.SharedObject = &SharedObjectRef{
						ObjectId:             *objectId,
						InitialSharedVersion: initialSharedVersion,
						Mutable:              mutable,
					}
				}
				callArg.Object = objectArg
			} else if unresolvedObj, ok := input.Value["UnresolvedObject"].(map[string]interface{}); ok {
				objectId, err := ConvertMgoAddressStringToBytes(model.MgoAddress(unresolvedObj["objectId"].(string)))
				if err != nil {
					return nil, err
				}
				callArg.UnresolvedObject = &UnresolvedObject{
					ObjectId: *objectId,
				}
			}
		} else if input.Type == "pure" {
			if pureVal, ok := input.Value["Pure"].(string); ok {
				// Handle base64-encoded string
				bytes, err := bcs.FromBase64(pureVal)
				if err != nil {
					return nil, err
				}
				callArg.Pure = &Pure{
					Bytes: bytes,
				}
			} else if pureVal, ok := input.Value["Pure"].([]interface{}); ok {
				// Handle array of numbers (legacy format)
				bytes := make([]byte, len(pureVal))
				for i, v := range pureVal {
					bytes[i] = byte(v.(float64))
				}
				callArg.Pure = &Pure{
					Bytes: bytes,
				}
			} else {
				callArg.UnresolvedPure = &UnresolvedPure{
					Value: input.Value,
				}
			}
		}

		td.V1.Kind.ProgrammableTransaction.Inputs = append(td.V1.Kind.ProgrammableTransaction.Inputs, callArg)
	}

	// Process transactions (commands)
	for _, tx := range serialized.Transactions {
		cmd := &Command{}

		switch tx.Kind {
		case "MakeMoveVec":
			elements := []*Argument{}
			for _, obj := range tx.Objects {
				arg := deserializeArgument(obj)
				elements = append(elements, arg)
			}

			makeMoveVec := &MakeMoveVec{
				Elements: elements,
			}

			if typeVal, ok := tx.Type.(map[string]interface{}); ok {
				if _, ok := typeVal["Some"]; ok {
					typeStr := typeVal["Some"].(string)
					makeMoveVec.Type = &typeStr
				}
			}

			cmd.MakeMoveVec = makeMoveVec

		case "MergeCoins":
			destination := deserializeArgument(tx.Destination)

			sources := []*Argument{}
			for _, src := range tx.Sources {
				arg := deserializeArgument(src)
				sources = append(sources, arg)
			}

			cmd.MergeCoins = &MergeCoins{
				Destination: destination,
				Sources:     sources,
			}

		case "MoveCall":
			// Parse target (package::module::function)
			parts := strings.Split(tx.Target, "::")
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid target format: %s", tx.Target)
			}

			packageBytes, err := ConvertMgoAddressStringToBytes(model.MgoAddress(parts[0]))
			if err != nil {
				return nil, err
			}

			arguments := []*Argument{}
			for _, arg := range tx.Arguments {
				arguments = append(arguments, deserializeArgument(arg))
			}

			typeArguments := []*TypeTag{}
			for _, typeArg := range tx.TypeArguments {
				tag, err := parseTypeTag(typeArg)
				if err != nil {
					return nil, err
				}
				typeArguments = append(typeArguments, tag)
			}

			cmd.MoveCall = &ProgrammableMoveCall{
				Package:       *packageBytes,
				Module:        parts[1],
				Function:      parts[2],
				Arguments:     arguments,
				TypeArguments: typeArguments,
			}

		case "Publish":
			modules := []model.MgoAddressBytes{}
			for _, mod := range tx.Modules {
				moduleBytes, err := ConvertBytesToMgoAddressBytes(mod)
				if err != nil {
					return nil, err
				}
				modules = append(modules, *moduleBytes)
			}

			dependencies := []model.MgoAddressBytes{}
			for _, dep := range tx.Dependencies {
				depBytes, err := ConvertMgoAddressStringToBytes(model.MgoAddress(dep))
				if err != nil {
					return nil, err
				}
				dependencies = append(dependencies, *depBytes)
			}

			cmd.Publish = &Publish{
				Modules:      modules,
				Dependencies: dependencies,
			}

		case "SplitCoins":
			coin := deserializeArgument(tx.Coin)

			amounts := []*Argument{}
			for _, amt := range tx.Amounts {
				amounts = append(amounts, deserializeArgument(amt))
			}

			cmd.SplitCoins = &SplitCoins{
				Coin:   coin,
				Amount: amounts,
			}

		case "TransferObjects":
			objects := []*Argument{}
			for _, obj := range tx.Objects {
				objects = append(objects, deserializeArgument(obj))
			}

			address := deserializeArgument(tx.Address)

			cmd.TransferObjects = &TransferObjects{
				Objects: objects,
				Address: address,
			}

		case "Upgrade":
			modules := []model.MgoAddressBytes{}
			for _, mod := range tx.Modules {
				moduleBytes, err := ConvertBytesToMgoAddressBytes(mod)
				if err != nil {
					return nil, err
				}
				modules = append(modules, *moduleBytes)
			}

			dependencies := []model.MgoAddressBytes{}
			for _, dep := range tx.Dependencies {
				depBytes, err := ConvertMgoAddressStringToBytes(model.MgoAddress(dep))
				if err != nil {
					return nil, err
				}
				dependencies = append(dependencies, *depBytes)
			}

			packageBytes, err := ConvertMgoAddressStringToBytes(model.MgoAddress(tx.PackageId))
			if err != nil {
				return nil, err
			}

			ticket := deserializeArgument(tx.Ticket)

			cmd.Upgrade = &Upgrade{
				Modules:      modules,
				Dependencies: dependencies,
				Package:      *packageBytes,
				Ticket:       ticket,
			}
		}

		td.V1.Kind.ProgrammableTransaction.Commands = append(td.V1.Kind.ProgrammableTransaction.Commands, cmd)
	}

	return td, nil
}

// Helper function to deserialize an argument
func deserializeArgument(arg interface{}) *Argument {
	if arg == nil {
		return nil
	}

	argMap, ok := arg.(map[string]interface{})
	if !ok {
		return nil
	}

	kind, ok := argMap["kind"].(string)
	if !ok {
		return nil
	}

	result := &Argument{}

	switch kind {
	case "Input":
		index := uint16(argMap["index"].(float64))
		result.Input = &index
	case "Result":
		index := uint16(argMap["index"].(float64))
		result.Result = &index
	case "NestedResult":
		index := uint16(argMap["index"].(float64))
		resultIndex := uint16(argMap["resultIndex"].(float64))
		result.NestedResult = &NestedResult{
			Index:       index,
			ResultIndex: resultIndex,
		}
	case "GasCoin":
		result.GasCoin = struct{}{}
	}

	return result
}

// Helper function to parse a type tag string
func parseTypeTag(typeStr string) (*TypeTag, error) {
	tag := &TypeTag{}

	switch typeStr {
	case "bool":
		tag.Bool = new(bool)
	case "u8":
		tag.U8 = new(bool)
	case "u16":
		tag.U16 = new(bool)
	case "u32":
		tag.U32 = new(bool)
	case "u64":
		tag.U64 = new(bool)
	case "u128":
		tag.U128 = new(bool)
	case "u256":
		tag.U256 = new(bool)
	case "address":
		tag.Address = new(bool)
	case "signer":
		tag.Signer = new(bool)
	default:
		// Handle vector and struct types
		if strings.HasPrefix(typeStr, "vector<") && strings.HasSuffix(typeStr, ">") {
			innerType := typeStr[7 : len(typeStr)-1]
			innerTag, err := parseTypeTag(innerType)
			if err != nil {
				return nil, err
			}
			tag.Vector = innerTag
		} else if strings.Contains(typeStr, "::") {
			// Parse struct type (e.g., "0x2::mgo::MGO")
			// We need to be careful about "::" inside angle brackets
			parts := parseStructTypeParts(typeStr)
			if len(parts) < 3 {
				return nil, fmt.Errorf("invalid struct type format: %s", typeStr)
			}

			addressBytes, err := ConvertMgoAddressStringToBytes(model.MgoAddress(parts[0]))
			if err != nil {
				return nil, err
			}

			// Check if there are type parameters
			name := parts[2]
			typeParams := []*TypeTag{}

			if strings.Contains(name, "<") && strings.HasSuffix(name, ">") {
				// Extract name and type parameters
				nameEnd := strings.Index(name, "<")
				paramsStr := name[nameEnd+1 : len(name)-1]
				name = name[:nameEnd]

				// Parse type parameters with proper handling of nested angle brackets
				paramStrs := parseTypeParameters(paramsStr)
				for _, paramStr := range paramStrs {
					paramTag, err := parseTypeTag(paramStr)
					if err != nil {
						return nil, err
					}
					typeParams = append(typeParams, paramTag)
				}
			}

			tag.Struct = &StructTag{
				Address:    *addressBytes,
				Module:     parts[1],
				Name:       name,
				TypeParams: typeParams,
			}
		} else {
			return nil, fmt.Errorf("unknown type: %s", typeStr)
		}
	}

	return tag, nil
}

// Helper function to parse type parameters with proper handling of nested angle brackets
func parseTypeParameters(paramsStr string) []string {
	if paramsStr == "" {
		return []string{}
	}

	var params []string
	var current strings.Builder
	depth := 0

	for i, char := range paramsStr {
		switch char {
		case '<':
			depth++
			current.WriteRune(char)
		case '>':
			depth--
			current.WriteRune(char)
		case ',':
			if depth == 0 {
				// We're at the top level, this comma separates parameters
				param := strings.TrimSpace(current.String())
				if param != "" {
					params = append(params, param)
				}
				current.Reset()
			} else {
				current.WriteRune(char)
			}
		case ' ':
			// Only add space if we're not at the beginning of a parameter or after a comma
			if current.Len() > 0 && depth > 0 {
				current.WriteRune(char)
			} else if current.Len() > 0 && depth == 0 {
				// Skip spaces at top level unless we're in the middle of a parameter
				if i+1 < len(paramsStr) && paramsStr[i+1] != ',' {
					current.WriteRune(char)
				}
			}
		default:
			current.WriteRune(char)
		}
	}

	// Add the last parameter
	param := strings.TrimSpace(current.String())
	if param != "" {
		params = append(params, param)
	}

	return params
}

// Helper function to parse struct type parts while respecting angle brackets
func parseStructTypeParts(typeStr string) []string {
	var parts []string
	var current strings.Builder
	depth := 0

	for i := 0; i < len(typeStr); i++ {
		char := typeStr[i]
		switch char {
		case '<':
			depth++
			current.WriteByte(char)
		case '>':
			depth--
			current.WriteByte(char)
		case ':':
			if depth == 0 && i+1 < len(typeStr) && typeStr[i+1] == ':' {
				// We found "::" at the top level
				part := current.String()
				if part != "" {
					parts = append(parts, part)
				}
				current.Reset()
				// Skip the next ':'
				i++
			} else {
				current.WriteByte(char)
			}
		default:
			current.WriteByte(char)
		}
	}

	// Add the last part
	part := current.String()
	if part != "" {
		parts = append(parts, part)
	}

	return parts
}
