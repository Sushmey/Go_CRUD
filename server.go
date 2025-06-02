package main

import (
	"github.com/gorilla/mux"
	"log"
	"log/slog"
	"net/http"
)

type Server interface{
	Start() error //Starts the server and returns errors if they pop up
	routes() // Define server routes
}

type MuxServer struct{
	gorilla *mux.Router
	Client
}


//Constructor to create new server
func NewServer(db Client) Server{

	server := &MuxServer{
		mux.NewRouter(),
		db,
	}
	server.routes()
	return server
}

func (s *MuxServer) Start() error{
	slog.Info("Serving at port 8080")
	log.Fatal(http.ListenAndServe(":8080",s.gorilla))
	return nil
}

