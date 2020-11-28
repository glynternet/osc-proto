package golang

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/glynternet/osc-proto/pkg/types"
	types2 "github.com/glynternet/osc-proto/pkg/types"
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

func (g Generator) Generate(types types.Types) ([]byte, error) {
	if len(types) == 0 {
		return nil, nil
	}
	if len(types) > 1 {
		return nil, errors.New("only generating for a single type is supported currently")
	}
	var out bytes.Buffer
	for name, fields := range types {
		if err := tmpl.Execute(&out, struct {
			Package            string
			TypeName           types2.TypeName
			TypeMethodReceiver string
			FieldName          types2.FieldName
			FieldType          types2.FieldType
		}{
			Package:            g.Package,
			TypeName:           types2.TypeName(strings.Title(string(name))),
			TypeMethodReceiver: strings.ToLower(string(name[0])),
			FieldName:          fields[0].FieldName,
			FieldType:          fields[0].FieldType,
		}); err != nil {
			return nil, errors.Wrap(err, "executing template")
		}
	}
	return out.Bytes(), nil
}
