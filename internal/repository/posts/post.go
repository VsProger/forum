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
	GetCategoryByName(name string) ([]*models.Category, error)
	CreateCategory(name string) error
	GetPostByID(id int) (*models.Post, error)
	GetPosts() ([]models.Post, error)
	CreateComment(comment models.Comment) error
	GetAllPostsByUserId(id int) ([]models.Post, error)
	AddReactionToPost(reaction models.Reaction) error
	AddReactionToComment(reaction models.Reaction) error
	CreateNotification(notification models.Notification) error
	GetNotificationsForUser(userID int) ([]models.Notification, error)
	MarkNotificationAsRead(notificationID int) error
	NotifyUser(userID int, message string) error
	GetUserCommentsByUserID(userID int) ([]models.Post, error)
	DeletePost(postID int) error
	UpdatePost(post models.Post) error

	GetUsers() ([]models.User, error)
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
	tx, err := r.DB.Begin()
	if err != nil {
		log.Printf("error starting transaction: %v", err)
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
	INSERT INTO Posts (AuthorID, Title, Text, ImageURL, CreationTime)
	VALUES (?, ?, ?, ?, datetime('now','+6 hours'));`
	res, err := tx.Exec(query, post.AuthorID, post.Title, post.Text, post.ImageURL)
	if err != nil {
		log.Printf("error inserting post: %v", err)
		return fmt.Errorf("error inserting post: %w", err)
	}

	postID, err := res.LastInsertId()
	for _, category := range post.Categories {
		_, err := tx.Exec(`
			INSERT INTO PostCategory (PostID, CategoryID)
			VALUES (?, ?)
		`, postID, category.ID)
		if err != nil {
			log.Printf("error inserting category: %v", err)
			return fmt.Errorf("error inserting category: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error committing transaction: %v", err)
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (r *PostRepo) GetPosts() ([]models.Post, error) {
	query := `SELECT p.ID, p.AuthorID, p.Title, p.Text, p.CreationTime, p.ImageURL, u.Username 
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
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Text, &post.CreationTime, &post.ImageURL, &post.Username); err != nil {
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
	stmt := `SELECT ID, AuthorID, Title, Text, LikeCount, DislikeCount, ImageURL, CreationTime FROM Posts
    ORDER BY CreationTime DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post

		err = rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Text, &p.LikeCount, &p.DislikeCount, &p.ImageURL, &p.CreationTime)
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

func (r *PostRepo) GetCategoryByName(name string) ([]*models.Category, error) {
	query := `
	SELECT ID, Name
	FROM Category
	WHERE Name = ?`
	rows, err := r.DB.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure the rows are closed after we're done with them

	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	// Check if there was an error during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
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
	queryPost := `SELECT p.ID, p.AuthorID, p.Title, p.Text, p.LikeCount, p.DislikeCount, p.ImageURL, p.CreationTime, u.Username 
	FROM Posts p
	JOIN User u ON p.AuthorID = u.ID
	WHERE p.ID = ?;`
	queryCategories := `SELECT ID, Name FROM Category WHERE ID IN (SELECT CategoryID FROM PostCategory WHERE PostID = ?)`

	post := &models.Post{}
	err := r.DB.QueryRow(queryPost, id).Scan(&post.ID, &post.AuthorID, &post.Title, &post.Text, &post.LikeCount, &post.DislikeCount, &post.ImageURL, &post.CreationTime, &post.Username)
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
	queryPost := `SELECT p.ID, p.AuthorID, p.Title, p.Text, p.LikeCount, p.DislikeCount, p.ImageURL, p.CreationTime, u.Username 
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
		err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Text, &post.LikeCount, &post.DislikeCount, &post.ImageURL, &post.CreationTime, &post.Username)
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

func (p *PostRepo) AddReactionToPost(reaction models.Reaction) error {
	var existingreaction int
	var existingreactionId int
	err := p.DB.QueryRow(`
        SELECT ID, Vote FROM Reaction 
        WHERE UserID = $1 AND PostID = $2`,
		reaction.UserID,
		reaction.PostID,
	).Scan(&existingreactionId, &existingreaction)
	if err == sql.ErrNoRows {
		_, err := p.DB.Exec(`
            INSERT INTO Reaction (UserID, PostID, Vote) 
            VALUES ($1, $2, $3)`,
			reaction.UserID,
			reaction.PostID,
			reaction.Vote,
		)
		return err
	} else if err != nil {
		return err
	} else {
		if existingreaction == reaction.Vote {
			_, err := p.DB.Exec("DELETE FROM Reaction WHERE ID = $1", existingreactionId)
			return err
		} else {
			_, err := p.DB.Exec(`
                UPDATE Reaction 
                SET Vote = $1 
                WHERE ID = $2`,
				reaction.Vote,
				existingreactionId,
			)
			return err
		}
	}
}

func (p *PostRepo) AddReactionToComment(reaction models.Reaction) error {
	var existingreaction int
	var existingreactionId int
	err := p.DB.QueryRow(`
        SELECT ID, Vote FROM Reaction 
        WHERE UserID = $1 AND CommentID = $2`,
		reaction.UserID,
		reaction.CommentID,
	).Scan(&existingreactionId, &existingreaction)
	if err == sql.ErrNoRows {
		_, err := p.DB.Exec(`
            INSERT INTO Reaction (UserID, CommentID, Vote) 
            VALUES ($1, $2, $3)`,
			reaction.UserID,
			reaction.CommentID,
			reaction.Vote,
		)
		return err
	} else if err != nil {
		return err
	} else {
		if existingreaction == reaction.Vote {
			_, err := p.DB.Exec("DELETE FROM Reaction WHERE ID = $1", existingreactionId)
			return err
		} else {
			_, err := p.DB.Exec(`
                UPDATE Reaction 
                SET Vote = $1 
                WHERE ID = $2`,
				reaction.Vote,
				existingreactionId,
			)
			return err
		}
	}
}

