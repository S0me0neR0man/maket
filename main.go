package main

import (
	"fmt"
	"sync"
)

// F - friendship
type F [2]string

type Graph map[string][]string

func (g Graph) Print(title string) {
	println(title)
	for k, vs := range g {
		fmt.Printf("%v: %v\n", k, vs)
	}
}

type SocialGraphSlicer interface {
	Subset(startFrom string, deep int) Graph
}

type SocialGraphCache struct {
	m sync.Map
}

func NewSocialGraphCache() *SocialGraphCache {
	return &SocialGraphCache{m: sync.Map{}}
}

func (g *SocialGraphCache) Store(friendship ...F) {
	for _, pair := range friendship {
		aName, bName := pair[0], pair[1]
		aFriends, bFriends := g.Load(aName), g.Load(bName)
		aFriends = append(aFriends, bName)
		bFriends = append(bFriends, aName)
		g.m.Store(pair[0], aFriends)
		g.m.Store(pair[1], bFriends)
	}
}

func (g *SocialGraphCache) Load(key string) []string {
	if d, ok := g.m.Load(key); ok {
		if result, okToSlice := d.([]string); okToSlice {
			return result
		}
	}

	return nil
}

func (g *SocialGraphCache) Print() {
	println("-- SocialGraphCache")
	g.m.Range(func(key, value interface{}) bool {
		fmt.Printf("%s = %v\n", key, value)
		return true
	})
}

func (g *SocialGraphCache) Subset(startFrom string, deep int) Graph {
	visited := make(map[string]bool)
	subset := make(Graph)

	friends := g.Load(startFrom)
	subset[startFrom] = friends

	for i := 0; i < deep && len(friends) > 0; i++ {
		neighbors := make([]string, 0)
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

type FriendComparator interface {
	Less(a, b string) bool
}

type Weighted map[string]int

func (w Weighted) Less(a, b string) bool {
	return w[a] < w[b]
}

type WeighingByCommonFriends struct {
	slicer SocialGraphSlicer
}

func NewWeighingByCommonFriends(self, b string, slicer SocialGraphSlicer) *WeighingByCommonFriends {
	return &WeighingByCommonFriends{slicer: slicer}
}

func (w *WeighingByCommonFriends) ToWeight(self, friend string) Weighted {
	graph := w.slicer.Subset(self, 1)
	graph.Print(fmt.Sprintf("--- Graph self='%s' deep=1", self))

	return nil
}

func main() {
	//slicer := newSocialGraphSlicer()
}

func newSocialGraphSlicer() SocialGraphSlicer {
	cache := NewSocialGraphCache()
	cache.Store(
		F{"main", "friend"},
		F{"main", "one"}, F{"main", "two"},
		F{"friend", "one"}, F{"friend", "three"}, F{"friend", "five"},
		F{"other1", "other2"}, F{"other3", "other2"})
	cache.Print()

	return cache
}
