package csharp

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/glynternet/osc-proto/pkg/generate"
	"github.com/glynternet/osc-proto/pkg/routers"
	"github.com/glynternet/osc-proto/pkg/types"
	"github.com/pkg/errors"
)

var fileTmpl = func() *template.Template {
	const tmplStr = `using System;
using System.Collections.Generic;

namespace {{.Namespace}} {

    // Code generated by osc-proto (version {{.OSCProtoVersion}}) DO NOT EDIT.
{{template "types" .Types}}{{template "routers" .Routers}}}
`
	t, err := template.New("csharp").Parse(tmplStr)
	if err != nil {
		panic(errors.Wrap(err, "parsing template"))
	}

	t, err = t.Parse(`
{{define "types"}}{{range .}}
    public readonly struct {{.StructName}} {{"{"}}{{if .Fields}}
{{range .Fields}}        private readonly {{.FieldType}} _{{.FieldName}};
{{end}}
        public {{.StructName}}({{.ConstructorParameters}}) {
{{range .Fields}}            _{{.FieldName}} = {{.FieldName}};
{{end}}        }
{{range .Fields}}
        public {{.FieldType}} {{.FieldNameGetter}}() {
            return _{{.FieldName}};
        }
{{end}}    {{end}}{{"}"}}

    public class {{.TypeUnmarshaller}} {

{{range .Fields}}        // <{{.FieldName}}:{{.OriginalFieldType}}>
{{end}}        public {{.StructName}} Unmarshal(List<object> data) {
            if (data.Count != {{len .Fields}}) {
                throw new ArgumentException($"Expected {{len .Fields}} item in arg list but got {data.Count}");
            }
{{range $i, $field := .Fields}}            var {{.FieldName}} = {{.FieldTypeParseFunc}}(data[{{$i}}].ToString());
{{end}}            return new {{.StructName}}({{.UnmarshalledConstructorCallArgs}});
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
{{end}}{{end}}
`)
	if err != nil {
		panic(errors.Wrap(err, "parsing types template"))
	}

	t, err = t.Parse(`
{{define "routers"}}{{if .}}
    public interface IMessageHandler {
        void Handle(List<object> messageArgs);
    }
{{end}}{{range .}}
    public class {{.RouterName}} : IMessageHandler {
{{range .Unmarshallers}}        {{.}}
{{end}}
        private readonly {{.HandlerInterface}} {{.HandlerVar}};

        private {{.RouterName}}() { }

        public {{.RouterName}}({{.HandlerInterface}} {{.RouterConstructorHandlerVar}}) {
            {{.HandlerVar}} = {{.RouterConstructorHandlerVar}};
        }

        public void Handle(List<object> data) {
            var route = data[0].ToString();
            switch (route) {
{{range .RouteCases}}                {{.}}
{{end}}                default:
                    throw new ArgumentOutOfRangeException("route", route, "Unsupported message routing arg");
            }
        }
    }

    public interface {{.HandlerInterface}} {
{{range .HandlerInterfaceMethods}}        {{.}}
{{end}}    }
{{end}}{{end}}
`)
	if err != nil {
		panic(errors.Wrap(err, "parsing routers template"))
	}

	return t
}()

type Generator struct {
	OSCProtoVersion string
	Namespace       string
}

type fieldTemplateVars struct {
	FieldName          types.FieldName
	FieldNameGetter    string
	OriginalFieldType  types.FieldType
	FieldType          string
	FieldTypeParseFunc string
}

type typeTemplateVars struct {
	StructName                      string
	TypeUnmarshaller                string
	ConstructorParameters           string
	UnmarshalledConstructorCallArgs string
	Fields                          []fieldTemplateVars
}

type routerTemplateVars struct {
	RouterName                  string
	Unmarshallers               []string
	RouteHandlers               []string
	HandlerInterface            string
	HandlerVar                  string
	RouterConstructorHandlerVar string
	RouteCases                  []string
	HandlerInterfaceMethods     []string
}

func (g Generator) Generate(definitions generate.Definitions) (map[string][]byte, error) {
	if len(definitions.Types) == 0 {
		return nil, nil
	}
	ttvs, err := generateTypeTemplateVars(definitions, typeConversions())
	if err != nil {
		return nil, errors.Wrap(err, "generating typeTemplateVars")
	}

	rtvs, err := generateRouterTemplateVars(definitions)
	if err != nil {
		return nil, errors.Wrap(err, "generating routerTemplateVars")
	}
	version := strings.TrimSpace(g.OSCProtoVersion)
	if version == "" {
		version = "unknown"
	}
	var out bytes.Buffer
	if err := fileTmpl.Execute(&out, struct {
		OSCProtoVersion string
		Namespace       string
		Types           []typeTemplateVars
		Routers         []routerTemplateVars
	}{
		OSCProtoVersion: version,
		Namespace:       g.Namespace,
		Types:           ttvs,
		Routers:         rtvs,
	}); err != nil {
		return nil, errors.Wrap(err, "executing template")
	}
	return map[string][]byte{
		g.Namespace + ".cs": out.Bytes(),
	}, nil
}

