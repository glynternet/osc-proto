package golang_test

import (
	"io/ioutil"
	"testing"

	"github.com/glynternet/osc-proto/pkg/generate"
	"github.com/glynternet/osc-proto/pkg/generate/generatetest"
	"github.com/glynternet/osc-proto/pkg/generate/golang"
	"github.com/glynternet/osc-proto/pkg/routers"
	"github.com/glynternet/osc-proto/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyTypesShouldYieldEmptyFile(t *testing.T) {
	in := types.Types{}
	var expected map[string][]byte
	out, err := golang.Generator{}.Generate(generate.Definitions{Types: in})
	require.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestUnsupportedFieldTypeShouldReturnError(t *testing.T) {
	_, err := golang.Generator{}.Generate(generate.Definitions{Types: types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "nope",
		}},
	}})
	require.EqualError(t, err, "generating typeTmplVars: generating interface slice elements for type:foo: unsupported field type:nope for field:fieldFoo")
}

func TestVersionCommentShouldBePopulated(t *testing.T) {
	in := types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "bool",
		}},
	}
	out, err := golang.Generator{
		OSCProtoVersion: "\U0001F9E8",
		Package:         "packageBar",
	}.Generate(generate.Definitions{Types: in})
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string][]byte{
		"packageBar.go": testData(t, "version_comment_should_be_populated.go"),
	}, out)
}

func TestSingleTypeShouldYieldResult(t *testing.T) {
	in := types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "bool",
		}, {
			FieldName: "fieldBar",
			FieldType: "bool",
		}, {
			FieldName: "fieldBaz",
			FieldType: "int32",
		}},
	}
	out, err := golang.Generator{Package: "packageBar"}.Generate(generate.Definitions{Types: in})
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string][]byte{
		"packageBar.go": testData(t, "single_type.go"),
	}, out)
}

func TestMultipleTypesShouldYieldResult(t *testing.T) {
	in := types.Types{
		"foo": {},
		"bar": {{
			FieldName: "fieldBool",
			FieldType: "bool",
		}, {
			FieldName: "fieldInt32",
			FieldType: "int32",
		}, {
			FieldName: "fieldString",
			FieldType: "string",
		}},
	}

	out, err := golang.Generator{
		Package: "packageName",
	}.Generate(generate.Definitions{Types: in})
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string][]byte{
		"packageName.go": testData(t, "multiple_types.go"),
	}, out)
}

func TestSingleTypeSingleFieldSingleRouterShouldYieldResult(t *testing.T) {
	out, err := golang.Generator{Package: "packageBar"}.Generate(generate.Definitions{
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
		"packageBar.go": testData(t, "single_type_with_single_router.go"),
	}, out)
}

func TestMultipleRoutersShouldYieldResult(t *testing.T) {
	out, err := golang.Generator{Package: "packageBar"}.Generate(generate.Definitions{
		Types: types.Types{
			"foo": {{
				FieldName: "fieldFoo",
				FieldType: "bool",
			}},
			"bar": {{
				FieldName: "fieldBar",
				FieldType: "bool",
			}},
		},
		Routers: map[routers.RouterName]routers.Routes{
			"bar": {
				"whoop": "foo",
			},
			"baz": {
				"whoop": "bar",
			},
		},
	})
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string][]byte{
		"packageBar.go": testData(t, "multiple_routers.go"),
	}, out)
}

func TestRouteReferencingUnknownTypeShouldError(t *testing.T) {
	_, err := golang.Generator{Package: "packageBar"}.Generate(generate.Definitions{
		Types: types.Types{"foo": {}},
		Routers: map[routers.RouterName]routers.Routes{
			"bar": {
				"baz": "nope",
			},
		},
	})
	require.EqualError(t, err, "generating route messages: cannot find fields for type:nope in router:bar route:baz")
}

func testData(t *testing.T, filename string) []byte {
	expected, err := ioutil.ReadFile("testdata/" + filename)
	require.NoError(t, err)
	return expected
}
