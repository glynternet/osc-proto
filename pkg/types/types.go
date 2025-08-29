package types

import (
	"fmt"
	"sort"
	"strings"
)

type Types map[TypeName]TypeFields

func (ts Types) SortedNames() []string {
	var names []string
	for name := range ts {
		names = append(names, string(name))
	}
	sort.Strings(names)
	return names
}

type TypeName string

type TypeFields []FieldDefinition

type FieldDefinition struct {
	FieldName
	FieldType
}

func (f *FieldDefinition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := make(map[FieldName]FieldType)
	if err := unmarshal(m); err != nil {
		return err
	}
	if fieldCount := len(m); fieldCount != 1 {
		fields := func() string {
			var names []string
			for name := range m {
				names = append(names, string(name))
			}
			// sort for deterministic error message
			sort.Strings(names)
			return strings.Join(names, " ")
		}()
		return fmt.Errorf("expected a single field for FieldDefinition but got %d: %s", fieldCount, fields)
	}
	for name, fieldType := range m {
		*f = FieldDefinition{
			FieldName: name,
			FieldType: fieldType,
		}
	}
	return nil
}

type FieldName string

type FieldType string

type Enums map[EnumName]EnumDefinition

func (es Enums) SortedNames() []string {
	var names []string
	for name := range es {
		names = append(names, string(name))
	}
	sort.Strings(names)
	return names
}

type EnumName string

type EnumDefinition struct {
	Values []EnumValue `yaml:"values"`
}

type EnumValue string

func (es Enums) Validate() error {
	for name, def := range es {
		if err := def.Validate(name); err != nil {
			return err
		}
	}
	return nil
}

func (ed EnumDefinition) Validate(name EnumName) error {
	if len(ed.Values) == 0 {
		return EmptyEnumDefinition{EnumName: name}
	}
	
	seen := make(map[EnumValue]bool)
	for _, value := range ed.Values {
		if seen[value] {
			return DuplicateEnumValue{EnumName: name, EnumValue: value}
		}
		seen[value] = true
	}
	return nil
}

func (ts Types) ValidateWithEnums(enums Enums) error {
	primitiveTypes := map[FieldType]bool{
		"bool":    true,
		"int32":   true,
		"float32": true,
		"string":  true,
	}
	
	for typeName, fields := range ts {
		for _, field := range fields {
			if !primitiveTypes[field.FieldType] {
				// Check if it's a valid enum
				enumName := EnumName(field.FieldType)
				if _, exists := enums[enumName]; !exists {
					return TypeError{
						TypeName: typeName,
						Err: UnknownEnumType{
							EnumName:  enumName,
							FieldName: field.FieldName,
						},
					}
				}
			}
		}
	}
	return nil
}
