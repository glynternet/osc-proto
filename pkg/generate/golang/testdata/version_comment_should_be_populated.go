package packageBar

// Code generated by osc-proto (version 🧨) DO NOT EDIT.

func FooMessageArgs(fieldFoo bool) []interface{} {
	return []interface{}{
		boolInt32(fieldFoo),
	}
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