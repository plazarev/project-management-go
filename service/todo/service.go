package todo

import (
	"project-manager-go/data/todo"
	"project-manager-go/service"
)

type TodoService struct {
	Tasks    *tasks
	Projects *projects
	Users    *users
	Tags     *tags
}

func NewTodoService(store *todo.TodoStore, tree *service.TreeService) *TodoService {
	return &TodoService{
		Tasks:    &tasks{tree, store},
		Projects: &projects{store},
		Users:    &users{store},
		Tags:     &tags{store},
	}
}
