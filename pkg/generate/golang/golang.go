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

// Code generated by osc-proto (version {{.OSCProtoVersion}}) DO NOT EDIT.
{{range .Types}}
func {{.TypeName}}MessageArgs({{.MethodParameters}}) []interface{} {
	return []interface{}{{"{"}}{{.InterfaceSliceElements}}{{"}"}}
}
{{end}}
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
	OSCProtoVersion string
	Package         string
}

func (g Generator) Generate(typesToGenerate types.Types) (map[string][]byte, error) {
	if len(typesToGenerate) == 0 {
		return nil, nil
	}

	type typeTmplVars struct {
		TypeName               types.TypeName
		MethodParameters       string
		InterfaceSliceElements string
	}

	var typeTmplVarss []typeTmplVars
	for _, name := range typesToGenerate.SortedNames() {
		fields := typesToGenerate[types.TypeName(name)]
		sliceElements, err := interfaceSliceElements(fields)
		if err != nil {
			return nil, errors.Wrapf(err, "generating interface slice elements for type:%s", name)
		}
		typeTmplVarss = append(typeTmplVarss, typeTmplVars{
			TypeName:               types.TypeName(strings.Title(name)),
			MethodParameters:       methodParametersString(fields),
			InterfaceSliceElements: sliceElements,
		})
	}

	var out bytes.Buffer
	version := strings.TrimSpace(g.OSCProtoVersion)
	if version == "" {
		version = "unknown"
	}
	if err := tmpl.Execute(&out, struct {
		OSCProtoVersion string
		Package         string
		Types           []typeTmplVars
	}{
		OSCProtoVersion: version,
		Package:         g.Package,
		Types:           typeTmplVarss,
	}); err != nil {
		return nil, errors.Wrap(err, "executing template")
	}

	return map[string][]byte{g.Package + ".go": out.Bytes()}, nil
}

func methodParametersString(fields types.TypeFields) string {
	var params []string
	for _, field := range fields {
		params = append(params, string(field.FieldName)+" "+string(field.FieldType))
	}
	return strings.Join(params, ", ")
}

func fieldArgFuncs() map[types.FieldType]func(string) string {
	return map[types.FieldType]func(string) string{
		// TODO(glynternet): upgrade to support receiving bools in UnityOSC
		//   so we don't have to do this as a boolInt32
		"bool":  func(argName string) string { return fmt.Sprintf("boolInt32(%s)", argName) },
		"int32": func(argName string) string { return argName },
	}
}

func interfaceSliceElements(fields types.TypeFields) (string, error) {
	if len(fields) == 0 {
		return "", nil
	}
	var args []string
	fieldArgFuncs := fieldArgFuncs()
	for _, field := range fields {
		arg, ok := fieldArgFuncs[field.FieldType]
		if !ok {
			return "", unsupportedFieldType{
				FieldType: field.FieldType,
				FieldName: field.FieldName,
			}
		}
		args = append(args, arg(string(field.FieldName))+",")
	}
	return "\n\t\t" + strings.Join(args, "\n\t\t") + "\n\t", nil
}

type unsupportedFieldType struct {
	types.FieldType
	types.FieldName
}

func (u unsupportedFieldType) Error() string {
	return fmt.Sprintf("unsupported field type:%s for field:%s", u.FieldType, u.FieldName)
}
