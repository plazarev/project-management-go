package data

import (
	"errors"

	"gorm.io/gorm"
)

type ProjectsStore struct {
	Processor IProjectProcessor
}

func NewProjectsStore(processor IProjectProcessor) *ProjectsStore {
	return &ProjectsStore{
		Processor: processor,
	}
}

func (s *ProjectsStore) GetOne(ctx *DBContext, id int, dest IProject) error {
	obj, err := s.getById(ctx, id)
	if err != nil {
		return err
	}
	dest.PutProject(obj)

	return nil
}

func (s *ProjectsStore) GetAll(ctx *DBContext, list IProjectsList) error {
	projects := make([]Project, 0)
	err := ctx.DB.
		Order("`index`").
		Find(&projects).
		Error
	if err != nil {
		return err
	}

	list.PutProjects(projects)

	return nil
}

func (s *ProjectsStore) GetFristId(ctx *DBContext) (int, error) {
	project := Project{}
	err := ctx.DB.
		Order("`index` ASC").
		Take(&project).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return 0, err
	}

	return project.ID, nil
}

func (s *ProjectsStore) GetLastId(ctx *DBContext) (int, error) {
	project := Project{}
	err := ctx.DB.
		Order("`index` DESC").
		Take(&project).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return 0, err
	}

	return project.ID, nil
}

func (s *ProjectsStore) Add(ctx *DBContext, upd IProject) (int, error) {
	obj := Project{}
	upd.FillProject(&obj)

	err := s.Processor.Handle(ctx, &obj)
	if err != nil {
		return 0, err
	}

	err = ctx.DB.Create(&obj).Error

	return obj.ID, err
}

func (s *ProjectsStore) Update(ctx *DBContext, id int, upd IProject) error {
	obj, err := s.getById(ctx, id)
	if err != nil {
		return err
	}
	upd.FillProject(&obj)
	err = ctx.DB.Save(&obj).Error

	return err
}

func (s *ProjectsStore) Delete(ctx *DBContext, id int) ([]int, error) {
	err := ctx.DB.Delete(&Project{}, id).Error
	if err != nil {
		return nil, err
	}

	// collect project items
	items := make([]Item, 0)
	err = ctx.DB.Find(&items, "project_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	ids := make([]int, len(items))
	for i := range items {
		ids[i] = items[i].ID
	}

	err = ctx.DB.Delete(&Item{}, "project_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return ids, err
}

func (s *ProjectsStore) getById(ctx *DBContext, id int) (Project, error) {
	obj := Project{}
	err := ctx.DB.Take(&obj, id).Error
	return obj, err
}

func (s *ProjectsStore) HandleProjectAddOperation(ctx *DBContext, obj *Project) error {
	// the Handler is called before the project is created, default values can be defined here

	project := Project{}
	err := ctx.DB.Order("`index` DESC").Take(&project).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		obj.Index = project.Index + 1
	}

	return nil
}
