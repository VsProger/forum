package service

import (
	"errors"
	"fmt"
	"log"
	"time"

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
	GetNotificationsByUserID(user_id int) ([]models.Notification, error)
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
	// Валидация комментария
	if err := pkg.ValidateComment(comment); err != nil {
		return err
	}

	// Создание комментария
	if err := s.postRepo.CreateComment(comment); err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	// Получение поста для отправки уведомления
	post, err := s.postRepo.GetPostByID(comment.PostID)
	if err != nil {
		return fmt.Errorf("comment created, but failed to retrieve post for notification: %w", err)
	}

	// Формирование уведомления
	notification := models.Notification{
		UserID:    post.AuthorID,
		PostID:    comment.PostID,
		CommentID: comment.ID,
		Type:      "new_comment",
		Message:   fmt.Sprintf("Your post '%s' received a new comment: %s", post.Title, comment.Text),
		CreatedAt: time.Now(),
		IsRead:    false,
	}

	// Сохранение уведомления в БД
	if err := s.postRepo.CreateNotification(notification); err != nil {
		return fmt.Errorf("comment created, but failed to save notification: %w", err)
	}

	// Асинхронная отправка уведомления
	go func() {
		if err := s.postRepo.NotifyUser(post.AuthorID, notification.Message); err != nil {
			// Логирование ошибки
			log.Printf("failed to send notification: %v", err)
		}
	}()

	return nil
}

func (s *postService) GetPostsByUserId(user_id int) ([]models.Post, error) {
	posts, err := s.postRepo.GetAllPostsByUserId(user_id)
	if err != nil {
		return posts, err
	}
	return posts, nil
}

func (s *postService) AddReaction(reaction models.Reaction) error {
	if err := s.postRepo.AddReactionToPost(reaction); err != nil {
		return fmt.Errorf("error adding or updating reaction: %w", err)
	}

	// Получение поста для уведомления
	post, err := s.postRepo.GetPostByID(reaction.PostID)
	if err != nil {
		return fmt.Errorf("reaction added, but failed to retrieve post for notification: %w", err)
	}

	// Формирование уведомления
	action := "liked"
	if reaction.Vote == -1 {
		action = "disliked"
	}
	message := fmt.Sprintf("Your post '%s' was %s by a user.", post.Title, action)

	notification := models.Notification{
		UserID:    post.AuthorID,
		PostID:    reaction.PostID,
		CommentID: reaction.ID,
		Type:      "new_comment",
		Message:   message,
		CreatedAt: time.Now(),
		IsRead:    false,
	}

	// Асинхронная отправка уведомления
	go func() {
		if err := s.postRepo.CreateNotification(notification); err != nil {
			log.Printf("failed to send notification: %v", err)
		}
	}()

	return nil
}

func (s *postService) GetNotificationsByUserID(user_id int) ([]models.Notification, error) {
	notifications, err := s.postRepo.GetNotificationsForUser(user_id)
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("failed to retrieve notifications: %w", err)
	}

	return notifications, nil
}
