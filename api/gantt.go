package api

import (
	"net/http"
	"project-manager-go/publisher"
	"project-manager-go/service/gantt"

	"github.com/go-chi/chi"
)

type GanttAPI struct {
	BaseAPI
	service  *gantt.GanttService
	pubGantt *publisher.GanttPublisher
	pubAll   *publisher.PublisherAPI
}

func NewGanttAPI(service *gantt.GanttService, pub *publisher.GanttPublisher, pubAll *publisher.PublisherAPI, prefix string) *GanttAPI {
	return &GanttAPI{
		BaseAPI: BaseAPI{
			prefix: prefix,
		},
		service:  service,
		pubGantt: pub,
		pubAll:   pubAll,
	}
}

func (api *GanttAPI) SetAPI(r *chi.Mux) {
	r.Get(api.route("/tasks"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		tasks, err := api.service.Tasks.GetAll(userCtx, nil)

		respond(w, &tasks, err)
	})

	r.Post(api.route("/task"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		task := gantt.Task{}
		err = parseForm(w, r, &task)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id, err := api.service.Tasks.Add(userCtx, nil, task)

		if respond(w, &ResponseTID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetGantt,
			}
			api.pubAll.PublishAddItem(pubCtx, id)
		}
	})

	r.Put(api.route("/task/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		task := gantt.Task{}
		err = parseForm(w, r, &task)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")

		err = api.service.Tasks.Update(userCtx, nil, id, task)

		if respond(w, &ResponseTID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetGantt,
			}
			api.pubAll.PublishUpdateItem(pubCtx, id)
		}
	})

	r.Delete(api.route("/task/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")
		children, err := api.service.Tasks.Delete(userCtx, nil, id)

		if respond(w, &ResponseTID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetGantt,
			}
			api.pubAll.PublishDeleteItem(pubCtx, id, children)
		}
	})

	r.Get(api.route("/links"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		links, err := api.service.Links.GetAll(userCtx, nil)

		respond(w, &links, err)
	})

	r.Post(api.route("/link"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		link := gantt.Link{}
		err = parseForm(w, r, &link)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id, err := api.service.Links.Add(userCtx, nil, link)

		if respond(w, ResponseTID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetGantt,
			}
			api.pubGantt.AddLink(pubCtx, id)
		}
	})

	r.Delete(api.route("/link/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")
		err = api.service.Links.Delete(userCtx, nil, id)

		if respond(w, &ResponseTID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetGantt,
			}
			api.pubGantt.DeleteLink(pubCtx, id)
		}
	})
}
