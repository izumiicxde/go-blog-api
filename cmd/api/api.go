package api

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/izumii.cxde/blog-api/service/blog"
	"github.com/izumii.cxde/blog-api/service/user"
	"gorm.io/gorm"
)

type APIServer struct {
	addr string
	db   *gorm.DB
}

func NewAPIServer(addr string, db *gorm.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

// signature
func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	blogStore := blog.NewStore(s.db)
	blogHandler := blog.NewHandler(blogStore, userStore)
	blogHandler.RegisterRoutes(subrouter)

	slog.Info("Listening on: ", slog.String("addr", s.addr))
	return http.ListenAndServe(s.addr, router)
}
