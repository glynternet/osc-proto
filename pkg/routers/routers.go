package routers

import "sort"

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

type Routes []string
