package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/leandroberetta/stoqr/stoqr-api/models"
	"github.com/leandroberetta/stoqr/stoqr-api/repositories"
	"github.com/leandroberetta/stoqr/stoqr-api/services"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func corsOriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		next.ServeHTTP(w, r)
	})
}

func options(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
}

func checkPostgresConfig() error {
	if host := os.Getenv("STOQR_API_DB_HOST"); host == "" {
		return errors.New(("db: postgres host is empty"))
	}
	if port := os.Getenv("STOQR_API_DB_PORT"); port == "" {
		return errors.New(("db: postgres port is empty"))
	}
	if user := os.Getenv("STOQR_API_DB_USER"); user == "" {
		return errors.New(("db: postgres user is empty"))
	}
	if password := os.Getenv("STOQR_API_DB_PASSWORD"); password == "" {
		return errors.New(("db: postgres password is empty"))
	}
	if name := os.Getenv("STOQR_API_DB_NAME"); name == "" {
		return errors.New(("db: postgres db name is empty"))
	}
	return nil
}

func generatePostgresConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("STOQR_API_DB_HOST"),
		os.Getenv("STOQR_API_DB_USER"),
		os.Getenv("STOQR_API_DB_PASSWORD"),
		os.Getenv("STOQR_API_DB_NAME"),
		os.Getenv("STOQR_API_DB_PORT"))
}

func connectDB() *gorm.DB {
	db := &gorm.DB{}
	if err := checkPostgresConfig(); err != nil {
		log.Println(err)
		log.Println("Falling back to SQLite")
		db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db, err = gorm.Open(postgres.Open(generatePostgresConnectionString()), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	}
	db.AutoMigrate(&models.Item{})
	return db
}

func main() {
	log.Println("Starting STOQR")

	db := connectDB()
	itemRepository := repositories.NewItemRepositorySQL(db)
	itemService := services.NewItemService(itemRepository)

	r := mux.NewRouter()
	r.HandleFunc("/api/items", options).Methods(http.MethodOptions)
	r.HandleFunc("/api/items", itemService.CreateItem).Methods(http.MethodPost)
	r.HandleFunc("/api/items", itemService.ReadItems).Methods(http.MethodGet)
	r.HandleFunc("/api/items", itemService.ReadItems).Methods(http.MethodGet).Queries("filter", "{filter}")
	r.HandleFunc("/api/items/{itemId}", options).Methods(http.MethodOptions)
	r.HandleFunc("/api/items/{itemId}", itemService.ReadItem).Methods((http.MethodGet))
	r.HandleFunc("/api/items/{itemId}", itemService.DeleteItem).Methods(http.MethodDelete)
	r.HandleFunc("/api/items/{itemId}", itemService.UpdateItem).Methods(http.MethodPut)
	r.HandleFunc("/api/items/withdraw/{itemId}", options).Methods(http.MethodOptions)
	r.HandleFunc("/api/items/withdraw/{itemId}", itemService.WithdrawItem).Methods(http.MethodGet)
	r.Use(corsOriginMiddleware)
	r.Use(mux.CORSMethodMiddleware(r))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Println("Serving at :8080")
	go srv.ListenAndServe()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutdown complete")
	os.Exit(0)
}
