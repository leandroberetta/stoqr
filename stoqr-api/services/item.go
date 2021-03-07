package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/leandroberetta/stoqr/stoqr-api/models"
	"github.com/leandroberetta/stoqr/stoqr-api/repositories"
	"github.com/leandroberetta/stoqr/stoqr-api/server"
)

// ItemService contains the business logic of items
type ItemService struct {
	Repository repositories.ItemRepository
}

// CreateItem is the api method for create an item
func (svc *ItemService) CreateItem(w http.ResponseWriter, r *http.Request) {
	item := models.Item{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	err = svc.Repository.CreateItem(&item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// ReadItem is the api method for get an item
func (svc *ItemService) ReadItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["itemId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	item, err := svc.Repository.ReadItem(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// ReadItems is the api method to get items, optionally filtered by name
func (svc *ItemService) ReadItems(w http.ResponseWriter, r *http.Request) {
	filter := r.FormValue("filter")
	items, err := svc.Repository.ReadItems(filter)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// UpdateItem is the api method to update an item
func (svc *ItemService) UpdateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["itemId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	item := models.Item{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	err = svc.Repository.UpdateItem(id, item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteItem is the api method for delete an item
func (svc *ItemService) DeleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["itemId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = svc.Repository.DeleteItem(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// WithdrawItem is the api method for withdraw an item
func (svc *ItemService) WithdrawItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["itemId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	item, err := svc.Repository.ReadItem(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if item.Actual != 0 {
		item.Actual = item.Actual - 1
	}
	err = svc.Repository.UpdateItem(id, item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// AddRoutes configures the items routes into a given router
func (svc *ItemService) AddRoutes(r *mux.Router) {
	r.HandleFunc("/api/items", server.Options).Methods(http.MethodOptions)
	r.HandleFunc("/api/items", svc.CreateItem).Methods(http.MethodPost)
	r.HandleFunc("/api/items", svc.ReadItems).Methods(http.MethodGet)
	r.HandleFunc("/api/items", svc.ReadItems).Methods(http.MethodGet).Queries("filter", "{filter}")
	r.HandleFunc("/api/items/{itemId}", server.Options).Methods(http.MethodOptions)
	r.HandleFunc("/api/items/{itemId}", svc.ReadItem).Methods((http.MethodGet))
	r.HandleFunc("/api/items/{itemId}", svc.DeleteItem).Methods(http.MethodDelete)
	r.HandleFunc("/api/items/{itemId}", svc.UpdateItem).Methods(http.MethodPut)
	r.HandleFunc("/api/items/withdraw/{itemId}", server.Options).Methods(http.MethodOptions)
	r.HandleFunc("/api/items/withdraw/{itemId}", svc.WithdrawItem).Methods(http.MethodGet)
}

// NewItemService creates a new item service
func NewItemService(repository repositories.ItemRepository) *ItemService {
	return &ItemService{Repository: repository}
}
