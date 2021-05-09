package domain

type UserRepo interface {
	Add(u User) (string, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]User, error)
}
