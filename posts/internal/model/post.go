package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Post struct {
	Id        uuid.UUID
	AuthorId  uuid.UUID
	Content   string
	CreatedTs time.Time
	UpdatedTs time.Time
}
