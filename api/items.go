package api

import (
	"net/http"
	"project-manager-go/common"
	"project-manager-go/publisher"
	"project-manager-go/service/items"

	"github.com/go-chi/chi"
)

type ReqItem struct {
	ID   common.TID `json:"id"`
	Item items.Item `json:"item"`
}

type BaseItemsAPI struct {
	BaseAPI
	service *items.BaseItemsService
	pubAll  *publisher.PublisherAPI
}

func NewBaseItemsAPI(
	service *items.BaseItemsService,
	pubAll *publisher.PublisherAPI,
	prefix string,
) *BaseItemsAPI {
	return &BaseItemsAPI{
		BaseAPI: BaseAPI{
			prefix: prefix,
		},
		service: service,
		pubAll:  pubAll,
	}
}

func (api *BaseItemsAPI) SetAPI(r *chi.Mux) {
	r.Get(api.route("/items"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		items, err := api.service.Items.GetAll(userCtx, nil)

		respond(w, &items, err)
	})

	r.Post(api.route("/items"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		item := ReqItem{}
		err = parseForm(w, r, &item)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id, err := api.service.Items.Add(userCtx, nil, item.Item)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:   userCtx.ID,
				DeviceID: userCtx.DeviceID,
			}
			api.pubAll.PublishAddItem(pubCtx, id)
		}
	})

	r.Put(api.route("/items/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		item := ReqItem{}
		err = parseForm(w, r, &item)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")

		err = api.service.Items.Update(userCtx, nil, id, item.Item)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:   userCtx.ID,
				DeviceID: userCtx.DeviceID,
			}
			api.pubAll.PublishUpdateItem(pubCtx, id)
		}
	})

	r.Delete(api.route("/items/{id}"), func(w http.ResponseWriter, r *http.Request) {
		userCtx, err := parseUserContext(r)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		id := parseNumberParam(r, "id")
		children, err := api.service.Items.Delete(userCtx, nil, id)

		if respond(w, &ResponseID{id}, err) {
			pubCtx := publisher.PublisherContext{
				UserID:   userCtx.ID,
				DeviceID: userCtx.DeviceID,
			}
			api.pubAll.PublishDeleteItem(pubCtx, id, children)
		}
	})
}
