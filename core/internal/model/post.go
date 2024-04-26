package model

type PostID = string

type Post struct {
	ID       PostID
	AuthorID UserID
	Content  string
}
