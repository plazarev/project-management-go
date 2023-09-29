package todo

import (
	uCtx "project-manager-go/api/context"
	"project-manager-go/data"
	"project-manager-go/data/todo"
)

type tags struct {
	store *todo.TodoStore
}

func (s *tags) GetAll(userCtx uCtx.UserContext, dbCtx *data.DBContext) (tags []string, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	tags, err = s.store.Tags.GetAll(dbCtx)

	return tags, err
}
