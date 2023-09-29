package scheduler

import (
	"project-manager-go/data/scheduler"
	"project-manager-go/service"
)

type SchedulerService struct {
	Events    *events
	Calendars *calendars
	Users     *users
}

func NewSchedulerService(store *scheduler.SchedulerStore, tree *service.TreeService) *SchedulerService {
	return &SchedulerService{
		Events: &events{
			tree:  tree,
			store: store,
		},
		Calendars: &calendars{store},
		Users:     &users{store},
	}
}
