package csharp_test

import (
	"io/ioutil"
	"testing"

	"github.com/glynternet/osc-proto/pkg/generate/csharp"
	"github.com/glynternet/osc-proto/pkg/generate/generatetest"
	"github.com/glynternet/osc-proto/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyTypesShouldYieldEmptyFile(t *testing.T) {
	in := types.Types{}
	var expected map[string][]byte
	out, err := csharp.Generator{}.Generate(in)
	require.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestNonboolFieldTypeShouldYieldError(t *testing.T) {
	in := types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "nonbool",
		}},
	}
	_, err := csharp.Generator{Namespace: "foo"}.Generate(in)
	require.EqualError(t, err, "type:foo has unsupported field type:nonbool for field:fieldFoo")
}

func TestSingleTypeSingleFieldShouldYieldResult(t *testing.T) {
	in := types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "bool",
		}},
	}
	out, err := csharp.Generator{Namespace: "namespaceBar"}.Generate(in)
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
		}},
	}
	out, err := csharp.Generator{
		OSCProtoVersion: "\U0001F9E8",
		Namespace:       "namespaceBar",
	}.Generate(in)
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
	out, err := csharp.Generator{Namespace: "namespaceBar"}.Generate(in)
	require.NoError(t, err)
	generatetest.AssertEqualContentLayout(t, map[string][]byte{
		"namespaceBar.cs": testData(t, "multiple_types.cs"),
	}, out)
}

func testData(t *testing.T, filename string) []byte {
	expected, err := ioutil.ReadFile("testdata/" + filename)
	require.NoError(t, err)
	return expected
}
