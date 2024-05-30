package model

type PostView struct {
	UserID       UserID `json:"user_id"`
	PostID       PostID `json:"post_id"`
	PostAuthorID UserID `json:"post_author_id"`
}
