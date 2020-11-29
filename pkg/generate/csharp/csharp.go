package csharp

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/glynternet/osc-proto/pkg/types"
	"github.com/pkg/errors"
)

var tmpl = func() *template.Template {
	const tmplStr = `using System;
using System.Collections.Generic;
using avvaunity.GOH.Unity.Message.Unmarshaller;

namespace {{.Namespace}} {

    public readonly struct {{.TypeName}} {
{{range .Fields}}        private readonly {{.FieldType}} _{{.FieldName}};
{{end}}
        public {{.TypeName}}({{.ConstructorParameters}}) {
{{range .Fields}}            _{{.FieldName}} = {{.FieldName}};
{{end}}        }
{{range .Fields}}
        public bool {{.FieldNameGetter}}() {
            return _{{.FieldName}};
        }
{{end}}    }

    public class {{.TypeName}}Unmarshaller : IMessageUnmarshaller<{{.TypeName}}> {

{{range .Fields}}        // <{{.FieldName}}:{{.FieldType}}>
{{end}}        public {{.TypeName}} Unmarshal(List<object> data) {
            if (data.Count != {{len .Fields}}) {
                throw new ArgumentException($"Expected {{len .Fields}} item in arg list but got {data.Count}");
            }
{{range $i, $field := .Fields}}            var {{.FieldName}} = Parse{{.FieldTypeParseMethodSuffix}}(data[{{$i}}].ToString());
{{end}}            return new {{.TypeName}}({{.UnmarshalledConstructorCallArgs}});
        }

        private static bool ParseBool(string value) {
            try {
                return int.Parse(value) != 0;
            }
            catch (Exception e) {
                throw new ArgumentException($"cannot parse {value} to bool", e);
            }
        }
    }
}
`
	t, err := template.New("csharp").Parse(tmplStr)
	if err != nil {
		panic(errors.Wrap(err, "parsing template"))
	}
	return t
}()

type Generator struct {
	Namespace string
}

func (g Generator) Generate(typesToGenerate types.Types) (map[string][]byte, error) {
	if len(typesToGenerate) == 0 {
		return nil, nil
	}
	if len(typesToGenerate) > 1 {
		return nil, errors.New("only generating for a single type is supported currently")
	}
	for name, fields := range typesToGenerate {
		type fieldTemplateVars struct {
			FieldName                  types.FieldName
			FieldNameGetter            string
			FieldType                  types.FieldType
			FieldTypeParseMethodSuffix string
		}
		var ftvs []fieldTemplateVars
		for _, field := range fields {
			if field.FieldType != "bool" {
				return nil, fmt.Errorf("type:%s contains non-bool field type:%s for field:%s", name, field.FieldType, field.FieldName)
			}
			ftvs = append(ftvs, fieldTemplateVars{
				FieldName:                  field.FieldName,
				FieldNameGetter:            strings.Title(string(field.FieldName)),
				FieldType:                  field.FieldType,
				FieldTypeParseMethodSuffix: strings.Title(string(field.FieldType)),
			})
		}
		var out bytes.Buffer
		if err := tmpl.Execute(&out, struct {
			Namespace                       string
			TypeName                        types.TypeName
			ConstructorParameters           string
			UnmarshalledConstructorCallArgs string
			Fields                          []fieldTemplateVars
		}{
			Namespace:                       g.Namespace,
			TypeName:                        types.TypeName(strings.Title(string(name))),
			ConstructorParameters:           constructorParameters(fields),
			UnmarshalledConstructorCallArgs: unmarshalledConstructorCallArgs(fields),
			Fields:                          ftvs,
		}); err != nil {
			return nil, errors.Wrap(err, "executing template")
		}
		return map[string][]byte{
			string(name) + ".cs": out.Bytes(),
		}, nil
	}
	return nil, errors.New("should be unreachable")
}

func constructorParameters(fields types.TypeFields) string {
	var params []string
	for _, field := range fields {
		params = append(params, string(field.FieldType)+" "+string(field.FieldName))
	}
	return strings.Join(params, ", ")
}

func unmarshalledConstructorCallArgs(fields types.TypeFields) string {
	var args []string
	for _, field := range fields {
		args = append(args, string(field.FieldName))
	}
	return strings.Join(args, ", ")
}
