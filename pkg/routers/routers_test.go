package routers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestRoutersYAMLAPI(t *testing.T) {
	in := `---
foo:
  - bar
  - baz
woop:
  - shoop
  - doop
`

	expected := Routers{
		"foo":  []string{"bar", "baz"},
		"woop": []string{"shoop", "doop"},
	}

	var out Routers
	require.NoError(t, yaml.Unmarshal([]byte(in), &out))
	require.Equal(t, expected, out)
}
