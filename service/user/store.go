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

func (s *Store) SendVerificationCode(email, otp, username string) error {
	if err := mail.SendMail(otp, email, username); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}
	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var u types.User
	// get the user with email
	res := s.db.First(&u, "email = ?", email)
	if res.Error != nil {
		return nil, res.Error
	}
	return &u, nil
}

func (s *Store) GetUserById(id int64) (*types.User, error) {
	var u types.User

	res := s.db.First(&u, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &u, nil
}

// signature method to create a new user\n
// @params: u(RegisterUserPayload) user info
func (s *Store) CreateUser(u types.RegisterUserPayload, otp string) error {
	//validate user
	errs := utils.Validate.Struct(u)
	if errs != nil {
		return errs.(validator.ValidationErrors)
	}

	// if no errors create user
	// hash the user password
	h, err := auth.HashPassword(u.Password)
	if err != nil {
		return err
	}
	// passing u of RegisterUserPayload is causing error with gorm
	user := types.User{
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		Password:      h, // hashed with bcrypt
		AvatarUrl:     u.AvatarUrl,
		Otp:           otp,
		OtpExpiration: time.Now().Add(time.Minute * 5),
		Verified:      false,
	}

	return s.db.Create(&user).Error
}

func (s *Store) UpdateUserById(id int64, u types.User) error {
	if errs := utils.Validate.Struct(u); errs != nil {
		return errs.(validator.ValidationErrors)
	}
	res := s.db.Model(&types.User{}).Where("id = ?", id).Updates(u)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *Store) DeleteUserById(id int64) error {
	return s.db.Delete(&types.User{}, id).Error
}
