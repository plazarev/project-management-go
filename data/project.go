package data

type IProject interface {
	FillProject(*Project)
	PutProject(Project)
}

type IWrapperProject[T any] interface {
	FillProject(*Project)
	PutProject(Project)
	*T
}

type IProjectsList interface {
	PutProjects([]Project)
}

type ProjectsList[T any, PT IWrapperProject[T]] struct {
	arr []T
}

func NewProjectsList[T any, PT IWrapperProject[T]]() *ProjectsList[T, PT] {
	return &ProjectsList[T, PT]{
		arr: make([]T, 0),
	}
}

func (l *ProjectsList[T, PT]) GetArray() []T {
	return l.arr
}

func (l *ProjectsList[T, PT]) PutProjects(projects []Project) {
	l.arr = make([]T, len(projects))
	for i := range projects {
		// create not nil instance
		project := PT(new(T))
		// fill with the properties
		project.PutProject(projects[i])
		l.arr[i] = *project
	}
}
