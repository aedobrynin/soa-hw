package model

type PostId = string

type Post struct {
	Id       PostId
	AuthorId UserId
	Content  string
}
