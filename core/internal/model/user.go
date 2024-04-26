package model

import "github.com/google/uuid"

type UserID = uuid.UUID

type User struct {
	ID             UserID
	Login          string
	HashedPassword []byte
	Name           *string
	Surname        *string
	Email          *string
	Phone          *string
}
