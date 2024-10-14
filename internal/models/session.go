package models

import "time"

type Session struct {
	ID      int
	Token   string
	ExpTime time.Time
	UserID  int
}
