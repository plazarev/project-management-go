package publisher

import (
	"project-manager-go/common"
	"project-manager-go/data"
	schedulerStore "project-manager-go/data/scheduler"
	schedulerService "project-manager-go/service/scheduler"

	"github.com/go-chi/chi"
)

type SchedulerEvent struct {
	EventBase
	Event    *schedulerService.Event    `json:"event,omitempty"`
	Calendar *schedulerService.Calendar `json:"calendar,omitempty"`
}

type SchedulerPublisher struct {
	BasePublisher
	store *schedulerStore.SchedulerStore
}

func NewSchedulerPublisher(store *schedulerStore.SchedulerStore, r *chi.Mux, prefix string, routes []string) *SchedulerPublisher {
	t := WidgetScheduler
	api := newRemoteAPI(r, prefix, routes, t)
	return &SchedulerPublisher{
		BasePublisher: BasePublisher{
			widgetType: t,
			api:        api,
		},
		store: store,
	}
}

func (p *SchedulerPublisher) AddItem(ctx PublisherContext, id int) (err error) {
	dbCtx := data.NewCtx(nil)

	ev := schedulerService.Event{}
	err = p.store.Events.GetOne(dbCtx, id, &ev)
	if err != nil {
		return err
	}

	p.api.Events.Publish(
		"events",
		&SchedulerEvent{
			EventBase: EventBase{
				Widget: ctx.FromWidget,
				From:   ctx.DeviceID,
				Type:   "add-event",
			},
			Event: &ev,
		},
	)

	return nil
}

func (p *SchedulerPublisher) UpdateItem(ctx PublisherContext, id int) (err error) {
	dbCtx := data.NewCtx(nil)

	ev := schedulerService.Event{}
	err = p.store.Events.GetOne(dbCtx, id, &ev)
	if err != nil {
		return err
	}

	p.api.Events.Publish(
		"events",
		&SchedulerEvent{
			EventBase: EventBase{
				Widget: ctx.FromWidget,
				From:   ctx.DeviceID,
				Type:   "update-event",
			},
			Event: &ev,
		},
	)

	return nil
}

func (p *SchedulerPublisher) DeleteItem(ctx PublisherContext, id int, children []int) error {
	ids := append([]int{id}, children...)
	for i, id := range ids {
		p.api.Events.Publish(
			"events",
			&SchedulerEvent{
				EventBase: EventBase{
					Widget: ctx.FromWidget,
					From:   ctx.DeviceID,
					Type:   "delete-event",
					Self:   i > 0,
				},
				Event: &schedulerService.Event{ID: common.TID(id)},
			},
		)
	}

	return nil
}

func (p *SchedulerPublisher) AddProject(ctx PublisherContext, id int) (err error) {
	dbCtx := data.NewCtx(nil)

	c := schedulerService.Calendar{}
	err = p.store.Calendars.GetOne(dbCtx, id, &c)
	if err != nil {
		return err
	}

	p.api.Events.Publish(
		"calendars",
		&SchedulerEvent{
			EventBase: EventBase{
				Widget: ctx.FromWidget,
				From:   ctx.DeviceID,
				Type:   "add-calendar",
			},
			Calendar: &c,
		},
	)

	return nil
}

func (p *SchedulerPublisher) UpdateProject(ctx PublisherContext, id int) (err error) {
	dbCtx := data.NewCtx(nil)

	c := schedulerService.Calendar{}
	err = p.store.Calendars.GetOne(dbCtx, id, &c)
	if err != nil {
		return err
	}

	p.api.Events.Publish(
		"calendars",
		&SchedulerEvent{
			EventBase: EventBase{
				Widget: ctx.FromWidget,
				From:   ctx.DeviceID,
				Type:   "update-calendar",
			},
			Calendar: &c,
		},
	)

	return nil
}

func (p *SchedulerPublisher) DeleteProject(ctx PublisherContext, id int, children []int) error {
	p.api.Events.Publish(
		"calendars",
		&SchedulerEvent{
			EventBase: EventBase{
				Widget: ctx.FromWidget,
				From:   ctx.DeviceID,
				Type:   "delete-calendar",
			},
			Calendar: &schedulerService.Calendar{ID: common.TID(id)},
		},
	)

	return nil
}

func (p *SchedulerPublisher) ChangeProject(ctx PublisherContext, id int, projectId int) (err error) {
	dbCtx := data.NewCtx(nil)

	children, err := p.store.Events.GetAllChildrenIDs(dbCtx, projectId, id)
	if err != nil {
		return err
	}

	ids := append([]int{id}, children...)
	for i, eventId := range ids {
		ev := schedulerService.Event{}
		err := p.store.Events.GetOne(dbCtx, eventId, &ev)
		if err != nil {
			return err
		}

		p.api.Events.Publish(
			"events",
			&SchedulerEvent{
				EventBase: EventBase{
					Widget: ctx.FromWidget,
					From:   ctx.DeviceID,
					Type:   "update-event",
					Self:   i > 0,
				},
				Event: &ev,
			},
		)
	}

	return nil
}
