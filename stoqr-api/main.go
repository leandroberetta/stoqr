package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Item is the model of the item object
type Item struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Desired int    `json:"desired"`
	Actual  int    `json:"actual"`
}

// ItemRepository interface define the methods to persist items
type ItemRepository interface {
	CreateItem(item Item) error
	ReadItem(id int) (Item, error)
	UpdateItem(id int, item Item) error
	DeleteItem(id int) error
	ReadItems(filter string) ([]Item, error)
}

// ItemRepositorySQL persist items into a SQL database
type ItemRepositorySQL struct {
	*gorm.DB
}

// CreateItem persists an item into a database
func (db *ItemRepositorySQL) CreateItem(item Item) error {
	result := db.Create(&item)
	return result.Error
}

// ReadItem gets an item from a database
func (db *ItemRepositorySQL) ReadItem(id int) (Item, error) {
	var item Item
	result := db.First(&item, id)
	return item, result.Error
}

// UpdateItem updates an item and persists it into a database
func (db *ItemRepositorySQL) UpdateItem(id int, updatedItem Item) error {
	var item Item
	result := db.Find(&item, id)
	if result.Error != nil {
		return result.Error
	}
	item.Name = updatedItem.Name
	item.Desired = updatedItem.Desired
	item.Actual = updatedItem.Actual
	db.Save(&item)
	return nil
}

// DeleteItem removes an item from a database
func (db *ItemRepositorySQL) DeleteItem(id int) error {
	result := db.Delete(&Item{}, id)
	return result.Error
}

// ReadItems gets items from a database optionally filtering by name
func (db *ItemRepositorySQL) ReadItems(filter string) ([]Item, error) {
	var items []Item
	if filter != "" {
		db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", filter)).Find(&items)
	} else {
		db.Find(&items)
	}
	return items, nil
}

// NewItemRepositorySQL returns a new ItemRepositorySQL instance
func NewItemRepositorySQL(db *gorm.DB) ItemRepository {
	return &ItemRepositorySQL{db}
}

// ItemService contains the business logic of items
type ItemService struct {
	Repository ItemRepository
}

func (svc *ItemService) createItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	item := Item{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&item)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	err = svc.Repository.CreateItem(item)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusCreated)
}

func (svc *ItemService) readItem(w http.ResponseWriter, r *http.Request) {
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

func (svc *ItemService) readItems(w http.ResponseWriter, r *http.Request) {
	filter := r.FormValue("filter")
	items, err := svc.Repository.ReadItems(filter)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (svc *ItemService) updateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["itemId"])
	if err != nil {
		log.Fatal(err)
	}
	item := Item{}
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

func (svc *ItemService) deleteItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
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

// NewItemService creates a new item service
func NewItemService(repository ItemRepository) *ItemService {
	return &ItemService{Repository: repository}
}

func corsAllowed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		next.ServeHTTP(w, r)
	})
}

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("STOQR_API_DB_HOST"),
		os.Getenv("STOQR_API_DB_USER"),
		os.Getenv("STOQR_API_DB_PASSWORD"),
		os.Getenv("STOQR_API_DB_NAME"),
		os.Getenv("STOQR_API_DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Item{})
	itemRepository := NewItemRepositorySQL(db)
	itemService := NewItemService(itemRepository)
	r := mux.NewRouter()
	r.HandleFunc("/api/items/{itemId}", itemService.readItem).Methods((http.MethodGet))
	r.HandleFunc("/api/items", itemService.createItem).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/items", itemService.readItems).Methods(http.MethodGet)
	r.HandleFunc("/api/items", itemService.readItems).Methods(http.MethodGet).Queries("filter", "{filter}")
	r.HandleFunc("/api/items/{itemId}", itemService.deleteItem).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/items/{itemId}", itemService.updateItem).Methods(http.MethodPut)
	r.Use(corsAllowed)
	log.Fatal(http.ListenAndServe(":8080", r))
}
