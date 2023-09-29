package api

import (
	"net/http"
	"project-manager-go/publisher"
	"project-manager-go/service/scheduler"

	"github.com/go-chi/chi"
)

type SchedulerAPI struct {
	BaseAPI
	service      *scheduler.SchedulerService
	pubScheduler *publisher.SchedulerPublisher
	pubAll       *publisher.PublisherAPI
}

func NewSchedulerAPI(
	service *scheduler.SchedulerService,
	pub *publisher.SchedulerPublisher,
	pubAll *publisher.PublisherAPI,
	prefix string,
) *SchedulerAPI {
	return &SchedulerAPI{
		BaseAPI: BaseAPI{
			prefix: prefix,
		},
		service:      service,
		pubScheduler: pub,
		pubAll:       pubAll,
	}
}

func (api *SchedulerAPI) SetAPI(r *chi.Mux) {
	r.Get(api.route("/events"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		events, err := api.service.Events.GetAll(userCtx, nil)

		respond(w, &events, err)
	})

	r.Post(api.route("/events"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		event := scheduler.Event{}
		err = parseForm(w, r, &event)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id, err := api.service.Events.Add(userCtx, nil, event)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetScheduler,
			}
			api.pubAll.PublishAddItem(pubCtx, id)
		}
	})

	r.Put(api.route("/events/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		event := scheduler.Event{}
		err = parseForm(w, r, &event)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")

		updProject, err := api.service.Events.Update(userCtx, nil, id, event)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetScheduler,
			}

			api.pubAll.PublishUpdateItem(pubCtx, id)

			if updProject {
				api.pubAll.PublishChangeProject(pubCtx, id, int(event.CalendarID))
			}
		}
	})

	r.Delete(api.route("/events/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")
		children, err := api.service.Events.Delete(userCtx, nil, id)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetScheduler,
			}
			api.pubAll.PublishDeleteItem(pubCtx, id, children)
		}
	})

	r.Get(api.route("/calendars"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		calendars, err := api.service.Calendars.GetAll(userCtx, nil)

		respond(w, &calendars, err)
	})

	r.Post(api.route("/calendars"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		calendar := scheduler.Calendar{}
		err = parseForm(w, r, &calendar)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id, err := api.service.Calendars.Add(userCtx, nil, calendar)

		if respond(w, ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetScheduler,
			}
			api.pubAll.PublishAddProject(pubCtx, id)
		}
	})

	r.Put(api.route("/calendars/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		calendar := scheduler.Calendar{}
		err = parseForm(w, r, &calendar)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")
		err = api.service.Calendars.Update(userCtx, nil, id, calendar)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetScheduler,
			}
			api.pubAll.PublishUpdateProject(pubCtx, id)
		}
	})

	r.Delete(api.route("/calendars/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")
		children, err := api.service.Calendars.Delete(userCtx, nil, id)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetScheduler,
			}
			api.pubAll.PublishDeleteProject(pubCtx, id, children)
		}
	})
}
