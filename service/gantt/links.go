package gantt

import (
	"project-manager-go/api/context"
	"project-manager-go/common"
	"project-manager-go/data"
	"project-manager-go/data/gantt"
)

type Link struct {
	ID     common.FuzzyInt `json:"id"`
	Source common.FuzzyInt `json:"source"`
	Target common.FuzzyInt `json:"target"`
	Type   common.FuzzyInt `json:"type"`
}

type links struct {
	store *gantt.GanttStore
}

func (s *links) GetAll(userCtx context.UserContext, dbCtx *data.DBContext) (arr []gantt.GanttLink, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	arr, err = s.store.Links.GetAll(dbCtx)

	return arr, err
}

func (s *links) Add(userCtx context.UserContext, dbCtx *data.DBContext, link Link) (id int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	upd := gantt.GanttLink{
		ID:     int(link.ID),
		Source: int(link.Source),
		Target: int(link.Target),
		Type:   int(link.Type),
	}

	id, err = s.store.Links.Add(dbCtx, upd)

	return id, err
}

func (s *links) Delete(userCtx context.UserContext, dbCtx *data.DBContext, id int) (err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	err = s.store.Links.Delete(dbCtx, id)

	return err
}
