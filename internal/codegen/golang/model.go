package golang

import (
	"fmt"
	"strconv"
)

type Model struct {
	NerveVersion string
	Package      string
	Imports      map[string]struct{}
	Structs      map[string]GoStruct
	Errors       map[string]GoError
	Methods      []GoMethod
	Enums        map[string]GoEnum
	Events       []GoEvent
}

type GoField struct {
	Name     string
	Type     string
	Tag      string
	Optional bool
}

type GoStruct struct {
	Fields []GoField
}

type GoError struct {
	Name     string
	Code     int
	DataType *string
}

type GoEvent struct {
	Name            string
	CapitalizedName string
	DataType        *string
}

type GoEnum struct {
	Name   string
	Type   string
	Values []string
}

type GoMethod struct {
	Name            string
	CapitalizedName string
	ParamsType      *string
	ResultType      *string
	Errors          []string
}

const (
	TimeImport = "time"
	UUIDImport = "github.com/google/uuid"
)

func derefOrNil(s *string) string {
	if s == nil {
		return "nil"
	}

	return *s
}

func (r *parser) mapJsonSchemaTypeToGo(typ string, format *string, enumValues []string, enumName string) (string, error) {
	if len(enumValues) > 0 {
		var goType string
		values := make([]string, len(enumValues))

		switch typ {
		case "string":
			if format != nil && *format == "uuid" {
				goType = "uuid.UUID"
			} else {
				goType = "string"
			}

			for i, v := range enumValues {
				values[i] = "\"" + v + "\""
			}

		case "integer":
			goType = "int"

			for _, v := range enumValues {
				_, err := strconv.Atoi(v)
				if err != nil {
					return "", fmt.Errorf("enum value %q is not a valid integer: %w", v, err)
				}

				values = append(values, v)
			}

		case "number":
			goType = "float64"

			for _, v := range enumValues {
				if _, err := strconv.ParseFloat(v, 64); err != nil {
					return "", fmt.Errorf("enum value %q is not a valid number: %w", v, err)
				}

				values = append(values, v)
			}

		case "boolean":
			goType = "bool"

			for _, v := range enumValues {
				if v != "true" && v != "false" {
					return "", fmt.Errorf("enum value %q is not a valid boolean", v)
				}

				values = append(values, v)
			}

		default:
			return "", fmt.Errorf("unsupported enum type: %q", typ)
		}

		r.model.Enums[enumName] = GoEnum{
			Name:   enumName,
			Type:   goType,
			Values: values,
		}

		return enumName, nil
	}

	switch typ {
	case "string":
		if format != nil {
			switch *format {
			case "uuid":
				r.model.Imports[UUIDImport] = struct{}{}
				return "uuid.UUID", nil

			case "date-time":
				r.model.Imports[TimeImport] = struct{}{}
				return "time.Time", nil
			}
		}

		return "string", nil

	case "integer":
		return "int", nil

	case "number":
		return "float64", nil

	case "boolean":
		return "bool", nil

	default:
		return "", fmt.Errorf("wrong type: typ: %s, format: %s", typ, derefOrNil(format))
	}
}
