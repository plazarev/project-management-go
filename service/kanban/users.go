package kanban

import (
	uCtx "project-manager-go/api/context"
	"project-manager-go/data"
	"project-manager-go/data/kanban"
)

type users struct {
	store *kanban.KanbanStore
}

func (s *users) GetAll(userCtx uCtx.UserContext, dbCtx *data.DBContext) (users []data.User, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	users, err = s.store.Users.GetAll(dbCtx)

	return users, err
}
