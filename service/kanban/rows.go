package kanban

import (
	"project-manager-go/api/context"
	"project-manager-go/data"
	"project-manager-go/data/kanban"
)

type rows struct {
	store *kanban.KanbanStore
}

func (s *rows) GetAll(userCtx context.UserContext, dbCtx *data.DBContext) (arr []Row, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	list := data.NewProjectsList[Row]()
	err = s.store.Rows.GetAll(dbCtx, list)
	if err != nil {
		return nil, err
	}

	return list.GetArray(), err
}

func (s *rows) Add(userCtx context.UserContext, dbCtx *data.DBContext, row Row) (id int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	index, err := s.store.Rows.GetMaxIndex(dbCtx)
	if err != nil {
		return 0, err
	}

	row.Index = index

	id, err = s.store.Rows.Add(dbCtx, &row)

	return id, err
}

func (s *rows) Update(userCtx context.UserContext, dbCtx *data.DBContext, id int, upd Row) (err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	if upd.Label == "" {
		return nil
	}

	p := Row{}
	err = s.store.Rows.GetOne(dbCtx, id, &p)
	if err != nil {
		return err
	}

	p.Label = upd.Label

	err = s.store.Rows.Update(dbCtx, id, &p)

	return err
}

func (s *rows) Delete(userCtx context.UserContext, dbCtx *data.DBContext, id int) (children []int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	children, err = s.store.Rows.Delete(dbCtx, id)

	return children, err
}

func (s *rows) Move(userCtx context.UserContext, dbCtx *data.DBContext, id, before int) (err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	row := Row{}
	err = s.store.Rows.GetOne(dbCtx, id, &row)
	if err != nil {
		return err
	}

	from := row.Index
	rowBefore := Row{}
	to := 0

	if before == 0 {
		// move to the end
		to, err = s.store.Rows.GetMaxIndex(dbCtx)
	} else {
		err = s.store.Rows.GetOne(dbCtx, before, &rowBefore)
		to = rowBefore.Index
	}
	if err != nil {
		return err
	}

	if from < to {
		// move down: should find the previous card and swap with it
		to--
		err = s.store.Rows.GetByIndex(dbCtx, to, &rowBefore)
		if err != nil {
			return err
		}
	}

	rowBefore.Index = row.Index
	row.Index = to

	err = s.store.Rows.Update(dbCtx, int(row.ID), &row)
	if err != nil {
		return err
	}

	err = s.store.Rows.Update(dbCtx, int(rowBefore.ID), &rowBefore)
	if err != nil {
		return err
	}

	return err
}
