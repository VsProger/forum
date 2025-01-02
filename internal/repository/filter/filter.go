package filter

import (
	"database/sql"
	"fmt"

	"github.com/VsProger/snippetbox/internal/models"
)

type FilterRepo struct {
	DB *sql.DB
}

type Filter interface {
	GetPostsByCategories(categories []int) ([]models.Post, error)
	GetUsersByLikedPosts(userId int) ([]models.Post, error)
	GetUsersByDislikedPosts(userID int) ([]models.Post, error)
}

func NewFilterRepo(db *sql.DB) *FilterRepo {
	return &FilterRepo{
		DB: db,
	}
}

func (f *FilterRepo) GetPostsByCategories(categories []int) ([]models.Post, error) {
	var inParams string
	for i := range categories {
		if i > 0 {
			inParams += ", "
		}
		inParams += "?"
	}

	// Подготовка запроса с плейсхолдерами
	query := fmt.Sprintf(`
    SELECT p.ID, p.Title, p.Text, p.CreationTime, p.AuthorID, u.Username, 
           GROUP_CONCAT(c.Name) as Categories
    FROM Posts p
    JOIN User u ON p.AuthorID = u.ID
    JOIN PostCategory pc ON p.ID = pc.PostID
    JOIN Category c ON pc.CategoryID = c.ID
    WHERE pc.CategoryID IN (%s)
    GROUP BY p.ID, p.Title, p.Text, p.CreationTime, p.AuthorID, u.Username
	`, inParams)

	// Создаем срез аргументов
	args := make([]interface{}, len(categories))
	for i, v := range categories {
		args[i] = v
	}

	// Используем args... для распаковки аргументов при передаче в f.DB.Query
	result := []models.Post{}

	for i := 0; i < len(args); i++ {
		rows, err := f.DB.Query(query, args[i]) // Использование args... для распаковки
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		// Обрабатываем результаты
		for rows.Next() {
			var post models.Post
			if err := rows.Scan(&post.ID, &post.Title, &post.Text, &post.CreationTime, &post.AuthorID, &post.Username, &post.Category); err != nil {
				return nil, err
			}
			post.Categories, err = f.getAllCategoriesByPostId(post.ID)
			if err != nil {
				return nil, err
			}
			result = append(result, post)
		}

		// Проверка на ошибки после завершения итерации
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (f *FilterRepo) GetUsersByLikedPosts(userID int) ([]models.Post, error) {
	query := `
	SELECT p.ID, p.Title, p.Text, p.CreationTime, p.AuthorID, u.Username
	FROM Posts p
	JOIN User u ON p.AuthorID = u.ID
	JOIN Reaction r ON p.ID = r.PostID
	WHERE r.UserID = $1 AND r.Vote = 1
	`
	result := []models.Post{}
	rows, err := f.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Text, &post.CreationTime, &post.AuthorID, &post.Username); err != nil {
			return nil, err
		}
		post.Categories, err = f.getAllCategoriesByPostId(post.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (f *FilterRepo) GetUsersByDislikedPosts(userID int) ([]models.Post, error) {
	query := `
	SELECT p.ID, p.Title, p.Text, p.CreationTime, p.AuthorID, u.Username
	FROM Posts p
	JOIN User u ON p.AuthorID = u.ID
	JOIN Reaction r ON p.ID = r.PostID
	WHERE r.UserID = $1 AND r.Vote = -1
	`
	result := []models.Post{}
	rows, err := f.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Text, &post.CreationTime, &post.AuthorID, &post.Username); err != nil {
			return nil, err
		}
		post.Categories, err = f.getAllCategoriesByPostId(post.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
