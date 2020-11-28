package types

import (
	"fmt"
	"strings"
)

type Types map[TypeName]TypeFields

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
