package items

import (
	uCtx "project-manager-go/api/context"
	"project-manager-go/data"
	"project-manager-go/service"
)

type items struct {
	tree     *service.TreeService
	items    *data.TreeStore
	projects *data.ProjectsStore
}

func (s *items) GetAll(userCtx uCtx.UserContext, dbCtx *data.DBContext) (arr []Item, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	list := data.NewItemsList[Item]()
	err = s.items.GetAll(dbCtx, list)

	return list.GetArray(), err
}

func (s *items) Add(userCtx uCtx.UserContext, dbCtx *data.DBContext, item Item) (id int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	id, err = s.items.Add(dbCtx, &item)

	proj, err := s.projects.GetFristId(dbCtx)
	if err != nil {
		return 0, nil
	}
	maxIndex, err := s.items.MaxBranchIndex(dbCtx, proj, 0)
	if err != nil {
		return 0, err
	}
	err = s.tree.Move(dbCtx, id, proj, 0, maxIndex+1)
	return id, err
}

func (s *items) Update(userCtx uCtx.UserContext, dbCtx *data.DBContext, id int, item Item) (err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	err = s.items.Update(dbCtx, id, &item)

	return err
}

func (s *items) Delete(userCtx uCtx.UserContext, dbCtx *data.DBContext, id int) (children []int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	children, err = s.items.DeleteCascade(dbCtx, id)

	return children, err
}
