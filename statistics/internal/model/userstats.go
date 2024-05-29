package model

import "github.com/google/uuid"

type UserStatistics struct {
	UserID     UserID
	LikesCount uint64
}

// TODO: reuse that kind of typedefs in other services
type UserID = uuid.UUID
