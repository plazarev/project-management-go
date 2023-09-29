package data

// Processor allows to define some default values for concrete model

type ProjectHandlerFunc func(ctx *DBContext, obj *Project) error

type IProjectProcessor interface {
	PushHandler(handlers ...ProjectHandlerFunc)
	Handle(ctx *DBContext, obj *Project) error
}

type ProjectsProcessor struct {
	handlers []ProjectHandlerFunc
}

func NewProjectsProcessor() *ProjectsProcessor {
	return &ProjectsProcessor{
		handlers: make([]ProjectHandlerFunc, 0),
	}
}

func (p *ProjectsProcessor) PushHandler(handlers ...ProjectHandlerFunc) {
	p.handlers = append(p.handlers, handlers...)
}

func (p *ProjectsProcessor) Handle(ctx *DBContext, obj *Project) error {
	for i := range p.handlers {
		err := p.handlers[i](ctx, obj)
		if err != nil {
			return err
		}
	}
	return nil
}
