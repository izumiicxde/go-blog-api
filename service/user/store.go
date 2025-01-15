package user

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/izumii.cxde/blog-api/mail"
	"github.com/izumii.cxde/blog-api/service/auth"
	"github.com/izumii.cxde/blog-api/types"
	"github.com/izumii.cxde/blog-api/utils"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// SendVerificationCode sends a verification code to the user's email address using smtp server with my gmail
func (s *Store) SendVerificationCode(email, otp, username string) error {
	// can just return the error. But I added a custom error message for better clarification
	if err := mail.SendMail(otp, email, username); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}
	return nil
}

// GetUserByEmail gets a user by their email address
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var u types.User
	res := s.db.First(&u, "email = ?", email)
	return &u, res.Error
}

// GetUserById gets a user by their id
func (s *Store) GetUserById(id int64) (*types.User, error) {
	var u types.User
	res := s.db.First(&u, id)
	return &u, res.Error
}

/*
signature method to create a new user
@params: u(RegisterUserPayload) user info
*/
func (s *Store) CreateUser(u types.RegisterUserPayload, otp string) error {
	//validate user
	errs := utils.Validate.Struct(u)
	if errs != nil {
		return errs.(validator.ValidationErrors)
	}

	// if no errors create user
	// hash the user password
	hashedPassword, err := auth.HashPassword(u.Password)
	if err != nil {
		return err
	}
	// passing u of RegisterUserPayload is causing error with gorm
	user := types.User{
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		Password:      hashedPassword, // hashed with bcrypt
		AvatarUrl:     u.AvatarUrl,
		Otp:           otp,
		OtpExpiration: time.Now().Add(time.Minute * 5),
		Verified:      false,
	}

	return s.db.Create(&user).Error
}

func (s *Store) UpdateUserById(id int64, u types.User) error {
	// Validate the struct before updating
	if errs := utils.Validate.Struct(u); errs != nil {
		return errs.(validator.ValidationErrors)
	}

	// Use `Select` to include all fields explicitly
	u.Otp = "0"
	res := s.db.Model(&types.User{}).
		Where("id = ?", id).
		Updates(u)

	if res.Error != nil {
		return res.Error
	}
	if res.Error.Error() == "record not found" {
		return fmt.Errorf("blog not found")
	}
	return nil
}

func (s *Store) DeleteUserById(id int64) error {
	return s.db.Delete(&types.User{}, id).Error
}
