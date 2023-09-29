package gantt

import (
	uCtx "project-manager-go/api/context"
	"project-manager-go/data"
	"project-manager-go/data/gantt"
	"project-manager-go/service"
)

type tasks struct {
	tree  *service.TreeService
	store *gantt.GanttStore
}

func (s *tasks) GetAll(userCtx uCtx.UserContext, dbCtx *data.DBContext) (arr []Task, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	list := data.NewItemsList[Task]()
	err = s.store.Tasks.GetAll(dbCtx, list)

	return list.GetArray(), err
}

func (s *tasks) Add(userCtx uCtx.UserContext, dbCtx *data.DBContext, task Task) (id int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	// the Gantt widget does not provide the ProjectID property, so it must be calculated
	var projectId int
	if task.ParentID == 0 {
		// cannot identify the project ID, so push to the end of the last project
		projectId, err = s.store.ProjectsStore.GetLastId(dbCtx)
		if err != nil {
			return 0, err
		}
	} else {
		// should find the parent and take the project ID from it
		parentNode := service.Node{}
		err := s.store.Tasks.GetOne(dbCtx, int(task.ParentID), &parentNode)
		if err != nil {
			return 0, err
		}
		projectId = parentNode.ProjectID
	}

	// get the maximum index in the branch
	maxBranchIndex, err := s.store.Tasks.MaxBranchIndex(dbCtx, projectId, int(task.ParentID))
	if err != nil {
		return 0, err
	}

	task.Index = maxBranchIndex + 1

	id, err = s.store.Tasks.Add(dbCtx, &task)
	if err != nil {
		return 0, err
	}

	// Task model does not provide implicit ProjectID property, so need to update it manually
	err = s.store.Tasks.UpdateFields(dbCtx, id, map[string]any{
		"project_id": projectId,
	})

	return id, err
}

func (s *tasks) Update(userCtx uCtx.UserContext, dbCtx *data.DBContext, id int, upd Task) (err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	err = s.store.Tasks.Update(dbCtx, id, &upd)

	return err
}

func (s *tasks) Delete(userCtx uCtx.UserContext, dbCtx *data.DBContext, id int) (children []int, err error) {
	dbCtx = data.NewTCtx(dbCtx)
	defer func() { err = dbCtx.End(err) }()

	children, err = s.store.Tasks.DeleteCascade(dbCtx, id)

	return children, err
}

func (s *tasks) ToGlobalIndex(dbCtx *data.DBContext, project, parent, index int) (int, error) {
	if parent != 0 {
		return index, nil
	}

	projects := data.NewProjectsList[Project]()
	err := s.store.ProjectsStore.GetAll(dbCtx, projects)
	if err != nil {
		return 0, err
	}

	globalIndex := 0
	for _, p := range projects.GetArray() {
		if p.ID == project {
			break
		}

		maxIndex, err := s.store.Tasks.MaxBranchIndex(dbCtx, p.ID, 0)
		if err != nil {
			return 0, err
		}

		globalIndex += maxIndex + 1
	}

	return globalIndex + index, nil
}
