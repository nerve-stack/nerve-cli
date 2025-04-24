package golang

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/nerve-stack/nerve-cli/internal/schema"
	"github.com/nerve-stack/nerve-cli/pkg/cases"
)

type parser struct {
	spec  *schema.Spec
	model *Model
}

func newParser(spec *schema.Spec) *parser {
	return &parser{
		spec: spec,
		model: &Model{
			Imports: map[string]struct{}{},
			Structs: map[string]GoStruct{},
			Errors:  map[string]GoError{},
			Methods: []GoMethod{},
			Events:  []GoEvent{},
			Enums:   map[string]GoEnum{},
		},
	}
}

func ParseSpec(spec *schema.Spec) (*Model, error) {
	p := newParser(spec)
	if err := p.parseSpec(); err != nil {
		return nil, err
	}

	return p.model, nil
}

func (r *parser) parseSpec() error {
	for _, methodName := range slices.Sorted(maps.Keys(r.spec.Methods)) {
		method := r.spec.Methods[methodName]
		gm := GoMethod{
			Name:            methodName,
			CapitalizedName: cases.ToCamelCase(methodName),
		}

		if method.Params != nil {
			var (
				params     schema.SchemaOrRef
				paramsName string
				err        error
			)

			if method.Params.Ref != nil {
				params, paramsName, err = schema.ResolveRef(r.spec, *method.Params.Ref)
				if err != nil {
					return err
				}
			} else {
				params = *method.Params
				paramsName = gm.CapitalizedName + "Params"
			}

			goType, err := r.parseSchema(paramsName, params)
			if err != nil {
				return err
			}

			gm.ParamsType = &goType
		}

		if method.Result != nil {
			var (
				result     schema.SchemaOrRef
				resultName string
				err        error
			)

			if method.Result.Ref != nil {
				result, resultName, err = schema.ResolveRef(r.spec, *method.Result.Ref)
				if err != nil {
					return err
				}
			} else {
				result = *method.Result
				resultName = gm.CapitalizedName + "Result"
			}

			goType, err := r.parseSchema(resultName, result)
			if err != nil {
				return err
			}

			gm.ResultType = &goType
		}

		if len(method.Errors) != 0 {
			for _, errorName := range method.Errors {
				gm.Errors = append(gm.Errors, errorName)

				if _, ok := r.model.Errors[errorName]; ok {
					continue
				}

				errorSchema, ok := r.spec.Errors[errorName]
				if !ok {
					return fmt.Errorf("no %s error provided", errorName)
				}

				var dataType *string

				if errorSchema.Data != nil {
					dt, err := r.parseSchema(cases.ToCamelCase(errorName)+"Data", *errorSchema.Data)
					if err != nil {
						return err
					}

					dataType = &dt
				}

				r.model.Errors[errorName] = GoError{
					Name:     cases.ToCamelCase(errorName),
					Code:     errorSchema.Code,
					DataType: dataType,
				}
			}
		}

		r.model.Methods = append(r.model.Methods, gm)
	}

	for _, eventName := range slices.Sorted(maps.Keys(r.spec.Events)) {
		event := r.spec.Events[eventName]
		ge := GoEvent{
			Name:            eventName,
			CapitalizedName: cases.ToCamelCase(eventName),
		}
		var (
			params     schema.SchemaOrRef
			paramsName string
			err        error
		)

		if event.Ref != nil {
			params, paramsName, err = schema.ResolveRef(r.spec, *event.Ref)
			if err != nil {
				return err
			}
		} else {
			params = event
			paramsName = ge.CapitalizedName + "Params"
		}

		goType, err := r.parseSchema(paramsName, params)
		if err != nil {
			return err
		}

		ge.DataType = &goType
		r.model.Events = append(r.model.Events, ge)
	}

	return nil
}

func (r *parser) parseSchema(name string, sch schema.SchemaOrRef) (string, error) {
	if sch.Ref != nil {
		refSchema, refName, err := schema.ResolveRef(r.spec, *sch.Ref)
		if err != nil {
			return "", err
		}

		return r.parseSchema(refName, refSchema)
	}

	if sch.Type == nil {
		return "", errors.New("schema type is nil")
	}

	switch *sch.Type {
	case "object":
		gs := GoStruct{
			Fields: []GoField{},
		}

		for _, fieldName := range slices.Sorted(maps.Keys(sch.Properties)) {
			fieldSchema := sch.Properties[fieldName]
			fieldGoName := cases.ToCamelCase(fieldName)
			fieldTypeName := name + fieldGoName // Optional: unique struct names

			goType, err := r.parseSchema(fieldTypeName, fieldSchema)
			if err != nil {
				return "", err
			}

			gs.Fields = append(gs.Fields, GoField{
				Name:     fieldGoName,
				Type:     goType,
				Tag:      fieldName,
				Optional: !slices.Contains(sch.Required, fieldName),
			})
		}

		r.model.Structs[name] = gs

		return name, nil

	case "array":
		if sch.Items == nil {
			return "", errors.New("array schema must have items")
		}

		itemType, err := r.parseSchema(name+"Item", *sch.Items)
		if err != nil {
			return "", err
		}

		return "[]" + itemType, nil

	default:
		return r.mapJsonSchemaTypeToGo(*sch.Type, sch.Format, sch.Enum, name)
	}
}
