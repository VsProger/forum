package posts

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/VsProger/snippetbox/internal/models"
)

type Posts interface {
	CreatePost(post models.Post) error
	GetCategoryByName(name string) (*models.Category, error)
	CreateCategory(name string) error
	GetPostByID(id int) (*models.Post, error)
	GetPosts() ([]models.Post, error)
	CreateComment(comment models.Comment) error
	GetAllPostsByUserId(id int) ([]models.Post, error)
}

type PostRepo struct {
	DB *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{
		DB: db,
	}
}

func (r *PostRepo) CreatePost(post models.Post) error {
	// Начало транзакции
	tx, err := r.DB.Begin()
	if err != nil {
		log.Printf("ошибка начала транзакции: %v", err)
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback() // Если что-то пойдет не так, транзакция будет отменена

	// Подготовка запроса на вставку данных в таблицу Posts
	query := `
	INSERT INTO Posts (AuthorID, Title, Text, CreationTime)
	VALUES (?, ?, ?, datetime('now','+6 hours'));`

	// Выполнение запроса на вставку в рамках транзакции
	res, err := tx.Exec(query, post.AuthorID, post.Title, post.Text)
	if err != nil {
		log.Printf("ошибка выполнения запроса на вставку: %v", err)
		return fmt.Errorf("ошибка выполнения запроса на вставку: %w", err)
	}

	// Получение последнего вставленного ID
	var postID int64
	postID, err = res.LastInsertId()
	if err != nil {
		log.Printf("ошибка получения последнего вставленного ID: %v", err)
		return fmt.Errorf("ошибка получения последнего вставленного ID: %w", err)
	}

	// Вставка категорий в таблицу PostCategory в рамках транзакции
	for _, category := range post.Categories {
		_, err := tx.Exec(`
			INSERT INTO PostCategory (PostID, CategoryID)
			VALUES (?, ?)
		`, postID, category.ID)
		if err != nil {
			log.Printf("ошибка вставки категории в таблицу PostCategory: %v", err)
			return fmt.Errorf("ошибка вставки категории в таблицу PostCategory: %w", err)
		}
	}

	// Если все успешно, коммитим транзакцию
	if err := tx.Commit(); err != nil {
		log.Printf("ошибка коммита транзакции: %v", err)
		return fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return nil
}

func (r *PostRepo) GetPosts() ([]models.Post, error) {
	query := `SELECT p.ID, p.AuthorID, p.Title, p.Text, p.CreationTime, u.Username 
	FROM Posts p
	JOIN User u ON p.AuthorID =u.ID`
	queryCategories := `SELECT ID, Name FROM Category WHERE ID IN (SELECT CategoryID FROM PostCategory WHERE PostID = ?)`
	posts := []models.Post{}
	rows, err := r.DB.Query(query)
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		post := models.Post{}
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Text, &post.CreationTime, &post.Username); err != nil {
			return posts, err
		}
		rows2, err := r.DB.Query(queryCategories, post.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting categories for post %d: %w", post.ID, err)
		}
		defer rows2.Close()

		var categories []models.Category
		for rows2.Next() {
			var category models.Category
			if err := rows2.Scan(&category.ID, &category.Name); err != nil {
				return nil, fmt.Errorf("error scanning category: %w", err)
			}
			categories = append(categories, category)
		}
		if err := rows2.Err(); err != nil {
			return nil, fmt.Errorf("error iterating over categories: %w", err)
		}

		post.Categories = categories
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func (m *PostRepo) Latest() ([]models.Post, error) {
	stmt := `SELECT ID, AuthorID, Title, Text, LikeCount, DislikeCount, CreationTime FROM Posts
    ORDER BY CreationTime DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post

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

func (r *PostRepo) GetCategoryByName(name string) (*models.Category, error) {
	query := `
	SELECT ID, Name
	FROM Category
	WHERE Name = ?`
	row := r.DB.QueryRow(query, name)
	category := models.Category{}
	if err := row.Scan(&category.ID, &category.Name); err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *PostRepo) CreateCategory(name string) error {
	query := "INSERT INTO Category (Name) VALUES (?)"
	_, err := r.DB.Exec(query, name)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) GetPostByID(id int) (*models.Post, error) {
	queryPost := `SELECT p.ID, p.AuthorID, p.Title, p.Text, p.LikeCount, p.DislikeCount, p.CreationTime, u.Username 
	FROM Posts p
	JOIN User u ON p.AuthorID = u.ID
	WHERE p.ID = ?;`
	queryCategories := `SELECT ID, Name FROM Category WHERE ID IN (SELECT CategoryID FROM PostCategory WHERE PostID = ?)`

	post := &models.Post{}
	err := r.DB.QueryRow(queryPost, id).Scan(&post.ID, &post.AuthorID, &post.Title, &post.Text, &post.LikeCount, &post.DislikeCount, &post.CreationTime, &post.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("post not found with ID %d", id)
		}
		return nil, fmt.Errorf("error scanning post: %w", err)
	}

	rows, err := r.DB.Query(queryCategories, id)
	if err != nil {
		return nil, fmt.Errorf("error getting categories for post %d: %w", id, err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("error scanning category: %w", err)
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over categories: %w", err)
	}

	post.Categories = categories

	reactionQuery := `
    SELECT 
        COALESCE(SUM(CASE WHEN Vote = 1 THEN 1 ELSE 0 END), 0) as Likes,
        COALESCE(SUM(CASE WHEN Vote = -1 THEN 1 ELSE 0 END), 0) as Dislikes
    FROM Reaction WHERE PostID = $1
    `
	err = r.DB.QueryRow(reactionQuery, id).Scan(&post.LikeCount, &post.DislikeCount)
	if err != nil {
		return post, err
	}
	commentsQuery := `
	SELECT 
		c.Id, c.Text, c.PostID, c.AuthorID, u.Username,
		COALESCE(SUM(CASE WHEN r.Vote = 1 THEN 1 ELSE 0 END), 0) as Likes,
		COALESCE(SUM(CASE WHEN r.Vote = -1 THEN 1 ELSE 0 END), 0) as Dislikes
	FROM Comment c
	JOIN User u ON c.AuthorID = u.ID
	LEFT JOIN Reaction r ON c.ID = r.CommentID
	WHERE c.PostID = $1
	GROUP BY c.ID, u.Username, c.Text, c.PostID, c.AuthorID
	`
	rows, err = r.DB.Query(commentsQuery, id)
	if err != nil {
		return post, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.Text, &comment.PostID, &comment.AuthorID, &comment.Username, &comment.LikeCount, &comment.DislikeCount); err != nil {
			return post, err
		}
		post.Comment = append(post.Comment, comment)
	}
	if err = rows.Err(); err != nil {
		return post, err
	}
	return post, nil
}

func (p *PostRepo) CreateComment(comment models.Comment) error {
	query := "INSERT INTO Comment (AuthorID, PostID, Text, Username) VALUES ($1, $2, $3, $4)"
	_, err := p.DB.Exec(query, comment.AuthorID, comment.PostID, comment.Text, comment.Username)
	fmt.Println(err)
	return err
}

func (r *PostRepo) GetAllPostsByUserId(id int) ([]models.Post, error) {
	queryPost := `SELECT p.ID, p.AuthorID, p.Title, p.Text, p.LikeCount, p.DislikeCount, p.CreationTime, u.Username 
	FROM Posts p
	JOIN User u ON p.AuthorID = u.ID
	WHERE p.AuthorID = ? ORDER BY CreationTime DESC;`

	queryCategories := `SELECT ID, Name FROM Category WHERE ID IN (SELECT CategoryID FROM PostCategory WHERE PostID = ?)`
	rows, err := r.DB.Query(queryPost, id)
	if err != nil {
		return nil, fmt.Errorf("error getting posts %d: %w", id, err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Text, &post.LikeCount, &post.DislikeCount, &post.CreationTime, &post.Username)
		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("post not found with ID %d", id)
			}
			return nil, fmt.Errorf("error scanning post: %w", err)
		}
		rows2, err := r.DB.Query(queryCategories, post.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting categories for post %d: %w", id, err)
		}
		defer rows2.Close()

		var categories []models.Category
		for rows2.Next() {
			var category models.Category
			if err := rows2.Scan(&category.ID, &category.Name); err != nil {
				return nil, fmt.Errorf("error scanning category: %w", err)
			}
			categories = append(categories, category)
		}
		if err := rows2.Err(); err != nil {
			return nil, fmt.Errorf("error iterating over categories: %w", err)
		}

		post.Categories = categories
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
