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

func (m *PostModel) Insert(authorID int, title string, text string, categoryID int) (int, error) {
	stmt := `INSERT INTO Posts (AuthorID, Title, Text, LikeCount, DislikeCount, CreationTime, CategoryID)
    VALUES(?, ?, ?, 0, 0, DATETIME('NOW'), ?)`

	result, err := m.DB.Exec(stmt, authorID, title, text, categoryID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PostModel) Get(id int) (Post, error) {
	var p Post

	stmt := `SELECT ID, AuthorID, Title, Text, LikeCount, DislikeCount, CreationTime, CategoryID FROM Posts
    WHERE ID = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.AuthorID, &p.Title, &p.Text, &p.LikeCount, &p.DislikeCount, &p.CreationTime, &p.CategoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Post{}, ErrNoRecord
		} else {
			return Post{}, err
		}
	}
	return p, nil
}

func (m *PostModel) Latest() ([]Post, error) {
	stmt := `SELECT ID, AuthorID, Title, Text, LikeCount, DislikeCount, CreationTime, CategoryID FROM Posts
    ORDER BY CreationTime DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post

		err = rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Text, &p.LikeCount, &p.DislikeCount, &p.CreationTime, &p.CategoryID)
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
