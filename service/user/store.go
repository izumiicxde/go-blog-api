package user

import (
	"github.com/go-playground/validator/v10"
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

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var u types.User
	// get the user with email
	res := s.db.First(&u, "email = ?", email)
	if res.Error != nil {
		return nil, res.Error
	}
	return &u, nil
}
func (s *Store) GetUserById(id int) (*types.User, error) {
	var u types.User

	res := s.db.First(&u, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &u, nil
}

func (s *Store) CreateUser(u types.RegisterUserPayload) error {
	//validate user
	errs := utils.Validate.Struct(u)
	if errs != nil {
		return errs.(validator.ValidationErrors)
	}
	// if no errors create user
	res := s.db.Create(&u)
	return res.Error
}
