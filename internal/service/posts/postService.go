package posts

import (
	"errors"

	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/internal/repository/posts"
	"github.com/VsProger/snippetbox/pkg"
)

type PostService interface {
	CreatePost(post models.Post) error
	CreateCategory(name string) error
	GetPostByID(id int) (*models.Post, error)
	GetPosts() ([]models.Post, error)
	CreateComment(comment models.Comment) error
}

type postService struct {
	postRepo posts.Posts
}

func NewPostService(postRepo posts.Posts) PostService {
	return &postService{
		postRepo: postRepo,
	}
}

func (s *postService) GetPosts() ([]models.Post, error) {
	return s.postRepo.GetPosts()
}

func (s *postService) CreatePost(post models.Post) error {
	// if len(post.Categories) == 0 {
	// 	post.Categories[0] = models.Category{Name: "Other"}
	// }
	// for i, category := range post.Categories {
	// 	t, err := s.postRepo.GetCategoryByName(category.Name)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	post.Categories[i] = *t
	// }

	return s.postRepo.CreatePost(post)
}

func (s *postService) GetPostByID(id int) (*models.Post, error) {
	return s.postRepo.GetPostByID(id)
}

func (s *postService) CreateCategory(name string) error {
	if name == "" {
		return errors.New("Name should be be empty!")
	}
	_, err := s.postRepo.GetCategoryByName(name)
	if err == nil {
		return errors.New("category already exists")
	}
	err = s.postRepo.CreateCategory(name)
	if err != nil {
		return err
	}

	return nil
}

func (s *postService) CreateComment(comment models.Comment) error {
	if err := pkg.ValidateComment(comment); err != nil {
		return err
	}
	return s.postRepo.CreateComment(comment)
}
