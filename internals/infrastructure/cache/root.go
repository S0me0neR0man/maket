package cache

import (
	"fmt"
	"sync"

	"maket/internals/domain/entities"
)

type GraphCache struct {
	m sync.Map
}

func NewGraphCache() *GraphCache {
	return &GraphCache{m: sync.Map{}}
}

func (g *GraphCache) Store(friendship ...entities.F) {
	for _, pair := range friendship {
		aName, bName := pair[0], pair[1]
		aFriends, bFriends := g.Load(aName), g.Load(bName)
		aFriends = append(aFriends, bName)
		bFriends = append(bFriends, aName)
		g.m.Store(pair[0], aFriends)
		g.m.Store(pair[1], bFriends)
	}
}

func (g *GraphCache) Load(key entities.ID) []entities.ID {
	if d, ok := g.m.Load(key); ok {
		if result, okToSlice := d.([]entities.ID); okToSlice {
			return result
		}
	}

	return nil
}

func (g *GraphCache) Print() {
	println(">> GraphCache")
	g.m.Range(func(key, value interface{}) bool {
		fmt.Printf("%s = %v\n", key, value)
		return true
	})
	println("<< end GraphCache")
}

func (g *GraphCache) Subset(startFrom entities.ID, deep int) entities.Graph {
	visited := make(map[entities.ID]bool)
	subset := make(entities.Graph)

	friends := g.Load(startFrom)
	subset[startFrom] = friends

	for i := 0; i < deep && len(friends) > 0; i++ {
		neighbors := make([]entities.ID, 0)
		for _, friend := range friends {
			if visited[friend] {
				continue
			}
			subset[friend] = g.Load(friend)
			neighbors = append(neighbors, subset[friend]...)
			visited[friend] = true
		}

		friends = neighbors
	}

	return subset
}
