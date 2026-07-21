package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"

	"github.com/fstln-dev/utils-lib-go/pkg/graphql/ident"
)

var jsonUnmarshaler = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()

func constructQuery(op Operation, query interface{}, variables Variables) string {
	var buf bytes.Buffer
	writeQuery(&buf, reflect.TypeOf(query), reflect.ValueOf(query), false)

	if len(variables) > 0 {
		return string(op) + "(" + queryArguments(variables) + ")" + buf.String()
	}
	return buf.String()
}

func writeQuery(w io.Writer, t reflect.Type, v reflect.Value, inline bool) error {
	switch t.Kind() {
	case reflect.Ptr:
		err := writeQuery(w, t.Elem(), ElemSafe(v), false)
		if err != nil {
			return fmt.Errorf("failed to write query for ptr `%v`: %w", t, err)
		}
	case reflect.Struct:
		// If the type implements json.Unmarshaler, it's a scalar. Don't expand it.
		if reflect.PtrTo(t).Implements(jsonUnmarshaler) {
			return nil
		}
		if !inline {
			io.WriteString(w, "{")
		}
		iter := 0
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			value, ok := f.Tag.Lookup("graphql")
			// Skip this field if the tag value is hyphen
			if value == "-" {
				continue
			}
			if iter != 0 {
				io.WriteString(w, ",")
			}
			iter++

			inlineField := f.Anonymous && !ok
			if !inlineField {
				if ok {
					io.WriteString(w, value)
				} else {
					io.WriteString(w, ident.ParseMixedCaps(f.Name).ToLowerCamelCase())
				}
			}
			// Skip writeQuery if the GraphQL type associated with the filed is scalar
			if isTrue(f.Tag.Get("scalar")) {
				continue
			}
			err := writeQuery(w, f.Type, FieldSafe(v, i), inlineField)
			if err != nil {
				return fmt.Errorf("failed to write query for struct field `%v`: %w", f.Name, err)
			}
		}
		if !inline {
			io.WriteString(w, "}")
		}
	case reflect.Slice:
		if t.Elem().Kind() != reflect.Array {
			err := writeQuery(w, t.Elem(), IndexSafe(v, 0), false)
			if err != nil {
				return fmt.Errorf("failed to write query for slice item `%v`: %w", t, err)
			}
			return nil
		}
		// handle [][2]interface{} like an ordered map
		if t.Elem().Len() != 2 {
			return fmt.Errorf("only arrays of len 2 are supported, got %v", t.Elem())
		}
		sliceOfPairs := v
		_, _ = io.WriteString(w, "{")
		for i := 0; i < sliceOfPairs.Len(); i++ {
			pair := sliceOfPairs.Index(i)
			// it.Value() returns interface{}, so we need to use reflect.ValueOf
			// to cast it away
			key, val := pair.Index(0), reflect.ValueOf(pair.Index(1).Interface())
			keyString, ok := key.Interface().(string)
			if !ok {
				return fmt.Errorf("expected pair (string, %v), got (%v, %v)",
					val.Type(), key.Type(), val.Type())
			}
			_, _ = io.WriteString(w, keyString)
			err := writeQuery(w, val.Type(), val, false)
			if err != nil {
				return fmt.Errorf("failed to write query for pair[1] `%v`: %w", val.Type(), err)
			}
		}
		_, _ = io.WriteString(w, "}")
	case reflect.Map:
		return fmt.Errorf("type %v is not supported, use [][2]interface{} instead", t)
	}
	return nil
}

func IndexSafe(v reflect.Value, i int) reflect.Value {
	if v.IsValid() && i < v.Len() {
		return v.Index(i)
	}
	return reflect.ValueOf(nil)
}

func ElemSafe(v reflect.Value) reflect.Value {
	if v.IsValid() {
		return v.Elem()
	}
	return reflect.ValueOf(nil)
}

func FieldSafe(valStruct reflect.Value, i int) reflect.Value {
	if valStruct.IsValid() {
		return valStruct.Field(i)
	}
	return reflect.ValueOf(nil)
}

func isTrue(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func queryArguments(variables map[string]interface{}) string {
	// Sort keys in order to produce deterministic output for testing purposes.
	// TODO: If tests can be made to work with non-deterministic output, then no need to sort.
	keys := make([]string, 0, len(variables))
	for k := range variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		io.WriteString(&buf, "$")
		io.WriteString(&buf, k)
		io.WriteString(&buf, ":")
		writeArgumentType(&buf, reflect.TypeOf(variables[k]), true)
		// Don't insert a comma here.
		// Commas in GraphQL are insignificant, and we want minified output.
		// See https://facebook.github.io/graphql/October2016/#sec-Insignificant-Commas.
	}
	return buf.String()
}

// writeArgumentType writes a minified GraphQL type for t to w.
// value indicates whether t is a value (required) type or pointer (optional) type.
// If value is true, then "!" is written at the end of t.
func writeArgumentType(w io.Writer, t reflect.Type, value bool) {
	if t.Kind() == reflect.Ptr {
		// Pointer is an optional type, so no "!" at the end of the pointer's underlying type.
		writeArgumentType(w, t.Elem(), false)
		return
	}

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		// List. E.g., "[Int]".
		io.WriteString(w, "[")
		writeArgumentType(w, t.Elem(), true)
		io.WriteString(w, "]")
	default:
		// Named type. E.g., "Int".
		name := t.Name()
		if name == "string" { // HACK: Workaround for https://github.com/shurcooL/githubv4/issues/12.
			name = "ID"
		}
		io.WriteString(w, name)
	}

	if value {
		// Value is a required type, so add "!" to the end.
		io.WriteString(w, "!")
	}
}
