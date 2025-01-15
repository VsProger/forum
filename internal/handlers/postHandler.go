package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/pkg"
)

const maxImageSize = 20 * 1024 * 1024

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	nameFunction := "CreatePost"
	tmpl, err := template.ParseFiles("ui/html/pages/createPost.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, nameFunction)
		return
	}

	if r.Method == http.MethodGet {
		if err := tmpl.Execute(w, nil); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
	} else if r.Method == http.MethodPost {
		// Parse the form
		if err := r.ParseMultipartForm(20 * 1024 * 1024); err != nil { // Limit 20MB
			h.handleError(w, nameFunction, http.StatusBadRequest, fmt.Errorf("unable to parse form: %v", err))
			return
		}

		// Get session and user
		session, err := r.Cookie("session")
		if err != nil {
			log.Fatal(err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		user, err := h.service.GetUserByToken(session.Value)
		if err != nil {
			ErrorHandler(w, http.StatusUnauthorized, nameFunction)
			return
		}

		categories := r.Form["categories"]
		post := models.Post{
			Title: r.FormValue("title"),
			Text:  r.FormValue("text"),
		}

		// Handle file upload
		file, _, err := r.FormFile("image")
		if err != nil && err.Error() != "http: no such file" { // If error is not related to missing file
			h.handleError(w, nameFunction, http.StatusBadRequest, fmt.Errorf("file upload error: %v", err))
			return
		}

		if file != nil {
			// Check if the file size exceeds 20 MB
			fileSize := r.ContentLength
			if fileSize > 20*1024*1024 { // 20 MB limit
				tmpl.Execute(w, struct {
					ErrorText string
				}{
					ErrorText: "The file is too large. Please upload an image smaller than 20 MB.",
				})
				return
			}

			// Read the first 512 bytes of the file to check the MIME type
			buf := make([]byte, 512)
			if _, err := io.ReadFull(file, buf); err != nil {
				h.handleError(w, nameFunction, http.StatusInternalServerError, fmt.Errorf("unable to read file: %v", err))
				return
			}

			// Detect the MIME type
			fileType := http.DetectContentType(buf)
			fmt.Println("Detected file type:", fileType) // Debugging line

			// Define allowed MIME types
			allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
			isValidType := false
			for _, t := range allowedTypes {
				if fileType == t {
					isValidType = true
					break
				}
			}

			// If the file type is not valid, show the error
			if !isValidType {
				tmpl.Execute(w, struct {
					ErrorText string
				}{
					ErrorText: "Unsupported file type. Please upload a JPG, PNG, or GIF image.",
				})
				return
			}

			// Reset the file pointer back to the beginning for further processing
			if _, err := file.Seek(0, io.SeekStart); err != nil {
				h.handleError(w, nameFunction, http.StatusInternalServerError, fmt.Errorf("unable to reset file pointer: %v", err))
				return
			}

			// Save the file if valid
			imageURL, err := h.uploadImage(file)
			if err != nil {
				h.handleError(w, nameFunction, http.StatusInternalServerError, err)
				return
			}
			post.ImageURL = imageURL
		}

		// Assign categories
		for _, name := range categories {
			post.Categories = append(post.Categories, models.Category{Name: name})
		}

		if err := pkg.VallidatePost(post); err != nil {
			result := map[string]interface{}{
				"Post":      post,
				"ErrorText": err.Error(),
			}
			if tmpl.Execute(w, result) != nil {
				ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			}
			return
		}

		post.AuthorID = user.ID
		if err := h.service.PostService.CreatePost(post); err != nil {
			log.Fatal(err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	nameFunction := "getPost"
	tmpl, err := template.ParseFiles("ui/html/pages/post.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "getPost")
		return
	}
	if r.Method == http.MethodGet {
		idStr := r.URL.Path[len("/posts/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		post, err := h.service.GetPostByID(id)
		if err != nil || idStr == "" || len(idStr) > 2 || id > 50 || id <= 0 {
			log.Fatal(err)
			ErrorHandler(w, http.StatusNotFound, nameFunction)
			return
		}
		var username string

		session, err := r.Cookie("session")
		if err == nil {
			user, err := h.service.GetUserByToken(session.Value)
			if err == nil {
				username = user.Username

				if err != nil {
					ErrorHandler(w, http.StatusInternalServerError, nameFunction)
					return
				}
			}
		}

		result := map[string]interface{}{
			"Post":          post,
			"Authenticated": username,
		}

		if err = tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "getPost")
			return
		}
	} else if r.Method == http.MethodPost {
		idStr := r.URL.Path[len("/posts/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		session, err := r.Cookie("session")
		if err != nil {
			ErrorHandler(w, http.StatusUnauthorized, nameFunction)
			return
		}
		user, err := h.service.Auth.GetUserByToken(session.Value)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		post, err := h.service.GetPostByID(id)
		if err != nil {
			if idStr == "" || len(idStr) > 2 || id > 50 || id <= 0 {
				log.Fatal(err)
				ErrorHandler(w, http.StatusNotFound, nameFunction)
				return
			}
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		result := map[string]interface{}{
			"Post":          post,
			"Authenticated": user.Username,
		}
		comment := models.Comment{
			Text:     r.FormValue("text"),
			PostID:   id,
			AuthorID: user.ID,
		}

		if err := h.service.CreateComment(comment); err != nil {
			if err == models.ErrEmptyComment || err == models.ErrInvalidComment || err == models.ErrNotAscii {
				ErrorHandler(w, http.StatusBadRequest, nameFunction)
				return
			}
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		path := "/posts/" + idStr
		http.Redirect(w, r, path, http.StatusSeeOther)
		if err := tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleError(w http.ResponseWriter, functionName string, statusCode int, err error) {
	log.Printf("Error in %s: %v", functionName, err)
	http.Error(w, err.Error(), statusCode)
}

func (h *Handler) userPosts(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("ui/html/pages/home.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "getPost")
		return
	}
	if r.Method == http.MethodGet {
		nameFunction := "userPosts"
		session, err := r.Cookie("session")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		user, err := h.service.GetUserByToken(session.Value)
		if err != nil {

			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		posts, err := h.service.GetPostsByUserId(user.ID)
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		result := map[string]interface{}{
			"Posts":       posts,
			"CurrentUser": user,
			"Username":    user.Username,
			"Role":        *user.Role,
		}
		if err = tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "userPosts")
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) userComments(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("ui/html/pages/home.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "getComments")
		return
	}
	if r.Method == http.MethodGet {
		nameFunction := "userComments"
		session, err := r.Cookie("session")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		user, err := h.service.GetUserByToken(session.Value)
		if err != nil {

			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		posts, err := h.service.GetUserCommentsByUserID(user.ID)
		if err != nil {
			log.Fatal(err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		result := map[string]interface{}{
			"Posts":    posts,
			"Username": user.Username,
		}
		if err = tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "userPosts")
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) addReaction(w http.ResponseWriter, r *http.Request) {
	nameFunction := "addReaction"
	if r.Method == http.MethodPost {
		session, err := r.Cookie("session")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		user, err := h.service.Auth.GetUserByToken(session.Value)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		postId, err := pkg.Atoi(r.FormValue("postId"))
		if err != nil {
			log.Fatal(err)
			ErrorHandler(w, http.StatusNotFound, nameFunction)
			return
		}
		var commentId int
		if r.FormValue("commentId") != "" {
			commentId, err = pkg.Atoi(r.FormValue("commentId"))
			if err != nil {
				ErrorHandler(w, http.StatusBadRequest, nameFunction)
				return
			}
		}
		vote, err := pkg.Atoi(r.FormValue("status"))
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		reaction := models.Reaction{
			UserID:    user.ID,
			PostID:    postId,
			CommentID: commentId,
			Vote:      vote,
		}
		if err := h.service.AddReaction(reaction); err != nil {
			if err == fmt.Errorf("specify either PostId or CommentId, not both") || strings.Contains(err.Error(), "Vote IN (-1, 1)") {
				ErrorHandler(w, http.StatusBadRequest, nameFunction)
				return
			} else if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				log.Fatal(err)
				ErrorHandler(w, http.StatusNotFound, nameFunction)
				return
			}
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}

		path := "/posts/" + r.FormValue("postId")
		http.Redirect(w, r, path, http.StatusSeeOther)
	} else {
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
		return
	}
}

func (h *Handler) uploadImage(file multipart.File) (string, error) {
	// Создаем папку для хранения изображений, если ее нет
	imageDir := "ui/static/uploads"
	if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("unable to create image directory: %w", err)
	}

	// Читаем первые 512 байт для определения MIME-типа
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("unable to read file content: %w", err)
	}
	// Определяем MIME тип
	fileType := http.DetectContentType(buffer)

	// Генерируем уникальное имя для файла
	var fileExtension string
	switch fileType {
	case "image/jpeg":
		fileExtension = ".jpg"
	case "image/png":
		fileExtension = ".png"
	case "image/gif":
		fileExtension = ".gif"
	default:
		return "", fmt.Errorf("unsupported file type: %s", fileType)
	}

	// Вернемся к началу файла, чтобы можно было его скопировать
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("unable to seek file: %w", err)
	}

	// Генерируем имя файла с нужным расширением
	fileName := fmt.Sprintf("%d%s", time.Now().Unix(), fileExtension)
	filePath := fmt.Sprintf("%s/%s", imageDir, fileName)

	// Открываем файл для записи
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %w", err)
	}
	defer outFile.Close()

	// Копируем содержимое файла
	if _, err := io.Copy(outFile, file); err != nil {
		return "", fmt.Errorf("unable to copy file content: %w", err)
	}

	return filePath, nil
}

func (h *Handler) GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	nameFunction := "notifications"

	session, err := r.Cookie("session")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, nameFunction)

		user, err := h.service.GetUserByToken(session.Value)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}

		// Получаем уведомления из базы данных
		notifications, err := h.service.GetNotificationsByUserID(user.ID)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error fetching notifications", http.StatusInternalServerError)
			return
		}

		// Отправляем уведомления как JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"notifications": notifications,
		})
	}
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	nameFunction := "DeletePost"

	fmt.Print(r.Method)

	if r.Method == http.MethodPost {
		idStr := r.URL.Path[len("/postsdelete/"):]

		if idStr == "" {
			log.Println("Error: ID is missing in the URL path")
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			log.Printf("Error converting ID: %v", err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}

		// Get session and user info
		session, err := r.Cookie("session")
		if err != nil {
			ErrorHandler(w, http.StatusUnauthorized, nameFunction)
			return
		}

		user, err := h.service.GetUserByToken(session.Value)
		if err != nil {
			ErrorHandler(w, http.StatusUnauthorized, nameFunction)
			return
		}

		// Fetch the post by ID to verify if it exists
		//post, err := h.service.GetPostByID(id)
		//if err != nil {
		//	ErrorHandler(w, http.StatusNotFound, nameFunction)
		//	return
		//}

		// Ensure the user is the author of the post
		//if post.AuthorID != user.ID {
		//	ErrorHandler(w, http.StatusForbidden, nameFunction)
		//	return
		//}

		// Delete the post (this may include deleting related data like reactions or comments)
		if err := h.service.PostService.DeletePost(id); err != nil {
			log.Fatal(err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		if *user.Role == "admin" && r.URL.Path == "/adminpage" {
			http.Redirect(w, r, "/adminpage", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

		// Respond with a success message or appropriate redirection
		w.WriteHeader(http.StatusNoContent) // 204 No Content for successful deletion
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) editPost(w http.ResponseWriter, r *http.Request) {
	nameFunction := "editpost"
	tmpl, err := template.ParseFiles("ui/html/pages/editpost.html")
	if err != nil {

		ErrorHandler(w, http.StatusInternalServerError, nameFunction)
		return
	}

	// Extract the post ID from the URL
	postID := r.URL.Path[len("/postsedit/"):] // assuming the URL path is "/postsedit/{id}"
	if postID == "" {
		ErrorHandler(w, http.StatusBadRequest, nameFunction)
		return
	}

	postIDInt, err := strconv.Atoi(postID)

	if r.Method == http.MethodGet {
		// Fetch the post from the database using the extracted postID
		post, err := h.service.PostService.GetPostByID(postIDInt)
		if err != nil {
			ErrorHandler(w, http.StatusNotFound, nameFunction)
			return
		}

		// Populate the template with the post data for editing
		result := map[string]interface{}{
			"Post": post,
		}

		if err := tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}

	} else if r.Method == http.MethodPost {
		// Parse the form
		if err := r.ParseMultipartForm(20 * 1024 * 1024); err != nil { // Limit 20MB
			h.handleError(w, nameFunction, http.StatusBadRequest, fmt.Errorf("unable to parse form: %v", err))
			return
		}

		// Get session and user
		session, err := r.Cookie("session")
		if err != nil {
			log.Fatal(err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		user, err := h.service.GetUserByToken(session.Value)
		if err != nil {
			ErrorHandler(w, http.StatusUnauthorized, nameFunction)
			return
		}

		// Get the categories from the form
		categories := r.Form["categories"]
		post := models.Post{
			Title: r.FormValue("title"),
			Text:  r.FormValue("text"),
		}

		// Handle file upload
		file, _, err := r.FormFile("image")
		if err != nil && err.Error() != "http: no such file" { // If error is not related to missing file
			h.handleError(w, nameFunction, http.StatusBadRequest, fmt.Errorf("file upload error: %v", err))
			return
		}

		if file != nil {
			// Check if the file size exceeds 20 MB
			fileSize := r.ContentLength
			if fileSize > 20*1024*1024 { // 20 MB limit
				tmpl.Execute(w, struct {
					ErrorText string
				}{
					ErrorText: "The file is too large. Please upload an image smaller than 20 MB.",
				})
				return
			}

			// Read the first 512 bytes of the file to check the MIME type
			buf := make([]byte, 512)
			if _, err := io.ReadFull(file, buf); err != nil {
				h.handleError(w, nameFunction, http.StatusInternalServerError, fmt.Errorf("unable to read file: %v", err))
				return
			}

			// Detect the MIME type
			fileType := http.DetectContentType(buf)
			fmt.Println("Detected file type:", fileType) // Debugging line

			// Define allowed MIME types
			allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
			isValidType := false
			for _, t := range allowedTypes {
				if fileType == t {
					isValidType = true
					break
				}
			}

			// If the file type is not valid, show the error
			if !isValidType {
				tmpl.Execute(w, struct {
					ErrorText string
				}{
					ErrorText: "Unsupported file type. Please upload a JPG, PNG, or GIF image.",
				})
				return
			}

			// Reset the file pointer back to the beginning for further processing
			if _, err := file.Seek(0, io.SeekStart); err != nil {
				h.handleError(w, nameFunction, http.StatusInternalServerError, fmt.Errorf("unable to reset file pointer: %v", err))
				return
			}

			// Save the file if valid
			imageURL, err := h.uploadImage(file)
			if err != nil {
				h.handleError(w, nameFunction, http.StatusInternalServerError, err)
				return
			}
			post.ImageURL = imageURL
		}

		// Assign categories
		for _, name := range categories {
			post.Categories = append(post.Categories, models.Category{Name: name})
		}

		// Validate post
		// if err := pkg.VallidatePost(post); err != nil {
		// 	result := map[string]interface{}{
		// 		"Post":      post,
		// 		"ErrorText": err.Error(),
		// 	}
		// 	if tmpl.Execute(w, result) != nil {
		// 		ErrorHandler(w, http.StatusInternalServerError, nameFunction)
		// 	}
		// 	return
		// }

		post.AuthorID = user.ID
		post.ID = postIDInt // Set the post ID from the URL

		// Update the post in the database
		if err := h.service.PostService.UpdatePost(post); err != nil {
			log.Fatal(err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}

		// Redirect to the updated post
		http.Redirect(w, r, fmt.Sprintf("/posts/%s", strconv.Itoa(post.ID)), http.StatusSeeOther)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
