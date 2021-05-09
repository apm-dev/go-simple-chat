package domain

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func NewUser(name, email, pass string) (*User, error) {
	const op string = "domain.user.NewUser"

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &User{
		Name:     name,
		Email:    email,
		Password: string(hash),
	}, nil
}
