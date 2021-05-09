package authing

import (
	"apm.dev/go-simple-chat/src/domain"
	"apm.dev/go-simple-chat/src/pkg/logger"
	"github.com/pkg/errors"
)

type Service interface {
	Register(name, email, pass string) (string, error)
	Login(email, pass string) (string, error)
	Authorize(token string) (*UserClaims, error)
}

func NewService(rp domain.UserRepo, jwt *JWTManager) Service {
	return &service{rp, jwt}
}

type service struct {
	repo domain.UserRepo
	jwt  *JWTManager
}

func (svc *service) Register(name, email, pass string) (string, error) {
	const op string = "domain.authing.service.Register"
	// create domain user object
	user, err := domain.NewUser(name, email, pass)
	if err != nil {
		logger.Error(errors.Wrap(err, op))
		return "", domain.ErrInternalServer
	}
	// persist user
	id, err := svc.repo.Add(*user)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return "", domain.ErrUserAlreadyExists
		}
		logger.Error(errors.Wrap(err, op))
		return "", domain.ErrInternalServer
	}

	return id, nil
}

func (svc *service) Login(email, pass string) (string, error) {
	const op string = "domain.authing.service.Login"
	// fetch user from db
	user, err := svc.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrUserNotFound
		}
		logger.Error(errors.Wrap(err, op))
		return "", domain.ErrInternalServer
	}
	// check password
	if !user.IsCorrectPassword(pass) {
		logger.Info(errors.Wrap(domain.ErrWrongCredentials, user.Email))
		return "", domain.ErrWrongCredentials
	}
	// generate jwt token with user claims
	token, err := svc.jwt.Generate(*user)
	if err != nil {
		logger.Error(errors.Wrap(err, op))
		return "", domain.ErrInternalServer
	}

	return token, nil
}

func (svc *service) Authorize(token string) (*UserClaims, error) {
	const op string = "domain.authing.service.Authorize"

	// verify and get claims of token
	claims, err := svc.jwt.Verify(token)
	// handle error
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			logger.Info(errors.Wrap(err, op))
			return nil, domain.ErrInvalidToken
		}
		logger.Error(errors.Wrap(err, op))
		return nil, domain.ErrInternalServer
	}

	return claims, nil
}
