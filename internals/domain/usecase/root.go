package usecase

import (
	"fmt"

	"maket/internals/domain/entities"
)

type Weight int

type Weighted map[entities.F]Weight

func (w *Weighted) Print() {
	println(">> Weighted")
	for key, value := range *w {
		fmt.Printf("%s = %v\n", key, value)
	}
	println("<< end Weighted")
}

type Weighter interface {
	ToWeight(graph entities.Graph) Weighted
}
