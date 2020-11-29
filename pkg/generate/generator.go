package generate

import "github.com/glynternet/osc-proto/pkg/types"

type Generator interface {
	Generate(types types.Types) (map[string][]byte, error)
}
