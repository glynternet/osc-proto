package types

import "fmt"

type TypeError struct {
	TypeName
	Err error
}

func (t TypeError) Error() string {
	return fmt.Sprintf("type:%s has error: %v", t.TypeName, t.Err)
}

type UnsupportedFieldType struct {
	FieldType
	FieldName
}

func (u UnsupportedFieldType) Error() string {
	return fmt.Sprintf("unsupported field type:%s for field:%s", u.FieldType, u.FieldName)
}

type UnknownEnumType struct {
	EnumName
	FieldName
}

func (u UnknownEnumType) Error() string {
	return fmt.Sprintf("unknown enum type:%s for field:%s", u.EnumName, u.FieldName)
}

type DuplicateEnumValue struct {
	EnumName
	EnumValue
}

func (d DuplicateEnumValue) Error() string {
	return fmt.Sprintf("duplicate enum value:%s in enum:%s", d.EnumValue, d.EnumName)
}

type EmptyEnumDefinition struct {
	EnumName
}

func (e EmptyEnumDefinition) Error() string {
	return fmt.Sprintf("empty enum definition:%s", e.EnumName)
}
