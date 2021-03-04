package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Server *http.Server
	Router *mux.Router
}

func NewServer() *Server {
	server := &Server{}
	server.Router = mux.NewRouter()
	server.Server = &http.Server{Addr: ":8080", Handler: server.Router}
	return server
}

func (server *Server) Start() {
	log.Println("Serving at :8080")
	go server.Server.ListenAndServe()
}

func (server *Server) Stop() {
	log.Println("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	server.Server.Shutdown(ctx)
}

func Options(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}
