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
