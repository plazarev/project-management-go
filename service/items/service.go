package items

import (
	"project-manager-go/data"
	"project-manager-go/service"
)

type BaseItemsService struct {
	Items *items
}

func NewBaseItemsService(tree *service.TreeService, projects *data.ProjectsStore, treeItems *data.TreeStore) *BaseItemsService {
	return &BaseItemsService{
		Items: &items{
			items:    treeItems,
			projects: projects,
			tree:     tree,
		},
	}
}
