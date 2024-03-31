package model

import "github.com/gofrs/uuid"

type User struct {
	Id             uuid.UUID
	Login          string
	HashedPassword []byte
	Name           *string
	Surname        *string
	Email          *string
	Phone          *string
}