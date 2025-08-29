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
  - foo: bar
    enabled: bool
`

	err := yaml.UnmarshalStrict([]byte(in), &types.Types{})
	require.EqualError(t, err, "expected a single field for FieldDefinition but got 2: enabled foo")
}

func TestEnumValidation(t *testing.T) {
	tests := []struct {
		name        string
		enums       types.Enums
		expectError string
	}{
		{
			name: "valid enum",
			enums: types.Enums{
				"PlayerState": types.EnumDefinition{
					Values: []types.EnumValue{"IDLE", "RUNNING", "JUMPING"},
				},
			},
			expectError: "",
		},
		{
			name: "empty enum should fail",
			enums: types.Enums{
				"Empty": types.EnumDefinition{
					Values: []types.EnumValue{},
				},
			},
			expectError: "empty enum definition:Empty",
		},
		{
			name: "duplicate enum values should fail",
			enums: types.Enums{
				"Duplicate": types.EnumDefinition{
					Values: []types.EnumValue{"IDLE", "RUNNING", "IDLE"},
				},
			},
			expectError: "duplicate enum value:IDLE in enum:Duplicate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.enums.Validate()
			if tt.expectError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectError)
			}
		})
	}
}

func TestTypesValidateWithEnums(t *testing.T) {
	enums := types.Enums{
		"PlayerState": types.EnumDefinition{
			Values: []types.EnumValue{"IDLE", "RUNNING"},
		},
	}

	tests := []struct {
		name        string
		types       types.Types
		expectError string
	}{
		{
			name: "valid type with primitive field",
			types: types.Types{
				"Player": types.TypeFields{
					{FieldName: "x", FieldType: "float32"},
				},
			},
			expectError: "",
		},
		{
			name: "valid type with enum field",
			types: types.Types{
				"Player": types.TypeFields{
					{FieldName: "state", FieldType: "PlayerState"},
				},
			},
			expectError: "",
		},
		{
			name: "invalid enum reference should fail",
			types: types.Types{
				"Player": types.TypeFields{
					{FieldName: "state", FieldType: "UnknownEnum"},
				},
			},
			expectError: "type:Player has error: unknown enum type:UnknownEnum for field:state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.types.ValidateWithEnums(enums)
			if tt.expectError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectError)
			}
		})
	}
}
