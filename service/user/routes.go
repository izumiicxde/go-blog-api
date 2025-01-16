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

	router.HandleFunc("/verify", h.handleVerification).Methods("POST")
	router.HandleFunc("/get-verification-code", h.handleSendVerificationCode).Methods("GET")
}

// HandleLogin handles the login request
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// get, parse and validate the payload
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
	// get, parse and validate the payload
	var u types.RegisterUserPayload
	if err := utils.ParseJSON(r, &u); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	//get user by email. check if user exists
	// if err is nil then the user exists.
	user, err := h.store.GetUserByEmail(u.Email)
	if err != nil {
		if !(err.Error() == "record not found") {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}
	// user cannot register with same email if any users exists with same email
	if (err == nil) || (user != nil) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user already exists"))
		return
	}

	// if user doesn't exist create user and send them the verification code.
	otp := auth.GenerateOTP()
	if err = h.store.CreateUser(u, otp); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err))
		return
	}

	// send the email to the user with the verification code.
	if err := h.store.SendVerificationCode(u.Email, otp, fmt.Sprintf("%s %s", u.FirstName, u.LastName)); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user created successfully, please verify your email"})
}

// this is to validate the verification code given by the user with the one sent
func (h *Handler) handleVerification(w http.ResponseWriter, r *http.Request) {
	// get , parse and validate the user payload
	var verificationPayload types.VerificationPayload
	if err := utils.ParseJSON(r, &verificationPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(verificationPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body %s", err.Error()))
		return
	}

	// get the user by email as provided by the front-end
	u, err := h.store.GetUserByEmail(verificationPayload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// check if user is already verified
	if u.Verified {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user already verified"))
		return
	}
	// validate the otp with the user provided one
	if !auth.ValidateOTP(verificationPayload.Otp, *u) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid otp"))
		return
	}
	// if the otp is correct. Then set verified to true.
	u.Verified = true
	u.Otp = ""
	u.OtpExpiration = time.Now() // to say the otp has expired. just-in-case
	// update the user with new field values
	if err := h.store.UpdateUserById(int64(u.ID), *u); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "user verified successfully"})
}

func (h *Handler) handleSendVerificationCode(w http.ResponseWriter, r *http.Request) {
	// get, parse and validate the email. to send the verification code
	var p struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := utils.ParseJSON(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// get the user.
	u, err := h.store.GetUserByEmail(p.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// check if user is already verified. Then we don't send the code again
	if u.Verified {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user already verified"))
		return
	}

	// Check if OTP expiration time has passed. This is to stop from too many requests
	if time.Now().Before(u.OtpExpiration) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("please wait until %s before requesting a new OTP", u.OtpExpiration.Format(time.RFC1123)))
		return
	}

	// generate the otp
	otp := auth.GenerateOTP()
	// send the email with otp to user.
	if err := h.store.SendVerificationCode(u.Email, otp, fmt.Sprintf("%s %s", u.FirstName, u.LastName)); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// update the user in db with the otp for future use
	u.Otp = otp
	u.OtpExpiration = time.Now().Add(time.Minute * 5) // 5 min expiration
	if err := h.store.UpdateUserById(int64(u.ID), *u); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "verification code sent successfully", "verification_code": otp})
}
