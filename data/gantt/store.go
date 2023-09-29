package gantt

import (
	"project-manager-go/data"
)

type GanttStore struct {
	Tasks         *tasks
	Links         *links
	ProjectsStore *data.ProjectsStore
}

func InitGanttStore(itemsTreeStore *data.TreeStore, projectsStore *data.ProjectsStore) *GanttStore {
	db := data.GetDB()

	// Add a custom model to the db. This model only relates to Gantt
	err := db.AutoMigrate(&GanttLink{})
	if err != nil {
		panic(err)
	}

	return &GanttStore{
		Tasks: &tasks{
			TreeStore: itemsTreeStore,
		},
		Links:         &links{},
		ProjectsStore: projectsStore,
	}
}

func (s *GanttStore) HandleTaskDeleteOperation(ctx *data.DBContext, obj *data.Item) error {
	// the Handler is called before the item is deleted, some relations can be cleared here

	err := ctx.DB.Where("source = ? OR target = ?", obj.ID, obj.ID).Delete(&GanttLink{}).Error

	return err
}
