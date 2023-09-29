package kanban

import (
	"project-manager-go/api/context"
	"project-manager-go/data"
	"project-manager-go/data/kanban"
)

type cols struct {
	store *kanban.KanbanStore
}

func (s *cols) GetAll(userCtx context.UserContext, dbCtx *data.DBContext) (arr []kanban.KanbanColumn, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	arr, err = s.store.Columns.GetAll(dbCtx)

	return arr, err
}

func (s *cols) Add(userCtx context.UserContext, dbCtx *data.DBContext, column kanban.KanbanColumn) (id int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	index, err := s.store.Columns.GetMaxIndex(dbCtx)
	if err != nil {
		return 0, err
	}

	column.Index = index

	id, err = s.store.Columns.Add(dbCtx, column)

	return id, err
}

func (s *cols) Update(userCtx context.UserContext, dbCtx *data.DBContext, id int, upd kanban.KanbanColumn) (err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	if upd.Label == "" {
		return nil
	}

	col, err := s.store.Columns.GetOne(dbCtx, id)
	if err != nil {
		return err
	}

	col.Label = upd.Label

	err = s.store.Columns.Update(dbCtx, id, col)

	return err
}

func (s *cols) Delete(userCtx context.UserContext, dbCtx *data.DBContext, id int) (children []int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	children, err = s.store.Columns.Delete(dbCtx, id)

	return children, err
}

func (s *cols) Move(userCtx context.UserContext, dbCtx *data.DBContext, id, before int) (err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	column, err := s.store.Columns.GetOne(dbCtx, id)
	if err != nil {
		return err
	}

	from := column.Index
	columnBefore := kanban.KanbanColumn{}
	to := 0

	if before == 0 {
		// move to the end
		to, err = s.store.Columns.GetMaxIndex(dbCtx)
	} else {
		columnBefore, err = s.store.Columns.GetOne(dbCtx, before)
		to = columnBefore.Index
	}
	if err != nil {
		return err
	}

	if from < to {
		// move down: should find the previous card and swap with it
		to--
		columnBefore, err = s.store.Columns.GetByIndex(dbCtx, to)
		if err != nil {
			return err
		}
	}

	columnBefore.Index = column.Index
	column.Index = to

	err = s.store.Columns.Update(dbCtx, column.ID, column)
	if err != nil {
		return err
	}

	err = s.store.Columns.Update(dbCtx, columnBefore.ID, columnBefore)
	if err != nil {
		return err
	}

	return err
}
