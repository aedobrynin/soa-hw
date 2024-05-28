package model

import "github.com/google/uuid"

type PostStatistics struct {
	LikesCnt uint64
	ViewsCnt uint64
}

type CutPostStatistics struct {
	PostID   PostID
	LikesCnt *uint64
	ViewsCnt *uint64
}

// TODO: reuse that kind of typedefs in other services
type PostID = uuid.UUID
