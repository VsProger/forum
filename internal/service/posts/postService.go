package posts

import (
	"errors"

	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/internal/repository/posts"
)

type PostService interface {
	CreatePost(post models.Post) error
	CreateCategory(name string) error
}

type postService struct {
	postRepo posts.Posts
}

func NewPostService(postRepo posts.Posts) PostService {
	return &postService{
		postRepo: postRepo,
	}
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
