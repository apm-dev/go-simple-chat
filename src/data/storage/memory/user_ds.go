package memory

import (
	"sync"

	"apm.dev/go-simple-chat/src/data/repository"
	"apm.dev/go-simple-chat/src/domain"
)

type userDS struct {
	sync.RWMutex
	users map[string]domain.User
}

func NewUserDS() repository.UserDataSource {
	return &userDS{
		users: make(map[string]domain.User),
	}
}

func (ds *userDS) Insert(u domain.User) (string, error) {
	ds.Lock()
	defer ds.Unlock()

	if _, ok := ds.users[u.Email]; ok {
		return "", domain.ErrUserAlreadyExists
	}

	ds.users[u.Email] = u
	return u.Email, nil
}

func (ds *userDS) GetByEmail(e string) (*domain.User, error) {
	ds.RLock()
	defer ds.RUnlock()

	u, ok := ds.users[e]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return u.Clone(), nil
}

func (ds *userDS) GetAll() ([]domain.User, error) {
	ds.RLock()
	defer ds.RUnlock()

	if len(ds.users) == 0 {
		return nil, domain.ErrUserNotFound
	}

	users := make([]domain.User, 0, len(ds.users))
	for _, u := range ds.users {
		users = append(users, u)
	}

	return users, nil
}
