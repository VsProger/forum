package posts

import (
	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/internal/repository/posts"
)

type PostService interface {
	CreatePost(post models.Post) error
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
