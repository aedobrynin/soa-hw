package model

import (
	"time"

	"github.com/google/uuid"
)

type UserID = string

type PostID = uuid.UUID

type Post struct {
	ID        PostID
	AuthorID  UserID
	Content   string
	CreatedTs time.Time
	UpdatedTs time.Time
}
