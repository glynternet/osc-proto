package golang

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/glynternet/osc-proto/pkg/types"
	"github.com/pkg/errors"
)

var tmpl = func() *template.Template {
	const tmplStr = `package {{.Package}}

func {{.TypeName}}MessageArgs({{.ArgName}} bool) []interface{} {
	return []interface{}{
		{{.GetFieldArg}},
	}
}

// boolInt32 returns an int32 representation of a bool.
// This is required for supporting OSC frameworks that
// don't support a boolean primitive
func boolInt32(value bool) int32 {
	if value {
		return 1
	}
	return 0
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

func (g Generator) Generate(typesToGenerate types.Types) (map[string][]byte, error) {
	if len(typesToGenerate) == 0 {
		return nil, nil
	}
	if len(typesToGenerate) > 1 {
		return nil, errors.New("only generating for a single type is supported currently")
	}
	for name, fields := range typesToGenerate {
		var out bytes.Buffer
		argName := string(fields[0].FieldName)
		if err := tmpl.Execute(&out, struct {
			Package     string
			TypeName    types.TypeName
			ArgName     string
			GetFieldArg string
		}{
			Package:     g.Package,
			TypeName:    types.TypeName(strings.Title(string(name))),
			ArgName:     argName,
			GetFieldArg: boolFieldArg(argName),
		}); err != nil {
			return nil, errors.Wrap(err, "executing template")
		}
		return map[string][]byte{
			string(name) + ".go": out.Bytes(),
		}, nil
	}
	return nil, errors.New("should be unreachable")
}

// TODO(glynternet): upgrade to support receiving bools in UnityOSC so we don't have to do this
func boolFieldArg(argName string) string {
	return fmt.Sprintf("boolInt32(%s)", argName)
}
