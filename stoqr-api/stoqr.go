package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/leandroberetta/stoqr/stoqr-api/database"
	"github.com/leandroberetta/stoqr/stoqr-api/repositories"
	"github.com/leandroberetta/stoqr/stoqr-api/server"
	"github.com/leandroberetta/stoqr/stoqr-api/services"
)

func main() {
	log.Println("Starting STOQR")

	database := database.Connect()

	itemRepository := repositories.NewItemRepositorySQL(database)
	itemService := services.NewItemService(itemRepository)

	server := server.NewServer()
	server.Router.Use(mux.CORSMethodMiddleware(server.Router))
	itemService.AddRoutes(server.Router)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	server.Start()

	<-ch

	server.Stop()

	log.Println("Shutdown complete")
}
