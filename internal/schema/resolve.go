package schema

import (
	"fmt"
	"strings"
)

func ResolveRef(spec *Spec, ref string) (SchemaOrRef, string, error) {
	const prefix = "#/schemas/"
	if !strings.HasPrefix(ref, prefix) {
		return SchemaOrRef{}, "", fmt.Errorf("unsupported ref: %s", ref)
	}

	name := strings.TrimPrefix(ref, prefix)

	schema, ok := spec.Schemas[name]
	if !ok {
		return SchemaOrRef{}, "", fmt.Errorf("entity %q not found", name)
	}

	return schema, name, nil
}
