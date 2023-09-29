package api

import (
	"net/http"
	"project-manager-go/common"
	kanbanStore "project-manager-go/data/kanban"
	"project-manager-go/publisher"
	"project-manager-go/service/kanban"

	"github.com/go-chi/chi"
)

type ReqCard struct {
	ID   common.TID  `json:"id"`
	Card kanban.Card `json:"card"`
}

type ReqRow struct {
	ID  common.TID `json:"id"`
	Row kanban.Row `json:"row"`
}

type ReqColumn struct {
	ID     common.TID `json:"id"`
	Column struct {
		Label string `json:"label"`
	} `json:"column"`
}

type KanbanAPI struct {
	BaseAPI
	service   *kanban.KanbanService
	pubKanban *publisher.KanbanPublisher
	pubAll    *publisher.PublisherAPI
}

func NewKanbanAPI(
	service *kanban.KanbanService,
	pub *publisher.KanbanPublisher,
	pubAll *publisher.PublisherAPI,
	prefix string,
) *KanbanAPI {
	return &KanbanAPI{
		BaseAPI: BaseAPI{
			prefix: prefix,
		},
		service:   service,
		pubKanban: pub,
		pubAll:    pubAll,
	}
}

func (api *KanbanAPI) SetAPI(r *chi.Mux) {
	r.Get(api.route("/cards"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		cards, err := api.service.Cards.GetAll(userCtx, nil)

		respond(w, &cards, err)
	})

	r.Post(api.route("/cards"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		card := ReqCard{}
		err = parseForm(w, r, &card)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id, err := api.service.Cards.Add(userCtx, nil, card.Card)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetKanban,
			}
			api.pubAll.PublishAddItem(pubCtx, id)
		}
	})

	r.Put(api.route("/cards/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		card := ReqCard{}
		err = parseForm(w, r, &card)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")

		err = api.service.Cards.Update(userCtx, nil, id, card.Card)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetKanban,
			}
			api.pubAll.PublishUpdateItem(pubCtx, id)
		}
	})

	r.Put(api.route("/cards/{id}/move"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		params := kanban.MoveCards{}
		err = parseForm(w, r, &params)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")
		batch := []kanban.MoveParams{params.MoveParams}
		if params.ID == 0 {
			batch = params.Batch
		}
		resp, err := api.service.Cards.Move(userCtx, nil, id, batch)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetKanban,
			}

			for _, params := range batch {
				// publish move operation to the kanban clients
				api.pubKanban.MoveCard(pubCtx, int(params.ID), int(params.Before), int(params.RowID), int(params.ColumnID))
			}
			for _, id := range resp.UpdatedProject {
				api.pubAll.PublishChangeProject(pubCtx, id, int(params.RowID))
			}
		}
	})

	r.Delete(api.route("/cards/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")
		children, err := api.service.Cards.Delete(userCtx, nil, id)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetKanban,
			}
			api.pubAll.PublishDeleteItem(pubCtx, id, children)
		}
	})

	r.Get(api.route("/rows"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		rows, err := api.service.Rows.GetAll(userCtx, nil)

		respond(w, &rows, err)
	})

	r.Post(api.route("/rows"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		row := kanban.Row{}
		err = parseForm(w, r, &row)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id, err := api.service.Rows.Add(userCtx, nil, row)

		if respond(w, ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetKanban,
			}
			api.pubAll.PublishAddProject(pubCtx, id)
		}
	})

	r.Put(api.route("/rows/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		row := ReqRow{}
		err = parseForm(w, r, &row)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")
		err = api.service.Rows.Update(userCtx, nil, id, row.Row)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetKanban,
			}
			api.pubAll.PublishUpdateProject(pubCtx, id)
		}
	})

	r.Put(api.route("/rows/{id}/move"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		params := kanban.MoveParams{}
		err = parseForm(w, r, &params)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")

		err = api.service.Rows.Move(userCtx, nil, id, int(params.Before))

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetKanban,
			}
			api.pubKanban.MoveRow(pubCtx, id, int(params.Before))
		}
	})

	r.Delete(api.route("/rows/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")
		children, err := api.service.Rows.Delete(userCtx, nil, id)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetKanban,
			}
			api.pubAll.PublishDeleteProject(pubCtx, id, children)
		}
	})

	r.Get(api.route("/columns"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		cols, err := api.service.Cols.GetAll(userCtx, nil)

		respond(w, &cols, err)
	})

	r.Post(api.route("/columns"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		form := ReqColumn{}
		err = parseForm(w, r, &form)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		column := kanbanStore.KanbanColumn{
			Label: form.Column.Label,
		}
		id, err := api.service.Cols.Add(userCtx, nil, column)

		if respond(w, ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetKanban,
			}
			api.pubKanban.AddColumn(pubCtx, id)
		}
	})

	r.Put(api.route("/columns/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		form := ReqColumn{}
		err = parseForm(w, r, &form)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")
		column := kanbanStore.KanbanColumn{
			ID:    id,
			Label: form.Column.Label,
		}
		err = api.service.Cols.Update(userCtx, nil, id, column)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetKanban,
			}
			api.pubKanban.UpdateColumn(pubCtx, id)
		}
	})

	r.Put(api.route("/columns/{id}/move"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		params := kanban.MoveParams{}
		err = parseForm(w, r, &params)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")

		err = api.service.Cols.Move(userCtx, nil, id, int(params.Before))

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetKanban,
			}
			api.pubKanban.MoveColumn(pubCtx, id, int(params.Before))
		}
	})

	r.Delete(api.route("/columns/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")
		children, err := api.service.Cols.Delete(userCtx, nil, id)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:     userCtx.ID,
				DeviceID:   userCtx.DeviceID,
				FromWidget: publisher.WidgetKanban,
			}
			api.pubAll.PublishDeleteItem(pubCtx, 0, children)
			api.pubKanban.DeleteColumn(pubCtx, id, children)
		}
	})

	r.Get(api.route("/users"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		users, err := api.service.Users.GetAll(userCtx, nil)

		respond(w, &users, err)
	})

}
