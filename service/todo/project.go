package todo

import (
	"project-manager-go/common"
	"project-manager-go/data"
)

type Project struct {
	ID    common.TID `json:"id"`
	Label string     `json:"label"`
}

func (p *Project) PutProject(project data.Project) {
	p.ID = common.TID(project.ID)
	p.Label = project.Label
}

func (p *Project) FillProject(proj *data.Project) {
	proj.Label = p.Label
}
