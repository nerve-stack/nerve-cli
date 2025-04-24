package schema

type Spec struct {
	Version string                         `json:"version" yaml:"version"`
	Info    Info                           `json:"info"    yaml:"info"`
	Schemas Schemas                        `json:"schemas" yaml:"schemas"`
	Methods map[string]Method              `json:"methods" yaml:"methods"`
	Events  map[string]SchemaOrRefWithDesc `json:"events"  yaml:"events"`
}

type Method struct {
	Description *string       `json:"description,omitempty" yaml:"description"`
	Params      *SchemaOrRef  `json:"params,omitempty"      yaml:"params,omitempty"`
	Result      *SchemaOrRef  `json:"result,omitempty"      yaml:"result,omitempty"`
	Errors      []SchemaOrRef `json:"errors,omitempty"      yaml:"errors,omitempty"`
}

type Info struct {
	Name    string  `json:"name"              yaml:"name"`
	Version *string `json:"version,omitempty" yaml:"version"`
}

type Schemas struct {
	Types    map[string]SchemaOrRefWithDesc `json:"types"    yaml:"types"`
	Entities map[string]SchemaOrRefWithDesc `json:"entities" yaml:"entities"`
	Errors   map[string]ErrorSchema         `json:"errors"   yaml:"errors"`
}

type ErrorSchema struct {
	Code        int          `json:"code"                  yaml:"code"`
	Description *string      `json:"description,omitempty" yaml:"description"`
	Data        *SchemaOrRef `json:"data,omitempty"        yaml:"data"`
}

type SchemaOrRef struct {
	Ref        *string                `json:"$ref,omitempty"       yaml:"$ref,omitempty"`
	Type       *string                `json:"type,omitempty"       yaml:"type,omitempty"`
	Format     *string                `json:"format,omitempty"     yaml:"format,omitempty"`
	Enum       []string               `json:"enum,omitempty"       yaml:"enum,omitempty"`
	Items      *SchemaOrRef           `json:"items,omitempty"      yaml:"items,omitempty"`
	Properties map[string]SchemaOrRef `json:"properties,omitempty" yaml:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"   yaml:"required,omitempty"`
}

type SchemaOrRefWithDesc struct {
	SchemaOrRef `yaml:",inline"`
	Description *string `json:"description,omitempty" yaml:"description"`
}
