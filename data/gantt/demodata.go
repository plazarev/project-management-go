package gantt

import (
	"path"
	"project-manager-go/data"
)

type GanttDataProvider struct {
	folder string
}

func NewGanttDataProvider(folder string) *GanttDataProvider {
	return &GanttDataProvider{folder}
}

func (p GanttDataProvider) Up(ctx *data.DBContext) {
	data.InitializeDemodata[GanttLink](ctx, path.Join(p.folder, "gantt_links.json"))
}

func (p GanttDataProvider) Down(ctx *data.DBContext) {
	err := ctx.DB.Delete(&GanttLink{}, "1 = 1").Error
	if err != nil {
		panic(err)
	}
}
