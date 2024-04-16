package usecase

import "maket/internals/domain/entities"

type FriendsRepository interface {
	GetFriendsOfFriend(self entities.ID, friend entities.ID, slicer entities.GraphSlicer) ([]entities.ID, error)
}