func (r *PostRepo) CreateNotification(notification models.Notification) error {
	query := `
		INSERT INTO Notifications (UserID, PostID, CommentID, Type, Message, CreatedAt, IsRead)
		VALUES (?, ?, ?, ?, ?, ?, false)
	`
	_, err := r.DB.Exec(query, notification.UserID, notification.PostID, notification.CommentID, notification.Type, notification.Message, notification.CreatedAt)
	if err != nil {
		return fmt.Errorf("error creating notification: %w", err)
	}
	return nil
}

func (r *PostRepo) GetNotificationsForUser(userID int) ([]models.Notification, error) {
	query := `
    SELECT ID, UserID, PostID, CommentID, Type, Message, CreatedAt, IsRead
    FROM Notifications
    WHERE UserID = ? ORDER BY CreatedAt DESC
    `
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.PostID, &notification.CommentID, &notification.Type, &notification.Message, &notification.CreatedAt, &notification.IsRead); err != nil {
			return nil, fmt.Errorf("error scanning notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	// Return notifications and nil for error (no error occurred)
	return notifications, nil
}

func (r *PostRepo) GetUserCommentsByUserID(userID int) ([]models.Post, error) {
	query := `
	SELECT DISTINCT 
    p.ID, 
    p.AuthorID, 
    p.Title, 
    p.Text, 
    p.LikeCount,
	p.DislikeCount,
	p.ImageURL,  
    p.CreationTime
FROM 
    Posts p
JOIN 
    Comment c 
ON 
    p.ID = c.PostID
WHERE 
    c.AuthorID = $1`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching notifications: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Text, &post.LikeCount, &post.DislikeCount, &post.ImageURL, &post.CreationTime); err != nil {
			return nil, fmt.Errorf("error scanning notification: %w", err)
		}
		posts = append(posts, post)
	}

	// Return notifications and nil for error (no error occurred)
	return posts, nil
}

func (r *PostRepo) MarkNotificationAsRead(notificationID int) error {
	query := `
    UPDATE Notifications
    SET IsRead = true
    WHERE ID = ?
    `
	_, err := r.DB.Exec(query, notificationID)
	if err != nil {
		return fmt.Errorf("error marking notification as read: %w", err)
	}
	return nil
}

func (r *PostRepo) NotifyUser(userID int, message string) error {
	log.Printf("Notification sent to user %d: %s", userID, message)
	return nil
}

func (r *PostRepo) DeletePost(postID int) error {
	// Start a transaction to ensure all related data is deleted correctly
	tx, err := r.DB.Begin()
	if err != nil {
		log.Printf("error starting transaction: %v", err)
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete related reactions
	_, err = tx.Exec("DELETE FROM Reaction WHERE PostID = ?", postID)
	if err != nil {
		log.Printf("error deleting reactions for post: %v", err)
		return fmt.Errorf("error deleting reactions for post: %w", err)
	}

	// Delete related comments
	_, err = tx.Exec("DELETE FROM Comment WHERE PostID = ?", postID)
	if err != nil {
		log.Printf("error deleting comments for post: %v", err)
		return fmt.Errorf("error deleting comments for post: %w", err)
	}

	// Finally, delete the post
	_, err = tx.Exec("DELETE FROM Posts WHERE ID = ?", postID)
	if err != nil {
		log.Printf("error deleting post: %v", err)
		return fmt.Errorf("error deleting post: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("error committing transaction: %v", err)
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (r *PostRepo) UpdatePost(post models.Post) error {
	// Prepare the dynamic query
	query := "UPDATE Posts SET "
	params := make([]interface{}, 0) // Create a slice of type []interface{}
	paramIndex := 1

	// Dynamically build parts of the SQL query
	if post.Title != "" {
		query += fmt.Sprintf("Title = $%d, ", paramIndex)
		params = append(params, post.Title)
		paramIndex++
	}
	if post.Text != "" {
		query += fmt.Sprintf("Text = $%d, ", paramIndex)
		params = append(params, post.Text)
		paramIndex++
	}
	if post.ImageURL != "" {
		query += fmt.Sprintf("ImageURL = $%d, ", paramIndex)
		params = append(params, post.ImageURL)
		paramIndex++
	}

	// Remove the trailing comma and space, then add WHERE clause
	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", paramIndex)
	params = append(params, post.ID)

	// Manually pass each element of params
	switch len(params) {
	case 1:
		_, err := r.DB.Exec(query, params[0])
		return err
	case 2:
		_, err := r.DB.Exec(query, params[0], params[1])
		return err
	case 3:
		_, err := r.DB.Exec(query, params[0], params[1], params[2])
		return err
	case 4:
		_, err := r.DB.Exec(query, params[0], params[1], params[2], params[3])
		return err
	default:
		return fmt.Errorf("unsupported number of parameters: %d", len(params))
	}
}

// GetUsers retrieves all users from the database
func (r *PostRepo) GetUsers() ([]models.User, error) {
	query := "SELECT ID, Username, Email, Role FROM User;"

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}
