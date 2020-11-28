package golang_test

import (
	"testing"

	"github.com/glynternet/osc-proto/pkg/generate/golang"
	"github.com/glynternet/osc-proto/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyTypesShouldYieldEmptyFile(t *testing.T) {
	in := types.Types{}
	var expected []byte
	out, err := golang.Generator{}.Generate(in)
	require.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestMultipleTypesShouldYieldUnsupportedError(t *testing.T) {
	in := types.Types{
		"foo": {},
		"bar": {},
	}

	_, err := golang.Generator{}.Generate(in)
	require.Error(t, err)
}

func TestSingleTypeShouldYieldResult(t *testing.T) {
	in := types.Types{
		"foo": {{
			FieldName: "fieldFoo",
			FieldType: "fieldFooType",
		}},
	}
	out, err := golang.Generator{Package: "packageBar"}.Generate(in)
	const expected = `package packageBar

type Foo struct {
	fieldFoo fieldFooType
}

func (f Foo) MessageArgs() []interface{} {
	return []interface{
		f.fieldFoo
	}
}
`
	require.NoError(t, err)
	assert.Equal(t, expected, string(out))
}
