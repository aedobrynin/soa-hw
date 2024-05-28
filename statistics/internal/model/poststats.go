package model

import "github.com/google/uuid"

type PostStatistics struct {
	likes_cnt uint64
	views_cnt uint64
}

// TODO: reuse that kind of typedefs in other services
type PostID = uuid.UUID
