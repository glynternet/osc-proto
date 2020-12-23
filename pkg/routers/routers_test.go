package routers

import (
	"testing"

	"github.com/glynternet/osc-proto/pkg/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestRoutersYAMLAPI(t *testing.T) {
	in := `---
foo:
  bar: barType
  baz: bazType
woop:
  shoop: shoopType
  doop: doopType
`

	expected := Routers{
		"foo": map[RouteName]types.TypeName{
			"bar": types.TypeName("barType"),
			"baz": types.TypeName("bazType"),
		},
		"woop": map[RouteName]types.TypeName{
			"shoop": types.TypeName("shoopType"),
			"doop":  types.TypeName("doopType"),
		},
	}

	var out Routers
	require.NoError(t, yaml.Unmarshal([]byte(in), &out))
	require.Equal(t, expected, out)
}
