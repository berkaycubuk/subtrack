package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/berkaycubuk/subtrack/internal/services"
)

type Server struct {
	httpServer *http.Server
	subSvc     *services.SubscriptionService
	sessions   *sessionStore
	username   string
	password   string
}

func NewServer(subSvc *services.SubscriptionService, username, password string) *Server {
	srv := &Server{
		subSvc:   subSvc,
		sessions: newSessionStore(),
		username: username,
		password: password,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /login", srv.handleLoginForm)
	mux.HandleFunc("POST /login", srv.handleLogin)
	mux.HandleFunc("GET /logout", srv.handleLogout)
	mux.HandleFunc("GET /{$}", srv.requireAuth(srv.handleDashboard))
	mux.HandleFunc("GET /add", srv.requireAuth(srv.handleAddForm))
	mux.HandleFunc("POST /add", srv.requireAuth(srv.handleAdd))
	mux.HandleFunc("GET /edit/{id}", srv.requireAuth(srv.handleEditForm))
	mux.HandleFunc("POST /edit/{id}", srv.requireAuth(srv.handleEdit))
	mux.HandleFunc("GET /delete/{id}", srv.requireAuth(srv.handleDeleteForm))
	mux.HandleFunc("POST /delete/{id}", srv.requireAuth(srv.handleDelete))

	srv.httpServer = &http.Server{
		Handler: mux,
	}

	return srv
}

func (s *Server) Start(addr string) error {
	s.httpServer.Addr = addr
	log.Printf("Web server starting on %s", addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
