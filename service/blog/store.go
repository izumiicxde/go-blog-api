package blog

import (
	"fmt"

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

/*
UpdateBlogById updates a blog by its id
@params:
userId - the id of the user
id - the id of the blog
b - the blog to update

@returns:
error - if there was an error
*/
func (s *Store) UpdateBlogById(userId, id int64, b types.Blog) error {
	res := s.db.Model(&types.Blog{}).Where("user_id = ? AND id = ? AND deleted_at is NULL", userId, id).Updates(b)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("no blog found")
	}
	return nil
}

func (s *Store) CreateBlog(b types.Blog) error {
	// Validate the blog object
	errs := utils.Validate.Struct(b)
	if errs != nil {
		errs = errs.(validator.ValidationErrors)
		return errs
	}

	// Find or create tags based on the provided tag names
	var tags []types.Tag
	for _, tagName := range b.Tags {
		var tag types.Tag
		if err := s.db.FirstOrCreate(&tag, types.Tag{Name: tagName.Name}).Error; err != nil {
			return err
		}
		tags = append(tags, tag)
	}
	// Assign the tags to the blog (many-to-many relationship)
	b.Tags = tags

	// Create the blog
	return s.db.Create(&b).Error
}

/*
GetBlogById returns a blog by its id
If the blog doesn't exist, it returns nil with an error
@params:

	id - the id of the blog

@returns:

	blog - the blog with the given id
	error - if there was an error
*/
func (s *Store) GetBlogById(id int64) (*types.Blog, error) {
	var b types.Blog
	// The preload is used to eager load the tags relationship
	if err := s.db.Preload("Tags").First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

/*
GetAllBlogs returns all the blogs for a given user
@params:

	userId - the id of the user

@returns:

	blogs - a slice of blogs
*/
func (s *Store) GetAllBlogs(userId int64, term string) (*[]types.Blog, error) {
	var blogs []types.Blog
	// Start building the query with the user_id filter and deleted_at check
	query := s.db.Preload("Tags").Where("user_id = ? AND deleted_at is NULL", userId)
	// If a search term is provided, filter the results based on the term
	if term != "" {
		query = query.Where("title LIKE ? OR description LIKE ? OR category LIKE ?", "%"+term+"%", "%"+term+"%", "%"+term+"%")
	}

	// Execute the query
	if err := query.Find(&blogs).Error; err != nil {
		return nil, err
	}

	// If no blogs are found, return an error
	if len(blogs) == 0 {
		return nil, fmt.Errorf("no blogs found")
	}

	// Return the found blogs
	return &blogs, nil
}

/*
SoftDeleteBlogById soft deletes a blog by its id
@params:
userId - the id of the user
id - the id of the blog
@returns:
error - if there was an error
*/
func (s *Store) SoftDeleteBlogById(userId, id int64) error {
	res := s.db.Where("user_id = ? AND id = ? ", userId, id).Delete(&types.Blog{})
	if res.RowsAffected == 0 {
		return fmt.Errorf("no blog found")
	}
	// if the error is not blog not found then just return the error
	return res.Error
}

/*
DeleteBlogPermanentlyById deletes a blog permanently by its id
@params:
userId - the id of the user
id - the id of the blog

@returns:
error - if there was an error
*/
func (s *Store) DeleteBlogPermanentlyById(userId, id int64) error {

	res := s.db.Unscoped().Where("user_id = ? AND id = ?", userId, id).Delete(&types.Blog{})
	if res.RowsAffected == 0 {
		return fmt.Errorf("no blog found")
	}
	return res.Error
}
