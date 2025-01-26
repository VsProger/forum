package admin

import (
	"database/sql"
	"fmt"
	"github.com/VsProger/snippetbox/internal/models"
	"log"
)

type Admin interface {
	GetUsers() ([]models.User, error)
	UpgradeUser(user_id int) error
	DowngradeUser(user_id int) error
	ReportPost(postID int, userID int, reason string) error
	GetReports() ([]models.Report, error)
	RequestRole(user_id int) error
	ApproveRequest(user_id int) error
	RejectRequest(user_id int) error
	GetRequests() ([]models.User, error)
	CheckRequest(user_id int) (bool, error)
}

type AdminRepo struct {
	DB *sql.DB
}

func NewAdminRepo(db *sql.DB) *AdminRepo {
	return &AdminRepo{
		DB: db,
	}
}

// GetUsers retrieves all users from the database
func (r *AdminRepo) GetUsers() ([]models.User, error) {
	query := "SELECT * FROM User WHERE Role != 'admin'"

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.GoogleID, &user.GitHubID, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	log.Printf("Retrieved %d users", len(users)) // Add logging to check the number of users retrieved

	return users, nil
}

func (r *AdminRepo) UpgradeUser(user_id int) error {
	var currentRole string
	query := "SELECT Role FROM User WHERE ID = ?"
	err := r.DB.QueryRow(query, user_id).Scan(&currentRole)
	if err != nil {
		return fmt.Errorf("failed to fetch user role: %w", err)
	}
	if currentRole == "admin" {
		return fmt.Errorf("cannot upgrade an admin user")
	}

	updateQuery := "UPDATE User SET Role = 'moderator' WHERE ID = ?"
	_, err = r.DB.Exec(updateQuery, user_id)
	if err != nil {
		return fmt.Errorf("failed to upgrade user: %w", err)
	}
	return nil
}

func (r *AdminRepo) DowngradeUser(user_id int) error {
	var currentRole string
	query := "SELECT Role FROM User WHERE ID = ?"
	err := r.DB.QueryRow(query, user_id).Scan(&currentRole)
	if err != nil {
		return fmt.Errorf("failed to fetch user role: %w", err)
	}
	if currentRole == "admin" {
		return fmt.Errorf("cannot downgrade an admin user")
	}

	updateQuery := "UPDATE User SET Role = 'user' WHERE ID = ?"
	_, err = r.DB.Exec(updateQuery, user_id)
	if err != nil {
		return fmt.Errorf("failed to upgrade user: %w", err)
	}
	return nil
}

func (r *AdminRepo) ReportPost(postID int, userID int, reason string) error {
	stmt := "INSERT INTO Report (PostID, UserID, Reason) VALUES (?, ?, ?)"
	_, err := r.DB.Exec(stmt, postID, userID, reason)
	if err != nil {
		return fmt.Errorf("error reporting post: %w", err)
	}
	return nil
}

func (r *AdminRepo) GetReports() ([]models.Report, error) {
	query := `
SELECT 
    u.ID AS UserID,
    u.Username AS UserName,
    u.Email AS UserEmail,
    p.ID AS PostID,
    p.Title AS PostTitle,
    r.Reason AS ReportReason
FROM 
    Report r
JOIN 
    User u ON r.UserID = u.ID
JOIN 
    Posts p ON r.PostID = p.ID

`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var report models.Report
		if err := rows.Scan(&report.UserID, &report.UserName, &report.UserEmail, &report.PostID, &report.PostTitle, &report.ReportReason); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		reports = append(reports, report)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return reports, nil
}

// Add user who requested a moderator role to the Requests table, if not already present
func (r *AdminRepo) RequestRole(user_id int) error {
	stmt := "INSERT INTO Requests (UserID) VALUES (?)"
	_, err := r.DB.Exec(stmt, user_id)
	if err != nil {
		return fmt.Errorf("failed to request role: %w", err)
	}
	return nil
}

func (r *AdminRepo) ApproveRequest(user_id int) error {
	// Check if the user has an existing request
	var count int
	query := "SELECT COUNT(*) FROM Requests WHERE UserID = ?"
	err := r.DB.QueryRow(query, user_id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check for existing request: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("no request found for user")
	}

	// Update the user in the User table to a moderator
	stmt := "UPDATE User SET Role = 'moderator' WHERE ID = ?"
	_, err = r.DB.Exec(stmt, user_id)
	if err != nil {
		return fmt.Errorf("failed to approve request: %w", err)
	}

	// Delete the user from the Requests table
	deleteStmt := "DELETE FROM Requests WHERE UserID = ?"
	_, err = r.DB.Exec(deleteStmt, user_id)
	if err != nil {
		return fmt.Errorf("failed to delete request: %w", err)
	}

	return nil
}

func (r *AdminRepo) RejectRequest(user_id int) error {
	// Check if the user has an existing request
	var count int
	query := "SELECT COUNT(*) FROM Requests WHERE UserID = ?"
	err := r.DB.QueryRow(query, user_id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check for existing request: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("no request found for user")
	}

	// Delete the user from the Requests table
	stmt := "DELETE FROM Requests WHERE UserID = ?"
	_, err = r.DB.Exec(stmt, user_id)
	if err != nil {
		return fmt.Errorf("failed to reject request: %w", err)
	}

	return nil
}

func (r *AdminRepo) GetRequests() ([]models.User, error) {
	query := `SELECT u.ID, u.Username, u.Email FROM Requests r JOIN User u ON r.UserID = u.ID`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)

	}
	defer rows.Close()

	var requests []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		requests = append(requests, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return requests, nil
}

func (r *AdminRepo) CheckRequest(user_id int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM Requests WHERE UserID = ?"
	err := r.DB.QueryRow(query, user_id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check for existing request: %w", err)
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
