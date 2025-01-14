package blog

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

func (s *Store) CreateBlog(b types.Blog) error {
	errs := utils.Validate.Struct(b)
	if errs != nil {
		errs = errs.(validator.ValidationErrors)
		return errs
	}

	// Find or create tags based on the provided tag names
	var tags []types.Tag
	for _, tagName := range b.Tags {
		var tag types.Tag
		if err := s.db.Where("name = ?", tagName).First(&tag).Error; err != nil {
			// If the tag doesn't exist, create it
			tag = types.Tag{Name: tagName.Name}
			if err := s.db.Create(&tag).Error; err != nil {
				return err
			}
		}
		tags = append(tags, tag)
	}

	// Assign the tags to the blog (many-to-many relationship)
	b.Tags = tags

	return s.db.Create(&b).Error
}
