package entities

import "fmt"

type ID string

// F friendship
type F [2]ID

type Graph map[ID][]ID

func (g Graph) Print(title string) {
	println(title)
	for k, vs := range g {
		fmt.Printf("%v: %v\n", k, vs)
	}
}

type GraphSlicer interface {
	Subset(startFrom ID, deep int) Graph
}
