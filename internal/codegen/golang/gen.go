package golang

import (
	"os"
	"text/template"

	"github.com/nerve-stack/nerve-cli/internal/codegen/sdk"
	"github.com/nerve-stack/nerve-cli/internal/schema"
)

// type tmplCtx struct {
// 	Package string
// 	Version string
// }

func Generate(spec *schema.Spec) error {
	funcMap := template.FuncMap{
		"title": sdk.Title,
	}
	// Create a new template with the function map FIRST
	tmpl := template.New("template.tmpl").Funcs(funcMap)

	// Then parse the template content
	tmpl, err := tmpl.ParseFS(templates, "templates/template.tmpl")
	if err != nil {
		return err
	}

	if err := tmpl.Execute(os.Stdout, spec); err != nil {
		return err
	}

	return nil
}
