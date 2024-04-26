package model

import "github.com/google/uuid"

type UserId = uuid.UUID

type User struct {
	Id             UserId
	Login          string
	HashedPassword []byte
	Name           *string
	Surname        *string
	Email          *string
	Phone          *string
}
