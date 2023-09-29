package gantt

import "project-manager-go/data"

type tasks struct {
	*data.TreeStore
}

func (s *tasks) GetAll(ctx *data.DBContext, dest data.IItemsList) error {
	tasks := make([]data.Item, 0)
	err := ctx.DB.
		Order("project_id, parent_id, `index`").
		Find(&tasks).
		Error
	if err != nil {
		return err
	}

	dest.PutItems(tasks)

	return nil
}
