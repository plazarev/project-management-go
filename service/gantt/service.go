package gantt

import (
	"project-manager-go/data/gantt"
	"project-manager-go/service"
)

type GanttService struct {
	Tasks *tasks
	Links *links
}

func NewGanttService(store *gantt.GanttStore, tree *service.TreeService) *GanttService {
	return &GanttService{
		Tasks: &tasks{
			tree:  tree,
			store: store,
		},
		Links: &links{store},
	}
}
