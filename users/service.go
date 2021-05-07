package users

import (
	"errors"
	"sync"
)

var (
	users = map[string]User{
		"apm": {Name: "Parsa", Username: "apm"},
		"ftm": {Name: "Fatemeh", Username: "ftm"},
	}
	lock sync.RWMutex

	ErrAlreadyExist = errors.New("user already exists")
	ErrNotFound     = errors.New("user not found")
)

type User struct {
	Name     string `json:"name" binding:"required,alphanumunicode,min=3"`
	Username string `json:"username" binding:"required,alphanum,min=3"`
}

func Add(u User) error {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := users[u.Username]; ok {
		return ErrAlreadyExist
	}

	users[u.Username] = u
	return nil
}

func Find(username string) (*User, error) {
	lock.RLock()
	defer lock.RUnlock()

	u, ok := users[username]
	if !ok {
		return nil, ErrNotFound
	}
	return &u, nil
}

func GetAll() []User {
	lock.RLock()
	defer lock.RUnlock()

	arr := make([]User, 0, len(users))
	for _, u := range users {
		arr = append(arr, u)
	}
	return arr
}
