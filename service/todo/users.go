package todo

import (
	uCtx "project-manager-go/api/context"
	"project-manager-go/data"
	"project-manager-go/data/todo"
)

type TodoUser struct {
	data.User
	Label string `json:"label"`
}

type users struct {
	store *todo.TodoStore
}

func (s *users) GetAll(userCtx uCtx.UserContext, dbCtx *data.DBContext) (arr []TodoUser, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	users, err := s.store.Users.GetAll(dbCtx)
	arr = make([]TodoUser, len(users))
	for i := range users {
		arr[i] = TodoUser{
			User:  users[i],
			Label: users[i].Name,
		}
	}

	return arr, err
}
