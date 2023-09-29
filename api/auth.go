package api

import (
	"context"
	"log"
	"net/http"
	"strconv"

	uCtx "project-manager-go/api/context"

	"github.com/go-chi/chi"
)

type AuthAPI struct{}

func (api *AuthAPI) SetAPI(r *chi.Mux) {
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		uid, _ := strconv.Atoi(r.URL.Query().Get("id"))
		device := newDeviceID()
		token, err := createUserToken(uid, device)
		if err != nil {
			log.Println("[token]", err.Error())
		}
		w.Write(token)
	})
}

func (api *AuthAPI) GetPrefix() string {
	return ""
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Remote-Token")
		if token == "" {
			if r.Method == http.MethodGet {
				token = r.URL.Query().Get("token")
			}
		}

		if token != "" {
			id, device, err := verifyUserToken([]byte(token))
			if err != nil {
				log.Println("[token]", err.Error())
			} else {
				r = r.WithContext(context.WithValue(context.WithValue(r.Context(), uCtx.UserIDKey, id), uCtx.DeviceIDKey, device))
			}
		}
		next.ServeHTTP(w, r)
	})
}
