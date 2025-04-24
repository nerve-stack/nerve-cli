package golang

import (
	"io"
	"text/template"

	"github.com/nerve-stack/nerve-cli/internal/codegen/sdk"
	"github.com/nerve-stack/nerve-cli/internal/schema"
)

//	type tmplCtx struct {
//		Package string
//		Version string
//	}

func GenServer(w io.Writer, spec *schema.Spec) error {
	funcMap := template.FuncMap{
		"title": sdk.Title,
	}

	tmpl := template.New("server.tmpl").Funcs(funcMap)

	tmpl, err := tmpl.ParseFS(templates, "templates/server.tmpl")
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(w, spec); err != nil {
		return err
	}

	return nil
}
