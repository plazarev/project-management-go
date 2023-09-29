package api

import (
	"encoding/json"
	"errors"
	"net/http"
	uCtx "project-manager-go/api/context"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

var dID = time.Now().Unix()

type ResponseID struct {
	ID int `json:"id"`
}

type ResponseTID struct {
	ID int `json:"tid"`
}

func respond(w http.ResponseWriter, data any, err error) bool {
	if err != nil {
		respondWithError(w, err.Error())
		return false
	}
	respondWithJSON(w, data)
	return true
}

func respondWithJSON(w http.ResponseWriter, data any) {
	w.Header().Add("Content-Type", "application/json")
	bytes, _ := json.Marshal(&data)
	w.Write(bytes)
}

func respondWithError(w http.ResponseWriter, message string) {
	respondWithErrorCode(w, message, http.StatusInternalServerError)
}

func respondWithErrorCode(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}

func parseForm(w http.ResponseWriter, r *http.Request, dest any) error {
	body := http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(body)
	err := dec.Decode(&dest)
	return err
}

func parseNumberParam(r *http.Request, key string) int {
	value := chi.URLParam(r, key)
	num, _ := strconv.Atoi(value)
	return num
}

func parseUserContext(r *http.Request) (uCtx.UserContext, error) {
	vUser := r.Context().Value(uCtx.UserIDKey)
	id, ok := vUser.(int)
	if vUser != nil && !ok {
		return uCtx.UserContext{}, errors.New("can not parse user ID")
	}

	vDevice := r.Context().Value(uCtx.DeviceIDKey)
	device, ok := vDevice.(int)
	if vDevice != nil && !ok {
		return uCtx.UserContext{}, errors.New("can not parse device ID")
	}

	ctx := uCtx.UserContext{
		ID:       id,
		DeviceID: device,
	}

	return ctx, nil
}

func newDeviceID() int64 {
	dID += 1
	return dID
}
