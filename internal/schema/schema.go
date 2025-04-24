package schema

type Spec struct {
	Version string                 `json:"version" yaml:"version"`
	Info    Info                   `json:"info"    yaml:"info"`
	Schemas map[string]SchemaOrRef `json:"schemas" yaml:"schemas"`
	Errors  map[string]ErrorSchema `json:"errors"  yaml:"errors"`
	Methods map[string]Method      `json:"methods" yaml:"methods"`
	Events  map[string]SchemaOrRef `json:"events"  yaml:"events"`
}

type Info struct {
	Name    string  `json:"name"    yaml:"name"`
	Version *string `json:"version" yaml:"version"`
}

type Method struct {
	Params *SchemaOrRef `json:"params" yaml:"params"`
	Result *SchemaOrRef `json:"result" yaml:"result"`
	Errors []string     `json:"errors" yaml:"errors"`
}

type ErrorSchema struct {
	Code        int          `json:"code"        yaml:"code"`
	Description *string      `json:"description" yaml:"description"`
	Data        *SchemaOrRef `json:"data"        yaml:"data"`
}

type SchemaOrRef struct {
	Ref        *string                `json:"$ref"       yaml:"$ref"`
	Type       *string                `json:"type"       yaml:"type"`
	Format     *string                `json:"format"     yaml:"format"`
	Enum       []string               `json:"enum"       yaml:"enum"`
	Items      *SchemaOrRef           `json:"items"      yaml:"items"`
	Properties map[string]SchemaOrRef `json:"properties" yaml:"properties"`
	Required   []string               `json:"required"   yaml:"required"`
}
