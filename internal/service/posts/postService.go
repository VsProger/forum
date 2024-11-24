package service

import (
	"errors"
	"fmt"

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
	GetPostsByUserId(user_id int) ([]models.Post, error)
	AddReaction(reaction models.Reaction) error
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
	// If no categories are provided, add a default one
	if len(post.Categories) == 0 {
		post.Categories = append(post.Categories, models.Category{Name: "Other"})
	}

	// Iterate through the categories provided in the post
	for i, category := range post.Categories {
		// Call GetCategoryByName to retrieve the category details
		categories, err := s.postRepo.GetCategoryByName(category.Name)
		if err != nil {
			return err
		}

		// If no category is found, return an error
		if len(categories) == 0 {
			return fmt.Errorf("category %s not found", category.Name)
		}

		// Assuming you want to use the first matching category
		post.Categories[i] = *categories[0]
	}

	// Now, save the post with its categories
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

func (s *postService) GetPostsByUserId(user_id int) ([]models.Post, error) {
	posts, err := s.postRepo.GetAllPostsByUserId(user_id)
	if err != nil {
		return posts, err
	}
	return posts, nil
}

func (s *postService) AddReaction(reaction models.Reaction) error {
	switch {
	case reaction.PostID != 0 && reaction.CommentID == 0:
		if err := s.postRepo.AddReactionToPost(reaction); err != nil {
			return fmt.Errorf("Error adding or updating reaction: %v", err)
		}
	case reaction.CommentID != 0 && reaction.PostID != 0:
		if err := s.postRepo.AddReactionToComment(reaction); err != nil {
			fmt.Println("HEREE")
			return err
		}
	default:
		return fmt.Errorf("specify either PostId or CommentId, not both")
	}
	return nil
}
