package model

type OrderBy = uint8

const (
	OrderByLikesCnt OrderBy = 0
	OrderByViewsCnt OrderBy = 1
)

type GetTopPostsRequest struct {
	Limit   uint64
	OrderBy OrderBy
}
