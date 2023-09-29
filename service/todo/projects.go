package todo

import (
	"project-manager-go/api/context"
	"project-manager-go/data"
	"project-manager-go/data/todo"
)

type projects struct {
	store *todo.TodoStore
}

func (s *projects) GetAll(userCtx context.UserContext, dbCtx *data.DBContext) (arr []Project, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	list := data.NewProjectsList[Project]()
	err = s.store.Projects.GetAll(dbCtx, list)
	if err != nil {
		return nil, err
	}

	return list.GetArray(), nil
}

func (s *projects) Add(userCtx context.UserContext, dbCtx *data.DBContext, project Project) (id int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	id, err = s.store.Projects.Add(dbCtx, &project)

	return id, err
}

func (s *projects) Update(userCtx context.UserContext, dbCtx *data.DBContext, id int, project Project) (err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	p := Project{}
	err = s.store.Projects.GetOne(dbCtx, id, &p)
	if err != nil {
		return err
	}

	err = s.store.Projects.Update(dbCtx, id, &project)

	return err
}

func (s *projects) Delete(userCtx context.UserContext, dbCtx *data.DBContext, id int) (children []int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	children, err = s.store.Projects.Delete(dbCtx, id)

	return children, err
}
