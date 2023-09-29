package scheduler

import (
	"project-manager-go/common"
	"project-manager-go/data"
	"project-manager-go/data/scheduler/scheduler_props"
)

type Calendar struct {
	ID          common.TID `json:"id"`
	Label       string     `json:"label"`
	Active      bool       `json:"active"`
	Description string     `json:"description"`

	Color *scheduler_props.Color `json:"color"`
}

func (c *Calendar) FillProject(proj *data.Project) {
	proj.Label = c.Label
	proj.SchedulerProjectProps = scheduler_props.SchedulerProjectProps{
		Scheduler_Color:       c.Color,
		Scheduler_Active:      c.Active,
		Scheduler_Description: c.Description,
	}
}

func (t *Calendar) PutProject(project data.Project) {
	t.ID = common.TID(project.ID)
	t.Label = project.Label
	t.Color = project.Scheduler_Color
	t.Active = project.Scheduler_Active
	t.Description = project.Scheduler_Description
}
