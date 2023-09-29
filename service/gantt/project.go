package gantt

import (
	"project-manager-go/data"
)

type Project struct {
	ID    int
	Index int
}

func (p *Project) PutProject(project data.Project) {
	p.ID = project.ID
	p.Index = project.Index
}

func (p *Project) FillProject(project *data.Project) {
	project.Index = p.Index
}
