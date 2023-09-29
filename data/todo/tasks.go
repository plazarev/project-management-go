package todo

import (
	"project-manager-go/data"
)

type tasks struct {
	*data.TreeStore
}

func (s *tasks) GetAll(ctx *data.DBContext, list data.IItemsList) error {
	items := make([]data.Item, 0)
	err := ctx.DB.
		Preload("AssignedUsers").
		Order("project_id, parent_id, `index`").
		Find(&items).
		Error
	if err != nil {
		return err
	}

	list.PutItems(items)

	return nil
}

func (s *tasks) GetByProject(ctx *data.DBContext, projectId int, list data.IItemsList) error {
	items := make([]data.Item, 0)
	err := ctx.DB.
		Preload("AssignedUsers").
		Order("project_id, parent_id, `index`").
		Find(&items, "project_id = ?", projectId).
		Error
	if err != nil {
		return err
	}

	list.PutItems(items)

	return nil
}
