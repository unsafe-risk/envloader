package envloader

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// dataProvider defines a function that provides the value for a given key.
type dataProvider func(key string) (string, error)

// BindStruct binds data from a provider to a struct's fields using reflection.
func BindStruct(config interface{}, provider dataProvider) error {
	// Get the value of the config.
	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to a struct")
	}

	// Iterate over the struct fields.
	for i := 0; i < v.Elem().NumField(); i++ {
		field := v.Elem().Field(i)
		fieldType := v.Elem().Type().Field(i)

		// Skip unexported fields.
		if !fieldType.IsExported() {
			continue
		}

		// Get the "env" tag value.
		tag := fieldType.Tag.Get("env")
		if tag == "" {
			continue
		}

		// Split the tag into the environment variable name and options.
		parts := strings.Split(tag, ",")
		envVarName := parts[0]

		// Check for the "required" option in the tag.
		required := false
		for _, opt := range parts[1:] {
			if opt == "required" {
				required = true
				break
			}
		}

		// Get the value from the data provider using the envVarName.
		envValue, err := provider(envVarName)
		if err != nil {
			if required {
				return fmt.Errorf("failed to get value for field %s: %w", fieldType.Name, err)
			} else {
				continue
			}
		}

		if envValue == "" && required {
			return fmt.Errorf("required environment variable %s is missing", envVarName)
		}

		// Set the field value.
		switch field.Kind() {
		case reflect.String:
			field.SetString(envValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if val, err := strconv.ParseInt(envValue, 10, 64); err != nil {
				return fmt.Errorf("failed to parse int value for field %s: %w", fieldType.Name, err)
			} else {
				field.SetInt(val)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if val, err := strconv.ParseUint(envValue, 10, 64); err != nil {
				return fmt.Errorf("failed to parse uint value for field %s: %w", fieldType.Name, err)
			} else {
				field.SetUint(val)
			}
		case reflect.Float32, reflect.Float64:
			if val, err := strconv.ParseFloat(envValue, 64); err != nil {
				return fmt.Errorf("failed to parse float value for field %s: %w", fieldType.Name, err)
			} else {
				field.SetFloat(val)
			}
		case reflect.Complex64, reflect.Complex128:
			if val, err := strconv.ParseComplex(envValue, 64); err != nil {
				return fmt.Errorf("failed to parse complex value for field %s: %w", fieldType.Name, err)
			} else {
				field.SetComplex(val)
			}
		case reflect.Bool:
			switch envValue {
			case "Y", "y", "Yes", "YES", "yes", "on":
				field.SetBool(true)
			case "N", "n", "No", "NO", "no", "off":
				field.SetBool(false)
			default:
				if val, err := strconv.ParseBool(envValue); err != nil {
					return fmt.Errorf("failed to parse bool value for field %s: %w", fieldType.Name, err)
				} else {
					field.SetBool(val)
				}
			}
		default:
			return fmt.Errorf("unsupported field type: %s", field.Kind())
		}
	}

	return nil
}
