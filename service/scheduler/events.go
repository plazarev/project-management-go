package scheduler

import (
	uCtx "project-manager-go/api/context"
	"project-manager-go/data"
	"project-manager-go/data/scheduler"
	"project-manager-go/service"
)

type events struct {
	tree  *service.TreeService
	store *scheduler.SchedulerStore
}

func (s *events) GetAll(userCtx uCtx.UserContext, dbCtx *data.DBContext) (arr []Event, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	list := data.NewItemsList[Event]()
	err = s.store.Events.GetAll(dbCtx, list)

	return list.GetArray(), err
}

func (s *events) Add(userCtx uCtx.UserContext, dbCtx *data.DBContext, event Event) (id int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	id, err = s.store.Events.Add(dbCtx, &event)

	return id, err
}

func (s *events) Update(userCtx uCtx.UserContext, dbCtx *data.DBContext, id int, upd Event) (changedProject bool, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	old := Event{}
	err = s.store.Events.GetOne(dbCtx, id, &old)
	if err != nil {
		return false, err
	}
	changedProject = old.CalendarID != upd.CalendarID

	if changedProject {
		s.tree.ChangeNodeProject(dbCtx, id, int(upd.CalendarID), nil)
	}

	err = s.store.Events.Update(dbCtx, id, &upd)

	return changedProject, err
}

func (s *events) Delete(userCtx uCtx.UserContext, dbCtx *data.DBContext, id int) (children []int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	children, err = s.store.Events.DeleteCascade(dbCtx, id)

	return children, err
}
