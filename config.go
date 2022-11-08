package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Returns a `string` that can be converted to the "zero" value for `int`,
// `string`, and `bool`.
//
// Returns an `error` for unsupported types.
func defaultValueFor(kind reflect.Kind) (string, error) {
	switch kind {
	case reflect.Int:
		return "0", nil
	case reflect.String:
		return "", nil
	case reflect.Bool:
		return "false", nil
	default:
		return "", fmt.Errorf("unsupported kind %v", kind)
	}
}

// Converts `value` into an `int`, `string`, or `bool` (depending on the
// reflect.Kind of `vfield`) and assigns it to `vfield`.
//
// Returns an `error` for unsupported types.
func set(vfield reflect.Value, value string) error {
	switch vfield.Kind() {
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		vfield.SetInt(int64(intValue))
	case reflect.String:
		vfield.SetString(value)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		vfield.SetBool(boolValue)
	default:
		return fmt.Errorf("unsupported kind %v", vfield.Kind())
	}
	return nil
}

// Initializes `c` (which should be a pointer to a struct) according to its
// field's tags.
//
// See https://github.com/runtimeverification/go-config/README.md for more
// information.
func Init(c interface{}) error {
	// An instance of *c as a reflect.Value.
	elem := reflect.ValueOf(c).Elem()

	// A reflect.Type instance for *c.
	t := reflect.TypeOf(c).Elem()

	// Find the tags specified in the struct for c.
	//
	for i := 0; i < t.NumField(); i++ {
		// The reflect.StructField object.
		sfield := t.Field(i)

		// Its tags.
		tags := sfield.Tag

		// The reflect.Value corresponding to this field.
		vfield := elem.Field(i)

		// The "default" tag.
		defaultValueString, ok := tags.Lookup("default")

		// If we did not find a "default" tag,
		// set `defaultValueString` to a pre-determined value.
		if !ok {
			val, err := defaultValueFor(sfield.Type.Kind())
			if err != nil {
				return err
			}
			defaultValueString = val
		}

		// The "env" tag.
		envTag, ok := tags.Lookup("env")

		// If the tag exists and a value is found in the environment,
		// change `defaultValueString` to that value.
		if ok {
			envValue, ok := os.LookupEnv(envTag)
			if ok {
				defaultValueString = envValue
			}
		}

		// The "cli" tag.
		cliTag, ok := tags.Lookup("cli")

		// If no "cli" tag is found, set to the current `defaultValueString`
		// and move on.
		if !ok {
			err := set(vfield, defaultValueString)
			if err != nil {
				return err
			}
			continue
		}

		// The "desc" tag.
		description := tags.Get("desc")

		// We don't care for missing values above ^^
		// because it will simply be set to the empty string.
		// An empty description is fine.

		// Use the `flags` package.
		switch sfield.Type.Kind() {
		case reflect.Int:
			intValue, err := strconv.Atoi(defaultValueString)
			if err != nil {
				return err
			}
			flag.IntVar(
				(*int)(vfield.Addr().UnsafePointer()),
				cliTag,
				intValue,
				description,
			)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(defaultValueString)
			if err != nil {
				return err
			}
			flag.BoolVar(
				(*bool)(vfield.Addr().UnsafePointer()),
				cliTag,
				boolValue,
				description,
			)
		case reflect.String:
			flag.StringVar(
				(*string)(vfield.Addr().UnsafePointer()),
				cliTag,
				defaultValueString,
				description,
			)
		}
	}

	flag.Parse()

	return nil
}
