package csharp_test

import (
	"io/ioutil"
	"testing"

	"github.com/glynternet/osc-proto/pkg/generate"
	"github.com/glynternet/osc-proto/pkg/generate/csharp"
	"github.com/glynternet/osc-proto/pkg/generate/generatetest"
	"github.com/glynternet/osc-proto/pkg/routers"
	"github.com/glynternet/osc-proto/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyTypesShouldYieldEmptyFile(t *testing.T) {
	in := types.Types{}
	var expected map[string][]byte
	out, err := csharp.Generator{}.Generate(generate.Definitions{Types: in})
	require.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestNoFieldsTypeShouldYieldCorrectEmptyStruct(t *testing.T) {
	in := types.Types{
		"bar": {},
	}
	out, err := csharp.Generator{Namespace: "namespaceBar"}.Generate(generate.Definitions{Types: in})
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string][]byte{
		"namespaceBar.cs": testData(t, "type_with_no_fields.cs"),
	}, out)
}

func TestNonboolFieldTypeShouldYieldError(t *testing.T) {
	in := types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "nonbool",
		}},
	}
	_, err := csharp.Generator{Namespace: "foo"}.Generate(generate.Definitions{Types: in})
	require.EqualError(t, err, "generating typeTemplateVars: type:foo has error: unsupported field type:nonbool for field:fieldFoo")
}

func TestSingleTypeSingleFieldShouldYieldResult(t *testing.T) {
	in := types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "bool",
		}},
	}
	out, err := csharp.Generator{Namespace: "namespaceBar"}.Generate(generate.Definitions{Types: in})
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string][]byte{
		"namespaceBar.cs": testData(t, "single_type_single_field.cs"),
	}, out)
}

func TestSingleTypeMultipleFieldsShouldYieldResult(t *testing.T) {
	in := types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "bool",
		}, {
			FieldName: "fieldBar",
			FieldType: "int32",
		}, {
			FieldName: "fieldString",
			FieldType: "string",
		}, {
			FieldName: "fieldFloat32",
			FieldType: "float32",
		}},
	}
	out, err := csharp.Generator{
		OSCProtoVersion: "\U0001F9E8",
		Namespace:       "namespaceBar",
	}.Generate(generate.Definitions{Types: in})
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string][]byte{
		"namespaceBar.cs": testData(t, "single_type_multiple_fields.cs"),
	}, out)
}

func TestMultipleTypesShouldYieldResult(t *testing.T) {
	in := types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "bool",
		}},
		"bar": {{
			FieldName: "fieldBar",
			FieldType: "bool",
		}},
	}
	out, err := csharp.Generator{Namespace: "namespaceBar"}.Generate(generate.Definitions{Types: in})
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string][]byte{
		"namespaceBar.cs": testData(t, "multiple_types.cs"),
	}, out)
}

func TestSingleTypeSingleFieldSingleRouterShouldYieldResult(t *testing.T) {
	out, err := csharp.Generator{Namespace: "namespaceBar"}.Generate(generate.Definitions{
		Types: types.Types{
			"foo": {{
				FieldName: "fieldFoo",
				FieldType: "bool",
			}},
		},
		Routers: map[routers.RouterName]routers.Routes{
			"bar": {
				"baz":   "foo",
				"whoop": "foo",
			},
		},
	})
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string][]byte{
		"namespaceBar.cs": testData(t, "single_type_with_single_router.cs"),
	}, out)
}

func TestSingleTypeSingleFieldMultipleRouterShouldYieldResult(t *testing.T) {
	out, err := csharp.Generator{Namespace: "namespaceBar"}.Generate(generate.Definitions{
		Types: types.Types{
			"foo": {{
				FieldName: "fieldFoo",
				FieldType: "bool",
			}},
		},
		Routers: map[routers.RouterName]routers.Routes{
			"bar": {
				"baz":   "foo",
				"whoop": "foo",
			},
			"baz": {
				"baz":   "foo",
				"whoop": "foo",
			},
		},
	})
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string][]byte{
		"namespaceBar.cs": testData(t, "single_type_with_multiple_routers.cs"),
	}, out)
}

func TestEnumGeneration(t *testing.T) {
	out, err := csharp.Generator{Namespace: "TestNamespace"}.Generate(generate.Definitions{
		Enums: types.Enums{
			"LogLevel": types.EnumDefinition{
				Values: []types.EnumValue{"DEBUG", "INFO", "ERROR"},
			},
		},
		Types: types.Types{
			"LogMessage": types.TypeFields{{
				FieldName: "level",
				FieldType: "LogLevel",
			}},
		},
	})
	require.NoError(t, err)
	require.Len(t, out, 1)
	
	content := string(out["TestNamespace.cs"])
	assert.Contains(t, content, "public enum LogLevel {")
	assert.Contains(t, content, "DEBUG = 0,")
	assert.Contains(t, content, "INFO = 1,")
	assert.Contains(t, content, "ERROR = 2")
	assert.Contains(t, content, "public static class LogLevelExtensions")
	assert.Contains(t, content, "public static LogLevel ParseLogLevel(string value)")
	assert.Contains(t, content, "LogLevelExtensions.ParseLogLevel")
}

func TestEmptyEnumsAndTypesShouldYieldNil(t *testing.T) {
	out, err := csharp.Generator{}.Generate(generate.Definitions{})
	require.NoError(t, err)
	assert.Nil(t, out)
}

func TestEnumFieldInRouter(t *testing.T) {
	out, err := csharp.Generator{Namespace: "Test"}.Generate(generate.Definitions{
		Enums: types.Enums{
			"PlayerState": types.EnumDefinition{
				Values: []types.EnumValue{"IDLE", "RUNNING"},
			},
		},
		Types: types.Types{
			"Player": types.TypeFields{{
				FieldName: "state",
				FieldType: "PlayerState",
			}},
		},
		Routers: map[routers.RouterName]routers.Routes{
			"game": {
				"updatePlayer": "Player",
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, out, 1)
	
	content := string(out["Test.cs"])
	assert.Contains(t, content, "void HandleUpdatePlayer(Player Player)")
	assert.Contains(t, content, "private readonly PlayerState _state")
	assert.Contains(t, content, "public Player(PlayerState state)")
	assert.Contains(t, content, "PlayerStateExtensions.ParsePlayerState")
}

func testData(t *testing.T, filename string) []byte {
	expected, err := ioutil.ReadFile("testdata/" + filename)
	require.NoError(t, err)
	return expected
}
