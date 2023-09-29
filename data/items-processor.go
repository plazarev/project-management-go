package data

// Processor allows to define some default values for concrete model

type ItemHandlerFunc func(ctx *DBContext, item *Item) error

type IItemProcessor interface {
	PushHandler(handlers ...ItemHandlerFunc)
	Handle(ctx *DBContext, item *Item) error
}

type ItemsProcessor struct {
	handlers []ItemHandlerFunc
}

func NewItemsProcessor() *ItemsProcessor {
	return &ItemsProcessor{
		handlers: make([]ItemHandlerFunc, 0),
	}
}

func (p *ItemsProcessor) PushHandler(handlers ...ItemHandlerFunc) {
	p.handlers = append(p.handlers, handlers...)
}

func (p *ItemsProcessor) Handle(ctx *DBContext, item *Item) error {
	for i := range p.handlers {
		err := p.handlers[i](ctx, item)
		if err != nil {
			return err
		}
	}
	return nil
}
