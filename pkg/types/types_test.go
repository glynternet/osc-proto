package types_test

import (
	"testing"

	"github.com/glynternet/osc-proto/pkg/types"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestSingleFieldYamlUnmarshal(t *testing.T) {
	in := `
blur:
  - enabled: bool`

	expected := types.Types{
		"blur": types.TypeFields{
			types.FieldDefinition{
				FieldName: "enabled",
				FieldType: "bool",
			},
		},
	}

	var out types.Types
	err := yaml.UnmarshalStrict([]byte(in), &out)
	require.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestMultipleKeysInSingleFieldYamlUnmarshalShouldFail(t *testing.T) {
	in := `
blur:
  - enabled: bool
    foo: bar`

	var out types.Types
	err := yaml.UnmarshalStrict([]byte(in), &out)
	require.EqualError(t, err, "expected a single field for FieldDefinition but got 2: enabled foo")
}
