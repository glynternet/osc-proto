package packageBar

import "github.com/hypebeast/go-osc/osc"

// Code generated by osc-proto (version unknown) DO NOT EDIT.

func BarMessageArgs(fieldBar bool) []interface{} {
	return []interface{}{
		boolInt32(fieldBar),
	}
}

func FooMessageArgs(fieldFoo bool) []interface{} {
	return []interface{}{
		boolInt32(fieldFoo),
	}
}

func BarWhoopFoo(fieldFoo bool) osc.Message {
	return osc.Message{
		Address:   "/bar",
		Arguments: append([]interface{}{"whoop"}, FooMessageArgs(fieldFoo)...),
	}
}

func BazWhoopBar(fieldBar bool) osc.Message {
	return osc.Message{
		Address:   "/baz",
		Arguments: append([]interface{}{"whoop"}, BarMessageArgs(fieldBar)...),
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
