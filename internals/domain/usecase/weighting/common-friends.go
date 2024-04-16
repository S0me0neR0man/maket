package weighting

import (
	"maket/internals/domain/entities"
	"maket/internals/domain/usecase"
)

type ByCommonFriends struct {
}

func NewByCommonFriends() *ByCommonFriends {
	return &ByCommonFriends{}
}

func (w *ByCommonFriends) ToWeight(graph entities.Graph) usecase.Weighted {
	visited := make(map[entities.ID]bool)
	weighted := make(usecase.Weighted)

	for person, friends := range graph {
		if visited[person] {
			continue
		}
		for _, friend := range friends {
			friendsOfFriend := graph[friend]
			commonFriends := w.countCommonFriends(friends, friendsOfFriend)
			weighted[entities.F{person, friend}] = usecase.Weight(commonFriends)
			weighted[entities.F{friend, person}] = usecase.Weight(commonFriends)
		}
	}

	return weighted
}

func (w *ByCommonFriends) countCommonFriends(a, b []entities.ID) int {
	m := make(map[entities.ID]bool)

	for _, value := range a {
		m[value] = true
	}

	result := 0
	for _, value := range b {
		if !m[value] {
			continue
		}
		result++
	}

	return result
}
