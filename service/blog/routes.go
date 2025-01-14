package blog

import (
	"context"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/izumii.cxde/blog-api/service/auth"
	"github.com/izumii.cxde/blog-api/types"
	"github.com/izumii.cxde/blog-api/utils"
)

type Handler struct {
	store     types.BlogStore
	userStore types.UserStore
}

func NewHandler(store types.BlogStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := auth.ParseJWTRequest(r)
		if err != nil || userId == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), types.UserIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	r := router.PathPrefix("/").Subrouter()
	r.HandleFunc("/blogs/create", h.handleBlogCreation).Methods("POST")

	r.Use(h.AuthMiddleware) // this is to apply the middleware to all the routes under this subrouter
}

func (h *Handler) handleBlogCreation(w http.ResponseWriter, r *http.Request) {
	// first authorize the user
	userId := r.Context().Value(types.UserIDKey).(int64)
	log.Println("USERID", userId)
	if userId == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	_, err := h.userStore.GetUserById(userId)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}

	// get the blog body
	var b types.Blog
	if err := utils.ParseJSON(r, &b); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// validate blog body
	if err := utils.Validate.Struct(b); err != nil {
		err = err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	b.UserId = uint(userId)
	// create the blog
	err = h.store.CreateBlog(b)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "blog created successfully"})
}
