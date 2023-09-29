package publisher

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"

	uCtx "project-manager-go/api/context"

	remote "github.com/mkozhukh/go-remote"
)

type PublisherAPI struct {
	r    *chi.Mux
	pubs []IPublisherAPI
}

func NewPublisherAPI(r *chi.Mux) *PublisherAPI {
	return &PublisherAPI{
		r:    r,
		pubs: make([]IPublisherAPI, 0),
	}
}

func (p *PublisherAPI) AddAPI(pubs ...IPublisherAPI) {
	p.pubs = append(p.pubs, pubs...)
}

func (p *PublisherAPI) PublishAddItem(ctx PublisherContext, id int) error {
	for i := range p.pubs {
		err := p.pubs[i].AddItem(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PublisherAPI) PublishUpdateItem(ctx PublisherContext, id int) error {
	for i := range p.pubs {
		err := p.pubs[i].UpdateItem(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PublisherAPI) PublishDeleteItem(ctx PublisherContext, id int, children []int) error {
	for i := range p.pubs {
		err := p.pubs[i].DeleteItem(ctx, id, children)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PublisherAPI) PublishAddProject(ctx PublisherContext, id int) error {
	for i := range p.pubs {
		err := p.pubs[i].AddProject(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PublisherAPI) PublishUpdateProject(ctx PublisherContext, id int) error {
	for i := range p.pubs {
		err := p.pubs[i].UpdateProject(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PublisherAPI) PublishDeleteProject(ctx PublisherContext, id int, children []int) error {
	for i := range p.pubs {
		err := p.pubs[i].DeleteProject(ctx, id, children)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PublisherAPI) PublishChangeProject(ctx PublisherContext, id int, projectId int) error {
	for i := range p.pubs {
		m, ok := p.pubs[i].(IMovePublisher)
		if !ok {
			continue
		}
		err := m.ChangeProject(ctx, id, projectId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PublisherAPI) PublishMove(ctx PublisherContext, id int, params MoveParams) error {
	for i := range p.pubs {
		m, ok := p.pubs[i].(ITreePublisher)
		if !ok {
			continue
		}
		err := m.Move(ctx, id, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func newRemoteAPI(r *chi.Mux, path string, routes []string, widget WidgetType) *remote.Server {
	if remote.MaxSocketMessageSize < 32000 {
		remote.MaxSocketMessageSize = 32000
	}

	api := remote.NewServer(&remote.ServerConfig{
		WebSocket: true,
	})

	api.Connect = func(r *http.Request) (context.Context, error) {
		id, _ := r.Context().Value(uCtx.UserIDKey).(int)
		if id == 0 {
			return nil, errors.New("access denied")
		}
		device, _ := r.Context().Value(uCtx.DeviceIDKey).(int)
		if device == 0 {
			return nil, errors.New("access denied")
		}

		return context.WithValue(
			context.WithValue(r.Context(), remote.UserValue, id),
			remote.ConnectionValue, device), nil
	}

	for _, route := range routes {
		api.Events.AddGuard(route, func(m *remote.Message, c *remote.Client) bool {
			e, ok := m.Content.(IPublisherEvent)
			if !ok {
				return false
			}

			sameDevice := e.GetFrom() == c.ConnID
			sameWidget := e.GetWidget() == widget

			return e.IsSelf() || !sameDevice || !sameWidget
		})
	}

	r.Get(path, api.ServeHTTP)
	r.Post(path, api.ServeHTTP)

	return api
}
