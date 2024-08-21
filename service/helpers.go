package service

import (
	"project-manager-go/data"
)

func UsersToIDs[T ~int](users []data.User) []T {
	if len(users) == 0 {
		return nil
	}

	ids := make([]T, len(users))
	for i := range users {
		ids[i] = T(users[i].ID)
	}

	return ids
}

func IDsToUsers[T ~int](ids []T) []data.User {
	if len(ids) == 0 {
		return nil
	}

	users := make([]data.User, len(ids))
	for i := range users {
		users[i] = data.User{
			ID: int(ids[i]),
		}
	}

	return users
}