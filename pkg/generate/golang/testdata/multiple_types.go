package packageName

// Code generated by osc-proto (version unknown) DO NOT EDIT.

func BarMessageArgs(fieldBool bool, fieldInt32 int32, fieldString string, fieldFloat32 float32) []interface{} {
	return []interface{}{
		boolInt32(fieldBool),
		fieldInt32,
		fieldString,
		fieldFloat32,
	}
}

func FooMessageArgs() []interface{} {
	return []interface{}{}
}

// boolInt32 returns an int32 representation of a bool.
// This is required for supporting OSC frameworks that
// don't support a boolean primitive
func boolInt32(value bool) int32 {
	if value {
		return 1
	}
	return 0
}
