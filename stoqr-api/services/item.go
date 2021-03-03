package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/leandroberetta/stoqr/stoqr-api/models"
	"github.com/leandroberetta/stoqr/stoqr-api/repositories"
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
		log.Fatal(err)
	}
	defer r.Body.Close()
	err = svc.Repository.CreateItem(&item)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// ReadItem is the api method for get an item
func (svc *ItemService) ReadItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["itemId"])
	if err != nil {
		log.Fatal(err)
	}
	item, err := svc.Repository.ReadItem(id)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// ReadItems is the api method to get items, optionally filtered by name
func (svc *ItemService) ReadItems(w http.ResponseWriter, r *http.Request) {
	filter := r.FormValue("filter")
	items, err := svc.Repository.ReadItems(filter)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// UpdateItem is the api method to update an item
func (svc *ItemService) UpdateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["itemId"])
	if err != nil {
		log.Fatal(err)
	}
	item := models.Item{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&item)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	err = svc.Repository.UpdateItem(id, item)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteItem is the api method for delete an item
func (svc *ItemService) DeleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["itemId"])
	if err != nil {
		log.Fatal(err)
	}
	err = svc.Repository.DeleteItem(id)
	if err != nil {
		log.Fatal(err)
	}
}

// WithdrawItem is the api method for withdraw an item
func (svc *ItemService) WithdrawItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["itemId"])
	if err != nil {
		log.Fatal(err)
	}
	item, err := svc.Repository.ReadItem(id)
	if err != nil {
		log.Fatal(err)
	}
	if item.Actual != 0 {
		item.Actual = item.Actual - 1
	}
	err = svc.Repository.UpdateItem(id, item)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// NewItemService creates a new item service
func NewItemService(repository repositories.ItemRepository) *ItemService {
	return &ItemService{Repository: repository}
}
