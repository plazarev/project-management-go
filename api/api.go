package api

import (
	"path"

	"github.com/go-chi/chi"
)

type IAPI interface {
	SetAPI(r *chi.Mux)
	GetPrefix() string
}

type BaseAPI struct {
	prefix string
}

type ServerAPI struct {
	r    *chi.Mux
	apis []IAPI
}

func NewServerAPI(r *chi.Mux) *ServerAPI {
	return &ServerAPI{
		r: r,
	}
}

func (s *ServerAPI) AddAPI(apis ...IAPI) {
	for i := range apis {
		s.apis = append(s.apis, apis[i])
		apis[i].SetAPI(s.r)
	}
}

func (b *BaseAPI) GetPrefix() string {
	return b.prefix
}

func (b *BaseAPI) route(s string) string {
	return path.Join(b.GetPrefix() + s)
}
