package kaeru

// TODO: custom error declarations
// TODO: collect all errors and then return

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type ParseAny interface {
	ParseAny(v any) error
}

type ParseInt interface {
	ParseInt(i int) error
}

type ParseString interface {
	ParseString(s string) error
}

type ParseInt8 interface {
	ParseInt8(i int8) error
}

type ParseInt16 interface {
	ParseInt16(i int16) error
}

type ParseInt32 interface {
	ParseInt32(i int32) error
}

type ParseInt64 interface {
	ParseInt64(i int64) error
}

type ParseUint8 interface {
	ParseUint8(i uint8) error
}

type ParseUint16 interface {
	ParseUint16(i uint16) error
}

type ParseUint32 interface {
	ParseUint32(i uint32) error
}

type ParseUint64 interface {
	ParseUint64(i uint64) error
}

type ParseFloat32 interface {
	ParseFloat32(f float32) error
}

type ParseFloat64 interface {
	ParseFloat64(f float64) error
}

type ParseStringMap interface {
	ParseStringMap(m map[string]string) error
}

type ParseMap interface {
	ParseMap(m map[string]any) error
}

type ParseSlice interface {
	ParseSlice(s []any) error
}

type ParseStringSlice interface {
	ParseStringSlice(s []string) error
}

type SetDefault interface {
	SetDefault()
}

func Parse(input any, output any) error {
	outVal := reflect.ValueOf(output)
	// Check if output is a pointer and is addressable
	// Is this correct?
	if outVal.Kind() != reflect.Ptr {
		return errors.New("output must be a pointer")
	}

	// Get the reflect Value and Type of both input and output
	inVal := reflect.ValueOf(input)
	outVal = outVal.Elem()

	return parseValue(inVal, outVal)
}

func ParseJson(r io.Reader, output any) error {
	decoder := json.NewDecoder(r)
	var v any
	err := decoder.Decode(&v)

	if err != nil {
		return err
	}

	return Parse(v, output)
}

func ParseJsonBytes(data []byte, output any) error {
	var v any
	err := json.Unmarshal(data, &v)

	if err != nil {
		return err
	}

	return Parse(v, output)
}

func parseValue(inVal reflect.Value, outVal reflect.Value) error {
	switch inVal.Kind() {
	case
		reflect.Array,
		reflect.Chan,
		reflect.Func,
		reflect.Pointer,
		reflect.Struct,
		reflect.UnsafePointer:
		panic("inVal is not a valid parseable value")
	}

	if !outVal.CanSet() {
		panic("outVal is not settable")
	}

	required := true
	if inVal.Kind() == reflect.Interface {
		inVal = inVal.Elem()
	}

	if outVal.Kind() == reflect.Pointer {
		if outVal.IsNil() {
			outVal.Set(reflect.New(outVal.Type().Elem()))
		}
		outVal = outVal.Elem()
		required = false
	}
	
	// Handle nil input values using default or returning error if required
	if !inVal.IsValid() {
		if defaultable, ok := outVal.Addr().Interface().(SetDefault); ok {
			defaultable.SetDefault()
		} else if required {
			return errors.New("inVal is nil but must be set")
		}

		return nil
	}
	
	if parser, ok := outVal.Addr().Interface().(ParseAny); ok {
		return parser.ParseAny(inVal.Interface())
	}

	// If types are the same we can just set them and call it a day
	if inVal.Type() == outVal.Type() {
		outVal.Set(inVal)
		return nil
	}

	inValKind := inVal.Kind()
	outValKind := outVal.Kind()

	if isPrimitive(inValKind) {
		return parsePrimitive(inVal, outVal)
	} else if inValKind == reflect.Map {
		return parseMap(inVal, outVal)
	} else if inValKind == reflect.Slice {
		return parseSlice(inVal, outVal)
	} else {
		return fmt.Errorf("unsupported kinds, in: %s, out: %s", inValKind, outValKind)
	}
}

