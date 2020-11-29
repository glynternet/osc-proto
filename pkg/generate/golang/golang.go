package golang

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/glynternet/osc-proto/pkg/types"
	"github.com/pkg/errors"
)

var tmpl = func() *template.Template {
	const tmplStr = `package {{.Package}}

type {{.TypeName}} struct {
	{{.FieldName}} {{.FieldType}}
}

func ({{.TypeMethodReceiver}} {{.TypeName}}) MessageArgs() []interface{} {
	return []interface{
		{{.TypeMethodReceiver}}.{{.FieldName}}
	}
}
`
	t, err := template.New("golang").Parse(tmplStr)
	if err != nil {
		panic(errors.Wrap(err, "parsing template"))
	}
	return t
}()

type Generator struct {
	Package string
}

func (g Generator) Generate(typesToGenerate types.Types) ([]byte, error) {
	if len(typesToGenerate) == 0 {
		return nil, nil
	}
	if len(typesToGenerate) > 1 {
		return nil, errors.New("only generating for a single type is supported currently")
	}
	var out bytes.Buffer
	for name, fields := range typesToGenerate {
		if err := tmpl.Execute(&out, struct {
			Package            string
			TypeName           types.TypeName
			TypeMethodReceiver string
			FieldName          types.FieldName
			FieldType          types.FieldType
		}{
			Package:            g.Package,
			TypeName:           types.TypeName(strings.Title(string(name))),
			TypeMethodReceiver: strings.ToLower(string(name[0])),
			FieldName:          types.FieldName(strings.Title(string(fields[0].FieldName))),
			FieldType:          fields[0].FieldType,
		}); err != nil {
			return nil, errors.Wrap(err, "executing template")
		}
	}
	return out.Bytes(), nil
}
