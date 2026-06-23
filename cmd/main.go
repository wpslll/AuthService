package main

import (
	"AuthService/internal/handlers"
	"AuthService/internal/repository"
	"AuthService/internal/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
)

type Handler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Auth(w http.ResponseWriter, r *http.Request)
	Validate(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	h Handler
}

func NewRouter(h Handler) Router {
	return Router{
		h: h,
	}
}

func (s *Router) CreateRouter() mux.Router {
	r := mux.NewRouter()
	r.Path("/auth").Methods("POST").HandlerFunc(s.h.Auth)
	r.Path("/register").Methods("POST").HandlerFunc(s.h.Register)
	r.Path("/validate").Methods("Get").HandlerFunc(s.h.Validate)
	return *r
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	dbUrl := os.Getenv("DB_URL")
	fmt.Println(dbUrl)
	m, err := migrate.New("file://internal/migrations", dbUrl)
	if err != nil {
		fmt.Println("Failed to create migration", err)
		return
	}
	if err := m.Up(); err != nil {
		fmt.Println("Failed to up mig", err)
	}
	defer m.Close()
	db, err := repository.NewDB(ctx, dbUrl)
	service := service.NewService(db)
	handler := handlers.NewHandler(&service)
	r := NewRouter(handler)
	router := r.CreateRouter()
	
	server := &http.Server{
		Addr:         ":" + os.Getenv("SERVER_PORT"),
		Handler:      &router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Internal server error: ", err)
			cancel()
		}
	}()
	<-ctx.Done()
	shutdownCtx, shtdc := context.WithTimeout(context.Background(), 10*time.Second)
	defer shtdc()
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Println("Failed to shutdown the server", err)
	}
}