func generateTypeTemplateVars(definitions generate.Definitions, typeConversions map[types.FieldType]typeConversion) ([]typeTemplateVars, error) {
	var ttvs []typeTemplateVars
	for _, name := range definitions.Types.SortedNames() {
		var ftvs []fieldTemplateVars
		ttypeName := types.TypeName(name)
		fields := definitions.Types[ttypeName]
		for _, field := range fields {
			conversions, ok := typeConversions[field.FieldType]
			if !ok {
				return nil, types.TypeError{
					TypeName: ttypeName,
					Err: types.UnsupportedFieldType{
						FieldType: field.FieldType,
						FieldName: field.FieldName,
					},
				}
			}
			ftvs = append(ftvs, fieldTemplateVars{
				FieldName:          field.FieldName,
				FieldNameGetter:    strings.Title(string(field.FieldName)),
				OriginalFieldType:  field.FieldType,
				FieldType:          conversions.ttype,
				FieldTypeParseFunc: conversions.parseDataFieldFunc,
			})
		}

		tn := arg(name)
		ttvs = append(ttvs, typeTemplateVars{
			StructName:                      tn.StructName(),
			TypeUnmarshaller:                tn.Unmarshaller(),
			ConstructorParameters:           constructorParameters(typeConversions, fields),
			UnmarshalledConstructorCallArgs: unmarshalledConstructorCallArgs(fields),
			Fields:                          ftvs,
		})
	}
	return ttvs, nil
}

func generateRouterTemplateVars(definitions generate.Definitions) ([]routerTemplateVars, error) {
	var tvs []routerTemplateVars
	for _, name := range definitions.Routers.SortedNames() {
		routes := definitions.Routers[routers.RouterName(name)]
		routerName := router(name)
		tvs = append(tvs, routerTemplateVars{
			RouterName:                  routerName.Class(),
			HandlerInterface:            routerName.HandlerInterface(),
			HandlerVar:                  routerName.HandlerVar(),
			RouterConstructorHandlerVar: routerName.RouterConstructorHandlerVar(),
			Unmarshallers:               routerUnmarshallers(routes),
			RouteHandlers:               routerRouteHandlers(routes),
			RouteCases:                  routerRouteCases(routerName, routes),
			HandlerInterfaceMethods:     routerHandlerInterfaceMethods(routes),
		})
	}
	return tvs, nil
}

type router routers.RouterName

func (r router) Class() string {
	return strings.Title(string(r)) + "Router"
}

func (r router) HandlerInterface() string {
	return "I" + strings.Title(string(r)) + "Handler"
}

func (r router) HandlerVar() string {
	return "_" + string(r) + "Handler"
}

func (r router) RouterConstructorHandlerVar() string {
	return string(r) + "Handler"
}

func routerUnmarshallers(routes routers.Routes) []string {
	var us []string
	for _, routeName := range routes.SortedNames() {
		routeArg := arg(routes[routers.RouteName(routeName)])
		us = append(us, fmt.Sprintf("private static readonly %s %s = new %s();",
			routeArg.Unmarshaller(),
			routeArgsUnmarshallerVar(routeName, routeArg),
			routeArg.Unmarshaller()))
	}
	return us
}

func routerRouteHandlers(routes routers.Routes) []string {
	var us []string
	for _, routeName := range routes.SortedNames() {
		routeArg := arg(routes[routers.RouteName(routeName)])
		us = append(us, fmt.Sprintf("private readonly %s _%s;",
			routeArgsUnmarshallerVar(routeName, routeArg),
			routeArg.Unmarshaller()))
	}
	return us
}

func routerRouteCases(routerName router, routes routers.Routes) []string {
	var us []string
	for _, routeName := range routes.SortedNames() {
		routeArg := arg(routes[routers.RouteName(routeName)])
		us = append(us, fmt.Sprintf(`case "%s":
                    %s.%s(%s.Unmarshal(data.GetRange(1, data.Count - 1)));
                    break;`,
			routeName,
			routerName.HandlerVar(),
			"Handle"+strings.Title(routeName),
			routeArgsUnmarshallerVar(routeName, routeArg)),
		)
	}
	return us
}

func routerHandlerInterfaceMethods(routes routers.Routes) []string {
	var ms []string
	for _, routeName := range routes.SortedNames() {
		routeArg := arg(routes[routers.RouteName(routeName)])
		ms = append(ms, fmt.Sprintf(`void Handle%s(%s %s);`,
			strings.Title(routeName),
			routeArg.StructName(),
			routeArg,
		))
	}
	return ms
}

func routeArgsUnmarshallerVar(name string, argType arg) string {
	return strings.Title(string(name)) + argType.Unmarshaller()
}

func constructorParameters(conversions map[types.FieldType]typeConversion, fields types.TypeFields) string {
	var params []string
	for _, field := range fields {
		params = append(params, conversions[field.FieldType].ttype+" "+string(field.FieldName))
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

type typeConversion struct {
	ttype              string
	parseDataFieldFunc string
}

func typeConversions() map[types.FieldType]typeConversion {
	return map[types.FieldType]typeConversion{
		"bool": {
			ttype:              "bool",
			parseDataFieldFunc: "ParseBool",
		},
		"int32": {
			ttype:              "int",
			parseDataFieldFunc: "int.Parse",
		},
		"string": {
			ttype: "string",
			// having an empty string should yield parenthesis around the value,
			// which is ugly but should still be compilable and correct.
			parseDataFieldFunc: "",
		},
		"float32": {
			ttype:              "float",
			parseDataFieldFunc: "float.Parse",
		},
	}
}

type arg types.TypeName

func (n arg) StructName() string {
	return strings.Title(string(n))
}

func (n arg) Unmarshaller() string {
	return n.StructName() + "Unmarshaller"
}
