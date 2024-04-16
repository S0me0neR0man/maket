package main

import (
	"fmt"

	"maket/internals/domain/entities"
	cache2 "maket/internals/infrastructure/cache"
	"maket/internals/usecase/friends"
)

func main() {
	slicer := newGraphCache()

	ids, err := friends.GetFriendsOfFriend("main", "friend", slicer)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("\n friends of friend %v\n\n", ids)
}

func newGraphCache() entities.GraphSlicer {
	cache := cache2.NewGraphCache()
	cache.Store(
		entities.F{"main", "friend"},
		entities.F{"main", "one"}, entities.F{"main", "two"},
		entities.F{"friend", "one"}, entities.F{"friend", "three"}, entities.F{"friend", "five"},
		entities.F{"one", "two"},
		entities.F{"other1", "other2"}, entities.F{"other3", "other2"})

	cache.Print()

	return cache
}
