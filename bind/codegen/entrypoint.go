package codegen

import (
	"fmt"
	"sort"
	"strconv"

	"blockwatch.cc/tzgo/micheline"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// parseEntrypoints add entrypoints found in script, in the ParseContext.
func parseEntrypoints(c *ParseContext, script *micheline.Script) error {
	sEntrypoints, err := script.Entrypoints(true)
	if err != nil {
		return err
	}

	// Sort them by ID first, so the auto generated incremental names are sorted.
	lEntrypoints := make([]micheline.Entrypoint, len(sEntrypoints))
	i := 0
	for _, e := range sEntrypoints {
		lEntrypoints[i] = e
		i++
	}
	sort.SliceStable(lEntrypoints, func(i, j int) bool { return lEntrypoints[i].Id < lEntrypoints[j].Id })

	for _, e := range lEntrypoints {
		if err := parseEntrypoint(c, &e); err != nil {
			return errors.Wrapf(err, "failed to parse `%s` entrypoint", e.Call)
		}
	}

	return nil
}

func parseEntrypoint(c *ParseContext, entrypoint *micheline.Entrypoint) error {
	e := &Entrypoint{
		Name: c.fmt.Method(entrypoint.Call),
	}

	// Iterate over entrypoint arguments
	for i, arg := range entrypoint.Typedef {
		if err := parseEntrypointArg(c, e, &arg, i); err != nil {
			if errors.Is(err, errSkipEntrypoint) {
				log.Tracef("Skipping entrypoint `%s`", entrypoint.Call)
				return nil
			}
			return errors.Wrap(err, "failed to parse entrypoint arguments")
		}
	}

	c.entrypoints = append(c.entrypoints, e)
	return nil
}

// errSkipEntrypoint is returned by parseEntrypointArg, when the whole entrypoint
// should be skipped.
var errSkipEntrypoint = errors.New("skip entrypoint")

func parseEntrypointArg(c *ParseContext, e *Entrypoint, arg *micheline.Typedef, index int) error {
	a := &Arg{}

	// `unit` and `contract` breaks the flow control
	if arg.Type == micheline.T_UNIT.String() {
		// `unit` type is used for entrypoints without arguments.
		if index != 0 {
			panic("found a \"unit\" parameter not at first index")
		}
		return nil
	}
	if arg.Type == micheline.T_CONTRACT.String() {
		// `contract` type is used in getters, that are not usable outside
		// of a smart contract.
		return errSkipEntrypoint
	}

	name := arg.Name
	if name == "" || isInt(name) {
		name = fmt.Sprintf("arg%d", index)
	}
	a.Name = c.fmt.Argument(name)

	var err error
	a.Type, err = parseType(c, arg, e.Name.Normalized)
	if err != nil {
		return err
	}

	e.Args = append(e.Args, a)
	return nil
}

func parseType(c *ParseContext, arg *micheline.Typedef, refName string) (Type, error) {
	if arg.Optional {
		return parseOptionType(c, arg, refName)
	}

	// From tzgo/micheline doc:
	// a.Type is "struct", "union", or an OpCode
	if arg.Type == "struct" {
		return parseStructType(c, arg, refName)
	} else if arg.Type == "union" {
		return parseUnionType(c, arg, refName)
	}

	// Here, we can assume that the type is an OpCode.
	// Else, that means that the micheline is malformed.
	op, err := micheline.ParseOpCode(arg.Type)
	if err != nil {
		return nil, errors.Wrap(err, "invalid opcode")
	}
	if !op.IsTypeCode() {
		return nil, errors.Errorf("not a type opcode: %s", arg.Type)
	}

	if IsScalarType(op) {
		return TypeFromScalar(op), nil
	}

	if IsContainerType(op) {
		return parseContainerType(c, arg, op, refName)
	}

	return nil, errors.Errorf("opCode not handled: %s", op)
}

func parseStructType(c *ParseContext, arg *micheline.Typedef, refName string) (*StructType, error) {
	s := c.registerAutoStruct(refName, "Arg")
	for _, field := range arg.Args {
		fieldType, err := parseType(c, &field, refName)
		if err != nil {
			return nil, err
		}
		fieldName := field.Name
		if fieldName == "" || isInt(fieldName) {
			fieldName = c.autoFieldName(refName)
		}
		s.Fields = append(s.Fields, &Field{
			Name: c.fmt.Field(fieldName),
			Type: fieldType,
		})
	}

	return s, nil
}

func parseUnionType(c *ParseContext, arg *micheline.Typedef, refName string) (*UnionType, error) {
	if len(arg.Args) != 2 {
		return nil, errors.Errorf("got union with %d args", len(arg.Args))
	}

	u := c.registerAutoUnion(refName)

	var err error

	u.LType, err = parseType(c, &arg.Args[0], refName)
	if err != nil {
		return nil, err
	}

	u.RType, err = parseType(c, &arg.Args[1], refName)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func parseOptionType(c *ParseContext, arg *micheline.Typedef, refName string) (*OptionType, error) {
	o := c.registerAutoOption(refName)

	// Process the underlying type, without `Optional`
	withoutOption := *arg
	withoutOption.Optional = false

	var err error
	o.Type, err = parseType(c, &withoutOption, refName)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func parseContainerType(c *ParseContext, arg *micheline.Typedef, op micheline.OpCode, refName string) (Type, error) {
	switch op {
	case micheline.T_LIST,
		micheline.T_SET:
		itemType, err := parseType(c, &arg.Args[0], refName)
		if err != nil {
			return nil, err
		}
		return &ListType{T: itemType}, nil
	}

	return nil, errors.Errorf("opCode not handled: %s", op)
}

// region micheline utility

func IsScalarType(op micheline.OpCode) bool {
	switch op {
	case micheline.T_BOOL,
		micheline.T_CONTRACT, // skipped
		micheline.T_INT,
		micheline.T_KEY,
		micheline.T_KEY_HASH,
		micheline.T_NAT,
		micheline.T_SIGNATURE,
		micheline.T_STRING,
		micheline.T_BYTES,
		micheline.T_MUTEZ,
		micheline.T_TIMESTAMP,
		micheline.T_UNIT, // skipped
		micheline.T_OPERATION,
		micheline.T_ADDRESS,
		micheline.T_CHAIN_ID,
		micheline.T_NEVER,
		micheline.T_BLS12_381_G1,
		micheline.T_BLS12_381_G2,
		micheline.T_BLS12_381_FR,
		micheline.T_SAPLING_STATE,
		micheline.T_SAPLING_TRANSACTION:
		return true
	default:
		return false
	}
}

func IsContainerType(op micheline.OpCode) bool {
	switch op {
	case micheline.T_MAP,
		micheline.T_LIST,
		micheline.T_SET,
		micheline.T_LAMBDA:
		return true
	default:
		return false
	}
}

func TypeFromScalar(op micheline.OpCode) Type {
	// We use strings for addresses for readability.
	if op == micheline.T_ADDRESS {
		return &StringType{}
	}

	switch op.PrimType() {
	case micheline.PrimInt:
		return &IntType{}
	case micheline.PrimString:
		return &StringType{}
	case micheline.PrimBytes:
		return &BytesType{}
	default:
		panic(fmt.Sprintf("invalid opCode: %s", op))
	}
}

// endregion

func isInt(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
