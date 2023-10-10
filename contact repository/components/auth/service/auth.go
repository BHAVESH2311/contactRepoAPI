package service

import (
	"contactapp/components/user/service"
	"contactapp/models/user"
	"contactapp/repository"

	"github.com/jinzhu/gorm"
)

type AuthService struct {
	db           *gorm.DB
	repository   repository.Repository
	userservice  *service.UserService
	associations []string
}

func NewAuthService(db *gorm.DB, repo repository.Repository) *AuthService {
	return &AuthService{
		db:           db,
		repository:   repo,
		userservice:  service.NewUserService(db, repo),
		associations: []string{},
	}
}

func (service *AuthService) Register(newUser *user.User) error {
	uow := repository.NewUnitOfWork(service.db, false)

	defer uow.Rollback()

	err := service.userservice.CreateUser(newUser)

	if err != nil {
		uow.Rollback()
		return err
	}

	uow.Commit()
	return nil
}



func (service *AuthService) Login(user *user.User) error {
	uow := repository.NewUnitOfWork(service.db, false)

	defer uow.Rollback()
	err := service.userservice.GetUserByEmail(user)

	if err != nil {
		uow.Rollback()
		return err
	}

	uow.Commit()
	return nil
}

