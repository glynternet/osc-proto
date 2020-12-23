package routers

import (
	"sort"

	"github.com/glynternet/osc-proto/pkg/types"
)

type Routers map[RouterName]Routes

func (ts Routers) SortedNames() []string {
	var names []string
	for name := range ts {
		names = append(names, string(name))
	}
	sort.Strings(names)
	return names
}

type RouterName string

type RouteName string

type Routes map[RouteName]types.TypeName

func (ts Routes) SortedNames() []string {
	var names []string
	for name := range ts {
		names = append(names, string(name))
	}
	sort.Strings(names)
	return names
}
