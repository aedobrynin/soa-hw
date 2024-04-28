package model

type PostLike struct {
	UserID UserID `json:"user_id"`
	PostID PostID `json:"post_id"`
}
