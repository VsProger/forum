package models

type Reaction struct {
	ID        int
	UserID    int
	PostID    int // Опционально, может быть nil
	CommentID int // Опционально, может быть nil
	Vote      int // Только 1 или -1
}
