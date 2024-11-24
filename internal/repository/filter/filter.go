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

	// Prepare the query with placeholders
	query := fmt.Sprintf(`
    SELECT p.ID, p.Title, p.Text, p.CreationTime, p.AuthorID, u.Username, 
           GROUP_CONCAT(c.Name) as Categories
    FROM Post p
    JOIN User u ON p.AuthorID = u.ID
    JOIN PostCategory pc ON p.ID = pc.PostID
    JOIN Category c ON pc.CategoryID = c.ID
    WHERE pc.CategoryID IN (%s)
    GROUP BY p.ID, p.Title, p.Text, p.CreationTime, p.AuthorID, u.Username
	`, inParams)

	// Create the args slice as []interface{}
	args := make([]interface{}, len(categories))
	for i, v := range categories {
		args[i] = v
	}

	// Instead of using args..., pass args directly as a slice
	result := []models.Post{}
	rows, err := f.DB.Query(query, args) // Pass args as a single slice, not spread
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the results and scan them into the result slice
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

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (f *FilterRepo) GetUsersByLikedPosts(userID int) ([]models.Post, error) {
	query := `
	SELECT p.ID, p.Title, p.Text, p.CreationTime, p.AuthorID, u.Username
	FROM Post p
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
