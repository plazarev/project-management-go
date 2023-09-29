package kanban

import (
	"errors"
	"project-manager-go/data"

	"gorm.io/gorm"
)

type rows struct {
	*data.ProjectsStore
}

func (s *rows) GetAll(ctx *data.DBContext, list data.IProjectsList) error {
	projects := make([]data.Project, 0)
	err := ctx.DB.Order("kanban_row_index").Find(&projects).Error
	if err != nil {
		return err
	}

	list.PutProjects(projects)

	return nil
}

func (s *rows) GetMaxIndex(ctx *data.DBContext) (int, error) {
	row := data.Project{}

	err := ctx.DB.Order("kanban_row_index DESC").Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	return row.Kanban_RowIndex + 1, err
}

func (s *rows) GetByIndex(ctx *data.DBContext, index int, dest data.IProject) error {
	row := data.Project{}
	err := ctx.DB.Order("kanban_row_index DESC").Take(&row, "kanban_row_index = ?", index).Error
	dest.PutProject(row)
	return err
}
