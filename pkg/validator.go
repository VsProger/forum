package pkg

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/VsProger/snippetbox/internal/models"
	"github.com/gofrs/uuid"
)

var (
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidUsername  = errors.New("invalid username")
	ErrInvalidEmail     = errors.New("invalid email address")
	ErrTitleNotAscii    = errors.New("Title is not in Ascii")
	ErrTextNotAscii     = errors.New("Text is not in Ascii")
	ErrCategoryNotFound = errors.New("Category not found")
	ErrTitleLength      = errors.New("Length of title should be between 4 and 30")
	ErrTextLength       = errors.New("Length of text should be between 4 and 600")
	ErrWordsLength      = errors.New("Each word should be less than 30 letters")
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// func PermittedValue[T comparable](value T, permittedValues ...T) bool {
// 	return slices.Contains(permittedValues, value)
// }

func VallidatePost(post models.Post) error {
	post.Title = strings.TrimSpace(post.Title)

	post.Text = strings.TrimSpace(post.Text)
	if ok := isTextAscii(post.Title); !ok {
		return ErrTitleNotAscii
	}
	if ok := isTextAscii(post.Text); !ok {
		return ErrTextNotAscii
	}
	if len(post.Categories) == 0 {
		return ErrCategoryNotFound
	}
	if len(post.Title) < 4 || len(post.Title) > 30 {
		return ErrTitleLength
	}
	if len(post.Text) < 4 || len(post.Text) > 600 {
		return ErrTextLength
	}

	return nil
}

func isTextAscii(text string) bool {
	for _, ch := range text {
		if (int(ch) < 32 || int(ch) > 126) && int(ch) != 10 && int(ch) != 9 {
			return false
		}
	}
	return true
}

func ValidateComment(comment models.Comment) error {
	trimmedText := strings.TrimSpace(comment.Text)
	if trimmedText == "" {
		return models.ErrEmptyComment
	}
	if len(trimmedText) < 4 || len(trimmedText) > 200 {
		return models.ErrInvalidComment
	}
	if !isTextAscii(trimmedText) {
		return models.ErrNotAscii
	}
	return nil
}

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]{1,64}@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,63}$`)

	if !re.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrInvalidPassword
	}
	hasUpperCase := false
	hasLowerCase := false
	hasDigit := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpperCase = true
		case 'a' <= char && char <= 'z':
			hasLowerCase = true
		case '0' <= char && char <= '9':
			hasDigit = true
		}
	}

	if !hasUpperCase || !hasLowerCase || !hasDigit {
		return ErrInvalidPassword
	}

	return nil
}

func ValidateUsername(username string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9_.-]{3,20}$`)
	if !re.MatchString(username) {
		return ErrInvalidUsername
	}
	return nil
}

func CheckPasswordHash(password, hash string) bool {
	//if password != hash {
	//	return false
	//}
	//return true
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken() string {
	u, err := uuid.NewV4()
	if err != nil {
		log.Print(err)
	}
	return u.String()
}

func Atoi(s string) (int, error) {
	if strings.TrimLeft(s, "0") != s || strings.Contains(s, "+") {
		return 0, ErrWordsLength
	}
	if num, err := strconv.Atoi(s); err == nil {
		return num, nil
	} else {
		return 0, err
	}
}
