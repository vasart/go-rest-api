package model

import "crypto/subtle"

type User struct {
	ID			string	`json:"id"`
	Username	string	`json:"username"`
	Password	string	`json:"password"`
}

type Credentials struct {
	Username	string	`json:"username"`
	Password	string	`json:"password"`	
}

type UserRepository interface {
	CreateUser(u *User) error
	GetByUsername(username string) (*User, error)
	Login(c Credentials) (*User, error)
}

func (u *User) CheckPassword(password string) bool {
	return subtle.ConstantTimeCompare([]byte(u.Password), []byte(password)) == 1
}