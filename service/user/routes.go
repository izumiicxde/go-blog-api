package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/izumii.cxde/blog-api/mail"
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
	router.HandleFunc("/verify", h.handleVerification).Methods("POST")
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
		Name:     "token",
		Value:    t,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HttpOnly: true,
		Secure:   true,
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
	user, err := h.store.GetUserByEmail(u.Email)
	if err != nil {
		if !(err.Error() == "record not found") {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}
	if err == nil || user != nil {
		if user.Verified {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user already exists"))
		}
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	otp := auth.GenerateOTP()
	if err = h.store.CreateUser(u, otp); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err))
		return
	}

	ok, err := mail.SendMail(otp, u.Email, fmt.Sprintf("%s %s", u.FirstName, u.LastName))
	if !ok || err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to send verification email: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user created successfully, please verify your email"})
}

func (h *Handler) handleVerification(w http.ResponseWriter, r *http.Request) {
	var VerificationPayload struct {
		Email string `json:"email" validate:"required,email"`
		Otp   string `json:"otp" validate:"required"`
	}
	if err := utils.ParseJSON(r, &VerificationPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(VerificationPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body %s", err.Error()))
		return
	}
	u, err := h.store.GetUserByEmail(VerificationPayload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if u.Verified {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user already verified"))
		return
	}
	if !auth.ValidateOTP(VerificationPayload.Otp, *u) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid otp"))
		return
	}
	u.Verified = true
	u.Otp = ""
	if err := h.store.UpdateUserById(int64(u.ID), *u); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "user verified successfully"})
}
