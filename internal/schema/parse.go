package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func parseJSONSchema(data []byte) (*Spec, error) {
	var spec Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse JSON schema: %w", err)
	}

	return &spec, nil
}

// Function to parse the schema file as YAML
func parseYAMLSchema(data []byte) (*Spec, error) {
	var spec Spec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse YAML schema: %w", err)
	}

	return &spec, nil
}

// Function to determine the schema format based on file extension and parse it
func Parse(fpath string) (*Spec, error) {
	// Read the schema file
	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil, fmt.Errorf("unable to read schema file: %w", err)
	}

	// Get the file extension
	ext := filepath.Ext(fpath)

	// Call the appropriate parsing function based on the file extension
	var spec *Spec

	switch ext {
	case ".json":
		spec, err = parseJSONSchema(data)

	case ".yaml", ".yml":
		spec, err = parseYAMLSchema(data)

	default:
		err = fmt.Errorf("unsupported file extension: %s", ext)
	}

	if err != nil {
		return nil, err
	}

	return spec, nil
}
