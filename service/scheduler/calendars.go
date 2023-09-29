package scheduler

import (
	"project-manager-go/api/context"
	"project-manager-go/data"
	"project-manager-go/data/scheduler"
)

type calendars struct {
	store *scheduler.SchedulerStore
}

func (s *calendars) GetAll(userCtx context.UserContext, dbCtx *data.DBContext) (arr []Calendar, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	list := data.NewProjectsList[Calendar]()
	err = s.store.Calendars.GetAll(dbCtx, list)
	if err != nil {
		return nil, err
	}

	return list.GetArray(), nil
}

func (s *calendars) Add(userCtx context.UserContext, dbCtx *data.DBContext, calendar Calendar) (id int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	id, err = s.store.Calendars.Add(dbCtx, &calendar)

	return id, err
}

func (s *calendars) Update(userCtx context.UserContext, dbCtx *data.DBContext, id int, calendar Calendar) (err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	p := Calendar{}
	err = s.store.Calendars.GetOne(dbCtx, id, &p)
	if err != nil {
		return err
	}

	err = s.store.Calendars.Update(dbCtx, id, &calendar)

	return err
}

func (s *calendars) Delete(userCtx context.UserContext, dbCtx *data.DBContext, id int) (children []int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	children, err = s.store.Calendars.Delete(dbCtx, id)

	return children, err
}
