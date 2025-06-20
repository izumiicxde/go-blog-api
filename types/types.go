package types

import (
	"time"

	"gorm.io/gorm"
)

// suggested by the compiler
type ContextKey string

const UserIDKey ContextKey = "userId"

// === === POST === ===
type BlogStore interface {
	CreateBlog(b Blog) error
	GetAllBlogs() (*[]Blog, error)
	GetBlogById(id int64) (*Blog, error)
	GetAllBlogsByUserId(userId int64, term string) (*[]Blog, error)
	UpdateBlogById(userId, id int64, b Blog) error
	SoftDeleteBlogById(userId, id int64) error
	DeleteBlogPermanentlyById(userId, id int64) error
}

type Tag struct {
	gorm.Model
	Name string `json:"name" gorm:"uniqueIndex" validate:"required,min=1,max=50"`
}

type BlogTag struct {
	BlogID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`
}

type Blog struct {
	gorm.Model
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"required,min=3,max=500"`
	Content     string `json:"content" validate:"required,min=3,max=3000"`
	Category    string `json:"category" validate:"required,min=3,max=255"`
	// tags are separated by commas for now. can validate to new table with many to many relation
	Tags   []Tag `json:"tags" validate:"required" gorm:"many2many:blog_tags;"`
	UserId uint  `json:"user_id" validate:"-"` // Foreign key reference
	// User   User  `gorm:"foreignKey:UserId"`           // Establish relationship
}

type VerificationPayload struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required"`
}

// === === USER  === ===
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int64) (*User, error)
	CreateUser(u RegisterUserPayload, otp string) error
	UpdateUserById(id int64, u User) error
	DeleteUserById(id int64) error
	SendVerificationCode(email, otp, username string) error
}

type User struct {
	gorm.Model
	FirstName     string    `json:"first_name" validate:"required,min=3,max=30"`
	LastName      string    `json:"last_name" validate:"required,max=30"`
	Email         string    `json:"email" gorm:"uniqueIndex" validate:"required,email"`
	Password      string    `json:"password" validate:"required"`
	AvatarUrl     string    `json:"avatar_url" validate:"required"`
	Blogs         []Blog    `gorm:"foreignKey:UserId"` // One-to-many relationship
	Otp           string    `json:"-" validate:"-"`
	OtpExpiration time.Time `json:"-" validate:"-"`
	Verified      bool      `json:"-" validate:"-" gorm:"default:false"`
}

type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=30"`
	LastName  string `json:"last_name" validate:"required,max=30"`
	Email     string `json:"email"  validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	AvatarUrl string `json:"avatar_url" validate:"required"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
