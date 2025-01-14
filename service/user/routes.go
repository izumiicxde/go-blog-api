package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/izumii.cxde/blog-api/service/auth"
	"github.com/izumii.cxde/blog-api/types"
	"github.com/izumii.cxde/blog-api/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

	var u types.LoginUserPayload
	if err := utils.ParseJSON(r, &u); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	errs := utils.Validate.Struct(u)
	if errs != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body %s", errs.Error()))
		return
	}

	// get user by email
	user, err := h.store.GetUserByEmail(u.Email)
	if err != nil || user == nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user not found: %w", err))
		return
	}
	// check if the password is correct
	if !auth.CompareHashPassword(user.Password, u.Password) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}
	// generate the token
	t, err := auth.GenerateJWTToken(*user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("auth failed: %w", err))
	}

	// if all is good set the cookie
	c := http.Cookie{
		Name:    "token",
		Value:   t,
		Expires: time.Now().Add(time.Hour * 24 * 7),
	}
	http.SetCookie(w, &c)

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "login successful"})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var u types.RegisterUserPayload
	if err := utils.ParseJSON(r, &u); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// get user by email
	_, err := h.store.GetUserByEmail(u.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with this email already exists"))
		return
	}
	if err = h.store.CreateUser(u); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user created successfully"})
}
