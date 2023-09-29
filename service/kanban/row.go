package kanban

import (
	"project-manager-go/common"
	"project-manager-go/data"
	"project-manager-go/data/kanban/kanban_props"
)

type Row struct {
	ID    common.TID `json:"id"`
	Label string     `json:"label"`
	Index int        `json:"index"`
}

func (c *Row) FillProject(proj *data.Project) {
	proj.Label = c.Label
	proj.KanbanProjectProps = kanban_props.KanbanProjectProps{
		Kanban_RowIndex: c.Index,
	}
}

func (t *Row) PutProject(project data.Project) {
	t.ID = common.TID(project.ID)
	t.Label = project.Label
	t.Index = project.Kanban_RowIndex
}
