package generate

import (
	"github.com/glynternet/osc-proto/pkg/routers"
	"github.com/glynternet/osc-proto/pkg/types"
)

type Generator interface {
	Generate(definitions Definitions) (map[string][]byte, error)
}

type Definitions struct {
	Types   types.Types     `yaml:"types"`
	Routers routers.Routers `yaml:"routers"`
}
