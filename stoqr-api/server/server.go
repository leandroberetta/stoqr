package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server represents the server structure
type Server struct {
	Server *http.Server
	Router *mux.Router
}

// NewServer creates a Server instance
func NewServer() *Server {
	server := &Server{}
	server.Router = mux.NewRouter()
	server.Server = &http.Server{Addr: ":8080", Handler: server.Router}
	return server
}

// Start starts the server
func (server *Server) Start() {
	log.Println("Serving at :8080")
	go server.Server.ListenAndServe()
}

// Stop stops the server
func (server *Server) Stop() {
	log.Println("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	server.Server.Shutdown(ctx)
}

// Options is a handler for the OPTIONS method used for CORS
func Options(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}
