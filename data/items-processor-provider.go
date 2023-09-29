package data

type ItemsProcessorProvider struct {
	ev map[string]IItemProcessor
}

func NewItemsProcessorProvider() *ItemsProcessorProvider {
	return &ItemsProcessorProvider{
		ev: make(map[string]IItemProcessor),
	}
}

func (p *ItemsProcessorProvider) PushProcessor(name string, processor IItemProcessor) {
	p.ev[name] = processor
}

func (p *ItemsProcessorProvider) Handle(name string, ctx *DBContext, obj *Item) error {
	processor := p.ev[name]
	err := processor.Handle(ctx, obj)
	return err
}
