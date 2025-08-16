package api

import (
	"digits_say/storage"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Config struct {
	ListenAddr string
	DB         storage.DBConfig
}

type Server struct {
	*Config
	Started time.Time

	Storage *storage.Storage
}

func NewServer(config Config) (*Server, error) {
	storage := storage.Storage{
		DBConfig: config.DB,
	}

	err := storage.ConnectToSurreal()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SurrealDB: %w", err)
	}

	return &Server{
		Config:  &config,
		Storage: &storage,
	}, nil
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /user/", s.GetUserByID)
	mux.HandleFunc("POST /user", s.RegisterNewUser)
	mux.HandleFunc("UPDATE /user", s.UpdateUser)
	mux.HandleFunc("PATCH /user", s.UpdateUser)
	mux.HandleFunc("GET /conscience", s.GetConscienceText)
	mux.HandleFunc("GET /subscribers", s.GetListOfSubscribers)
	mux.HandleFunc("GET /commonday", s.GetCommonDayText)

	mux.HandleFunc("GET /health", s.Health)
	
	server := http.Server{
		Addr: s.Config.ListenAddr,
		Handler: mux,
	}

	s.Started = time.Now().UTC()
	slog.Info("api is running", "address", s.Config.ListenAddr)
	err:=server.ListenAndServe()
	if err!=nil{
		s.Storage.Close()
		slog.Error("API stoped", "error", err)
	}

}