func isPrimitive(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

// inVal and outVal must be a valid primitive kind
func parsePrimitive(inVal reflect.Value, outVal reflect.Value) error {
	if !isPrimitive(inVal.Kind()) {
		panic("inVal must be a primitive")
	}

	switch inVal.Kind() {
	case reflect.String:
		if parser, ok := outVal.Addr().Interface().(ParseString); ok {
			return parser.ParseString(inVal.String())
		}
	case reflect.Bool:
		if outVal.Kind() == reflect.Bool {
			outVal.Set(inVal.Convert(outVal.Type()))
			return nil
		}
	case reflect.Int8:
		if parser, ok := outVal.Addr().Interface().(ParseInt8); ok {
			return parser.ParseInt8(int8(inVal.Int()))
		}
		fallthrough
	case reflect.Int16:
		if parser, ok := outVal.Addr().Interface().(ParseInt16); ok {
			return parser.ParseInt16(int16(inVal.Int()))
		}
		fallthrough
	case reflect.Int32:
		if parser, ok := outVal.Addr().Interface().(ParseInt32); ok {
			return parser.ParseInt32(int32(inVal.Int()))
		}
		fallthrough
	case reflect.Int64:
		if parser, ok := outVal.Addr().Interface().(ParseInt64); ok {
			return parser.ParseInt64(inVal.Int())
		}
		fallthrough
	case reflect.Int:
		if parser, ok := outVal.Addr().Interface().(ParseInt); ok {
			return parser.ParseInt(int(inVal.Int()))
		}
	case reflect.Uint8:
		if parser, ok := outVal.Addr().Interface().(ParseUint8); ok {
			return parser.ParseUint8(uint8(inVal.Uint()))
		}
		fallthrough
	case reflect.Uint16:
		if parser, ok := outVal.Addr().Interface().(ParseUint16); ok {
			return parser.ParseUint16(uint16(inVal.Uint()))
		}
		fallthrough
	case reflect.Uint32:
		if parser, ok := outVal.Addr().Interface().(ParseUint32); ok {
			return parser.ParseUint32(uint32(inVal.Uint()))
		}
		fallthrough
	case reflect.Uint64:
		if parser, ok := outVal.Addr().Interface().(ParseUint64); ok {
			return parser.ParseUint64(inVal.Uint())
		}
	case reflect.Float32:
		if parser, ok := outVal.Addr().Interface().(ParseFloat32); ok {
			return parser.ParseFloat32(float32(inVal.Float()))
		}
		fallthrough
	case reflect.Float64:
		if parser, ok := outVal.Addr().Interface().(ParseFloat64); ok {
			return parser.ParseFloat64(inVal.Float())
		}
	}

	if inVal.CanConvert(outVal.Type()) {
		outVal.Set(inVal.Convert(outVal.Type()))
		return nil
	}

	return fmt.Errorf("inVal %s is not parseable to outVal %s", inVal.Type(), outVal.Type())
}

func parseMapToMap(inVal reflect.Value, outVal reflect.Value) error {
	if inVal.Kind() != reflect.Map {
		panic("inVal must be a map")
	}

	if outVal.Kind() != reflect.Map {
		panic("outVal must be a map")
	}

	inMapKeys := inVal.MapKeys()
	outMap := reflect.MakeMapWithSize(outVal.Type(), inVal.Len())
	outMapKeyType := outMap.Type().Key()
	outMapValueType := outMap.Type().Elem()

	for i := 0; i < len(inMapKeys); i++ {
		inKey := inMapKeys[i]
		inValue := inVal.MapIndex(inKey)
		outKey := reflect.New(outMapKeyType).Elem()
		outValue := reflect.New(outMapValueType).Elem()

		if err := parseValue(inKey, outKey); err != nil {
			return fmt.Errorf("error parsing map key %s: %w", inKey, err)
		}

		if err := parseValue(inValue, outValue); err != nil {
			return fmt.Errorf("error parsing map value %s: %w", inValue, err)
		}

		outMap.SetMapIndex(outKey, outValue)
	}

	outVal.Set(outMap)

	return nil
}

func parseMapToStruct(inVal reflect.Value, outVal reflect.Value) error {
	if inVal.Kind() != reflect.Map {
		panic("inVal must be a map")
	}

	if outVal.Kind() != reflect.Struct {
		panic("outVal must be a struct")
	}

	outType := outVal.Type()
	for i := 0; i < outVal.NumField(); i++ {
		field := outVal.Field(i)
		fieldType := outType.Field(i)
		fieldName := fieldType.Name

		tag := fieldType.Tag.Get("parse")

		if tag != "" {
			fieldName = tag
		}

		// Check if the field is exported
		if !field.CanSet() {
			continue
		}

		// Look for the field in the input map
		mapValue := inVal.MapIndex(reflect.ValueOf(fieldName))

		// Recur for nested structs or primitives
		if err := parseValue(mapValue, field); err != nil {
			return fmt.Errorf("error parsing field %s: %w", fieldName, err)
		}
	}

	return nil
}

func parseMap(inVal reflect.Value, outVal reflect.Value) error {
	if inVal.Kind() != reflect.Map {
		panic("inVal must be a map")
	}

	if m, ok := inVal.Interface().(map[string]string); ok {
		if parser, ok := outVal.Addr().Interface().(ParseStringMap); ok {
			return parser.ParseStringMap(m)
		}
	}

	if m, ok := inVal.Interface().(map[string]any); ok {
		if parser, ok := outVal.Addr().Interface().(ParseMap); ok {
			return parser.ParseMap(m)
		}
	}

	if outVal.Kind() == reflect.Struct {
		return parseMapToStruct(inVal, outVal)
	}

	if outVal.Kind() == reflect.Map {
		return parseMapToMap(inVal, outVal)
	}

	return fmt.Errorf("inVal %s is not parseable to outVal %s", inVal.Type(), outVal.Type())
}

func parseSliceToSlice(inVal reflect.Value, outVal reflect.Value) error {
	if inVal.Kind() != reflect.Slice {
		panic("inVal must be slice")
	}

	if outVal.Kind() != reflect.Slice {
		panic("outVal must be slice")
	}

	outSlice := reflect.MakeSlice(outVal.Type(), inVal.Len(), inVal.Cap())
	for i := 0; i < inVal.Len(); i++ {
		elem := outSlice.Index(i)
		if err := parseValue(inVal.Index(i), elem); err != nil {
			return err
		}
	}

	outVal.Set(outSlice)
	return nil
}

func parseSliceToArray(inVal reflect.Value, outVal reflect.Value) error {
	if inVal.Kind() != reflect.Slice {
		panic("inVal must be slice")
	}

	if outVal.Kind() != reflect.Array {
		panic("outVal must be array")
	}

	inLen := inVal.Len()
	outLen := outVal.Len()

	// Check if the input slice is longer than the output array
	if inLen > outLen {
		return fmt.Errorf("input slice (length %d) is longer than output array (length %d)", inLen, outLen)
	}

	// Copy elements from the input slice to the output array
	for i := 0; i < outLen; i++ {
		var inValIndexValue reflect.Value
		if i < inLen {
			inValIndexValue = inVal.Index(i)
		} else {
			inValIndexValue = reflect.ValueOf(nil)
		}

		if err := parseValue(inValIndexValue, outVal.Index(i)); err != nil {
			return fmt.Errorf("error parsing element at index %d: %w", i, err)
		}
	}

	return nil
}

// Parse slice input to slice output
func parseSlice(inVal reflect.Value, outVal reflect.Value) error {
	if inVal.Kind() != reflect.Slice {
		panic("inVal must be slice")
	}

	if s, ok := inVal.Interface().([]string); ok {
		if parser, ok := outVal.Addr().Interface().(ParseStringSlice); ok {
			return parser.ParseStringSlice(s)
		}
	}

	if parser, ok := outVal.Addr().Interface().(ParseSlice); ok {
		return parser.ParseSlice(inVal.Interface().([]any))
	}

	if outVal.Kind() == reflect.Slice {
		return parseSliceToSlice(inVal, outVal)
	}

	if outVal.Kind() == reflect.Array {
		return parseSliceToArray(inVal, outVal)
	}

	return fmt.Errorf("inVal %s is not parseable to outVal %s", inVal.Type(), outVal.Type())
}
