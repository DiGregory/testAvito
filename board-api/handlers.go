package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/DiGregory/testAvito/board-api/advert"
	"github.com/go-chi/chi"
	"strings"
	"strconv"
	"database/sql"
)

func (a *apiApp) getAdvertsHandler(w http.ResponseWriter, r *http.Request) {
	var offset = 0
	var err error

	offsetString := r.URL.Query().Get("offset")
	if offsetString != "" {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	sortBy := r.URL.Query().Get("sort_by")
	orderBy := r.URL.Query().Get("order_by")

	ads, err := a.AdvertStorage.GetAdverts(offset,sortBy, orderBy )
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(ads)
	switch err {
	case nil:
		w.Write(jsonResp)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (a *apiApp) getSingleAdvertHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	adID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fieldsRaw := r.URL.Query().Get("fields")
	fields := strings.Split(fieldsRaw, ",")

	ad, err := a.AdvertStorage.GetSingleAdvert(adID, fields)
	switch  {
	case err==sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		return
	case err!=nil:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	adJSON, err := json.Marshal(ad)
	switch err {
	case nil:
		w.Write(adJSON)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (a *apiApp) createAdvertHandler(w http.ResponseWriter, r *http.Request) {
	jsonReq, err := ioutil.ReadAll(r.Body)
	badID := -1
	if err != nil {
		fmt.Println(err)
		sendCreatingResponse(w, badID, http.StatusBadRequest)
		return
	}

	id, err := a.AdvertStorage.CreateAdvert(jsonReq)
	switch {
	case err == advert.InvalidFieldErr:
		sendCreatingResponse(w, badID, http.StatusBadRequest)
		return
	case err != nil:
		sendCreatingResponse(w, badID, http.StatusInternalServerError)
		return
	default:
		sendCreatingResponse(w, *id, http.StatusOK)
		return
	}

}

func sendCreatingResponse(w http.ResponseWriter, id int, statusCode int) {
	var data = map[string]interface{}{}
	data["id"] = id
	data["status_code"] = statusCode
	jsonResp, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(statusCode)
	w.Write(jsonResp)
}
