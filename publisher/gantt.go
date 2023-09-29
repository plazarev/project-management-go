package publisher

import (
	"project-manager-go/data"
	ganttStore "project-manager-go/data/gantt"
	"project-manager-go/service"
	ganttService "project-manager-go/service/gantt"

	"github.com/go-chi/chi"
)

type GanttEvent struct {
	EventBase
	Data any `json:"data"`
}

type GanttPublisher struct {
	BasePublisher
	service *ganttService.GanttService
	store   *ganttStore.GanttStore
}

func NewGanttPublisher(service *ganttService.GanttService, store *ganttStore.GanttStore, r *chi.Mux, prefix string, routes []string) *GanttPublisher {
	t := WidgetGantt
	api := newRemoteAPI(r, prefix, routes, t)
	return &GanttPublisher{
		BasePublisher: BasePublisher{
			widgetType: t,
			api:        api,
		},
		service: service,
		store:   store,
	}
}

func (p *GanttPublisher) AddItem(ctx PublisherContext, id int) (err error) {
	dbCtx := data.NewCtx(nil)

	task := ganttService.Task{}
	err = p.store.Tasks.GetOne(dbCtx, id, &task)
	if err != nil {
		return err
	}

	p.api.Events.Publish(
		"tasks",
		&GanttEvent{
			EventBase: EventBase{
				From:   ctx.DeviceID,
				Type:   "add-task",
				Widget: ctx.FromWidget,
			},
			Data: &task,
		},
	)

	p.Move(ctx, id, MoveParams{
		ParentID: int(task.ParentID),
	})

	return nil
}

func (p *GanttPublisher) UpdateItem(ctx PublisherContext, id int) error {
	dbCtx := data.NewCtx(nil)

	task := ganttService.Task{}
	err := p.store.Tasks.GetOne(dbCtx, id, &task)
	if err != nil {
		return err
	}

	p.api.Events.Publish(
		"tasks",
		&GanttEvent{
			EventBase: EventBase{
				From:   ctx.DeviceID,
				Type:   "update-task",
				Widget: ctx.FromWidget,
			},
			Data: &task,
		},
	)

	return nil
}

func (p *GanttPublisher) DeleteItem(ctx PublisherContext, id int, children []int) error {
	p.api.Events.Publish(
		"tasks",
		&GanttEvent{
			EventBase: EventBase{
				From:   ctx.DeviceID,
				Type:   "delete-task",
				Widget: ctx.FromWidget,
			},
			Data: &ganttService.Task{ID: id},
		},
	)

	return nil
}

func (p *GanttPublisher) AddProject(ctx PublisherContext, id int) error {
	// not used
	return nil
}

func (p *GanttPublisher) UpdateProject(ctx PublisherContext, id int) error {
	// not used
	return nil
}

func (p *GanttPublisher) DeleteProject(ctx PublisherContext, id int, children []int) error {
	for i := range children {
		p.api.Events.Publish(
			"tasks",
			&GanttEvent{
				EventBase: EventBase{
					From:   ctx.DeviceID,
					Type:   "delete-task",
					Widget: ctx.FromWidget,
				},
				Data: &ganttService.Task{ID: children[i]},
			},
		)
	}
	return nil
}

func (p *GanttPublisher) ChangeProject(ctx PublisherContext, id int, projectId int) error {
	dbCtx := data.NewCtx(nil)

	node := service.Node{}
	err := p.store.Tasks.GetOne(dbCtx, id, &node)
	if err != nil {
		return err
	}

	index, err := p.service.Tasks.ToGlobalIndex(dbCtx, node.ProjectID, node.ParentID, node.Index)
	if err != nil {
		return err
	}

	params := MoveParams{
		ID:       id,
		Index:    index,
		ParentID: 0,
	}

	p.api.Events.Publish(
		"tasks",
		&TodoEvent{
			EventBase: EventBase{
				Widget: ctx.FromWidget,
				From:   ctx.DeviceID,
				Type:   "move-task",
			},
			Data: &params,
			ID:   id,
		},
	)

	return nil
}

// Custom methods:

func (p *GanttPublisher) Move(ctx PublisherContext, id int, params MoveParams) error {
	dbCtx := data.NewCtx(nil)

	node := service.Node{}
	err := p.store.Tasks.GetOne(dbCtx, id, &node)
	if err != nil {
		return err
	}

	index, err := p.service.Tasks.ToGlobalIndex(dbCtx, node.ProjectID, node.ParentID, node.Index)
	if err != nil {
		return err
	}

	params = MoveParams{
		ID:       id,
		Index:    index,
		ParentID: params.ParentID,
	}

	p.api.Events.Publish(
		"tasks",
		&TodoEvent{
			EventBase: EventBase{
				Widget: ctx.FromWidget,
				From:   ctx.DeviceID,
				Type:   "move-task",
			},
			Data: &params,
			ID:   id,
		},
	)

	return nil
}

func (p *GanttPublisher) AddLink(ctx PublisherContext, id int) error {
	dbCtx := data.NewCtx(nil)

	link, err := p.store.Links.GetOne(dbCtx, id)
	if err != nil {
		return err
	}

	p.api.Events.Publish(
		"links",
		&GanttEvent{
			EventBase: EventBase{
				From:   ctx.DeviceID,
				Type:   "add-link",
				Widget: ctx.FromWidget,
			},
			Data: &link,
		},
	)

	return nil
}

func (p *GanttPublisher) UpdateLink(ctx PublisherContext, id int) error {
	dbCtx := data.NewCtx(nil)

	link, err := p.store.Links.GetOne(dbCtx, id)
	if err != nil {
		return err
	}

	p.api.Events.Publish(
		"links",
		&GanttEvent{
			EventBase: EventBase{
				From:   ctx.DeviceID,
				Type:   "update-link",
				Widget: ctx.FromWidget,
			},
			Data: &link,
		},
	)

	return nil
}

func (p *GanttPublisher) DeleteLink(ctx PublisherContext, id int) error {
	p.api.Events.Publish(
		"links",
		&GanttEvent{
			EventBase: EventBase{
				From:   ctx.DeviceID,
				Type:   "delete-link",
				Widget: ctx.FromWidget,
			},
			Data: ganttStore.GanttLink{
				ID: id,
			},
		},
	)

	return nil
}
