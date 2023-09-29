package scheduler

import (
	"project-manager-go/data"
	"project-manager-go/data/scheduler/scheduler_props"
)

type SchedulerStore struct {
	Events    *events
	Calendars *calendars
	Users     *users
}

func InitSchedulerStore(itemsTreeStore *data.TreeStore, projectsStore *data.ProjectsStore) *SchedulerStore {
	return &SchedulerStore{
		Events: &events{
			TreeStore: itemsTreeStore,
		},
		Calendars: &calendars{
			ProjectsStore: projectsStore,
		},
		Users: &users{},
	}
}

func (s *SchedulerStore) HandleProjectAddOperation(ctx *data.DBContext, obj *data.Project) error {
	// the Handler is called before the project is created, default values can be defined here

	if obj.Scheduler_Color == nil {
		obj.Scheduler_Color = &scheduler_props.Color{
			Background: "#5890DC",
			Border:     "#2D74D3",
		}
	}

	obj.Scheduler_Active = true

	return nil
}
