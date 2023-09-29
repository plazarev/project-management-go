package todo

import "project-manager-go/data"

type TodoStore struct {
	Tasks    *tasks
	Projects *projects
	Users    *users
	Tags     *tags
}

func InitTodoStore(itemsTreeStore *data.TreeStore, projectsStore *data.ProjectsStore) *TodoStore {
	return &TodoStore{
		Tasks: &tasks{
			TreeStore: itemsTreeStore,
		},
		Projects: &projects{
			ProjectsStore: projectsStore,
		},
		Users: &users{},
		Tags:  &tags{},
	}
}
