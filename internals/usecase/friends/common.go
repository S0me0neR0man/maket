package friends

import (
	"errors"
	"slices"
	"sort"

	"maket/internals/domain/entities"
	"maket/internals/domain/usecase"
	"maket/internals/domain/usecase/weighting"
)

type Friends struct {
	slicer entities.GraphSlicer
	self   entities.ID

	graph    entities.Graph
	weighted usecase.Weighted
	render   func() ([]entities.ID, error)
	err      error
}

func NewFriendsRepoSubset(self entities.ID, slicer entities.GraphSlicer) *Friends {
	return &Friends{
		self:   self,
		slicer: slicer,
	}
}

func (f *Friends) FriendsOfFriend(friend entities.ID) *Friends {
	f.err = nil
	f.render = nil
	f.weighted = nil

	f.graph = f.slicer.Subset(f.self, 2)

	if !slices.Contains(f.graph[f.self], friend) {
		f.err = errors.New("not friend")
		return f
	}

	f.render = func() ([]entities.ID, error) {
		return f.renderFOF(friend)
	}

	return f
}

func (f *Friends) CalcWeight(weighter usecase.Weighter) *Friends {
	f.weighted = weighter.ToWeight(f.graph)
	f.weighted.Print()

	return f
}

func (f *Friends) renderFOF(friend entities.ID) ([]entities.ID, error) {
	if f.err != nil {
		return nil, f.err
	}

	ids := make([]entities.ID, 0)
	for _, person := range f.graph[friend] {
		if person == friend || person == f.self {
			continue
		}
		ids = append(ids, person)
	}

	if f.weighted != nil {
		sort.Slice(ids, func(i, j int) bool {
			a, b := ids[i], ids[j]
			return f.weighted[entities.F{f.self, b}] < f.weighted[entities.F{f.self, a}]
		})
	}

	return ids, nil
}

func (f *Friends) GetIDs() ([]entities.ID, error) {
	if f.render == nil {
		return nil, errors.New("nothing to render")
	}
	if f.err != nil {
		return nil, f.err
	}

	return f.render()
}

func GetFriendsOfFriend(self entities.ID, friend entities.ID, slicer entities.GraphSlicer) ([]entities.ID, error) {
	subset := NewFriendsRepoSubset(self, slicer)
	byCommonFriends := weighting.NewByCommonFriends()

	return subset.
		FriendsOfFriend(friend).
		CalcWeight(byCommonFriends).
		GetIDs()
}
