package blog

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	r.HandleFunc("/blogs", h.handleBlogCreation).Methods("POST") // For creating a blog

	r.HandleFunc("/blogs", h.handleGetAllBlogs).Methods("GET")      // For fetching all blogs
	r.HandleFunc("/blogs/{id}", h.handleGetBlogById).Methods("GET") // For fetching a single blog by ID

	r.HandleFunc("/blogs/{id}", h.handleBlogUpdate).Methods("PATCH") // For updating a blog by ID

	r.HandleFunc("/blogs/soft/{id}", h.handleBlogSoftDeletion).Methods("DELETE")   // Soft delete
	r.HandleFunc("/blogs/delete/{id}", h.handleBlogHardDeletion).Methods("DELETE") // Hard delete

	r.Use(h.AuthMiddleware) // this is to apply the middleware to all the routes under this subrouter
}

func (h *Handler) handleBlogHardDeletion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid blog id: %w", err))
		return
	}
	userId := r.Context().Value(types.UserIDKey).(int64)
	if userId == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// soft delete the blog
	if err := h.store.DeleteBlogPermanentlyById(userId, blogId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete blog: %w", err))
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "permanently deleted the blog"})
}

// this sets the deleted_at field to the current time. and doesn't completely remove the blog from db
func (h *Handler) handleBlogSoftDeletion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid blog id: %w", err))
		return
	}
	userId := r.Context().Value(types.UserIDKey).(int64)
	if userId == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// soft delete the blog
	if err := h.store.SoftDeleteBlogById(userId, blogId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete blog: %w", err))
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "blog soft delete success"})
}

func (h *Handler) handleBlogUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid blog id: %w", err))
		return
	}
	userId := r.Context().Value(types.UserIDKey).(int64)
	if userId == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// get the blog body
	var b types.Blog
	if err := utils.ParseJSON(r, &b); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}
	// validate the blog body
	if err := utils.Validate.Struct(b); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if err := h.store.UpdateBlogById(userId, blogId, b); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update blog: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "blog updated successfully"})
}

func (h *Handler) handleGetAllBlogs(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	userId := r.Context().Value(types.UserIDKey).(int64)
	if userId == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// get all the blogs for the user
	blogs, err := h.store.GetAllBlogs(userId, term)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error getting blogs: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, blogs)
}

func (h *Handler) handleGetBlogById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// get the blog by id
	b, err := h.store.GetBlogById(blogId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, b)
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
