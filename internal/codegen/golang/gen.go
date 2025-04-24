package golang

import (
	"bytes"
	"text/template"

	"github.com/nerve-stack/nerve-cli/pkg/cases"
)

var funcs = template.FuncMap{
	"camel": cases.ToCamelCase,
}

func renderTemplate(templateName string, data any) (string, error) {
	tmplContent, err := templatesFS.ReadFile(templateName)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(templateName).
		Funcs(funcs).
		Parse(string(tmplContent))
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	err = tmpl.Execute(&buffer, data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func RenderModelToBuffer(model *Model) (map[string]string, error) {
	output := make(map[string]string)

	mainSource, err := renderTemplate("templates/main.tmpl", model)
	if err != nil {
		return nil, err
	}

	output["gen.go"] = mainSource

	return output, nil
}
