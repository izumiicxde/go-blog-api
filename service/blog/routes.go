package blog

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/izumii.cxde/blog-api/service/auth"
	"github.com/izumii.cxde/blog-api/types"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := auth.ParseJWTRequest(r)
		if err != nil || userId == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	r := router.PathPrefix("/").Subrouter()
	r.HandleFunc("/blogs/create", h.handleBlogCreation).Methods("POST")
}

func (h *Handler) handleBlogCreation(w http.ResponseWriter, r *http.Request) {

}
