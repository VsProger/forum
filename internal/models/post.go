package models

import (
	"database/sql"
	"errors"
	"time"
)

type Post struct {
	ID           int
	AuthorID     int
	Title        string
	Text         string
	LikeCount    int
	DislikeCount int
	CreationTime time.Time
	CategoryID   int
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(authorID int, title string, text string, categoryIDs []int) (int, error) {
	stmt := `INSERT INTO Posts (AuthorID, Title, Text, LikeCount, DislikeCount, CreationTime)
    VALUES(?, ?, ?, 0, 0, DATETIME('NOW'))`

	result, err := m.DB.Exec(stmt, authorID, title, text)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	stmt = `INSERT INTO PostCategory (PostID, CategoryID) VALUES (?, ?)`

	for _, categoryID := range categoryIDs {
		_, err := m.DB.Exec(stmt, postID, categoryID)
		if err != nil {
			return 0, err // Consider transaction rollback logic here
		}
	}

	return int(postID), nil
}

func (m *PostModel) Get(id int) (Post, []Category, error) {
	var p Post

	// Fetch the post details
	stmt := `SELECT ID, AuthorID, Title, Text, LikeCount, DislikeCount, CreationTime FROM Posts WHERE ID = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.AuthorID, &p.Title, &p.Text, &p.LikeCount, &p.DislikeCount, &p.CreationTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Post{}, nil, ErrNoRecord
		} else {
			return Post{}, nil, err
		}
	}

	// Fetch the category IDs associated with the post
	categoryStmt := `SELECT c.ID, c.Name FROM Category AS c JOIN PostCategory AS pc ON c.ID = pc.CategoryID WHERE pc.PostID = ?`
	rows, err := m.DB.Query(categoryStmt, id)
	if err != nil {
		return Post{}, nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return Post{}, nil, err
		}
		categories = append(categories, cat)
	}

	if err = rows.Err(); err != nil {
		return Post{}, nil, err
	}

	return p, categories, nil
}

func (m *PostModel) Latest() ([]Post, error) {
	stmt := `SELECT ID, AuthorID, Title, Text, LikeCount, DislikeCount, CreationTime FROM Posts
    ORDER BY CreationTime DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post

		err = rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Text, &p.LikeCount, &p.DislikeCount, &p.CreationTime)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
