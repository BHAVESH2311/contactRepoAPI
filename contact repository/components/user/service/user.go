package service

import (
	"contactapp/errors"
	"contactapp/models/user"
	"contactapp/repository"

	"time"

	"github.com/jinzhu/gorm"
)


type UserService struct{
	db *gorm.DB
	repository repository.Repository
	associations []string
}


func NewUserService(db *gorm.DB, repo repository.Repository)*UserService{
	return &UserService{
		db: db, 
		repository: repo,
		associations: []string{},
	
	}
}

func (s *UserService) CreateUser(newUser *user.User)error{

	uow:=repository.NewUnitOfWork(s.db,false)

	defer uow.Rollback()

	err:=s.repository.Add(uow,newUser)
	if err!=nil{
		uow.Rollback()
		return err
	}
	uow.Commit()
	return nil
}

func(s *UserService) CreateAdmin(newUser *user.User)error{
	uow:=repository.NewUnitOfWork(s.db,false)
	defer uow.Rollback()
	err:=s.repository.Add(uow,newUser)
	if err!=nil{
		uow.Rollback()
		return err
	}
	uow.Commit()
	return nil
}


func (s *UserService) GetAllUsers(limit, offset int) ([]user.User, int, error) {
	uow := repository.NewUnitOfWork(s.db, true)
	defer uow.Rollback()

	var users = []user.User{}
	var total int

	err := s.repository.GetAll(uow,&users,&total,limit,offset)
	if err != nil {
		uow.Rollback()
		return nil, -1, err
	}
	uow.Commit()
	return users, total, nil
}



func (s *UserService) doesUserExist(Id uint)error{
	exists,err :=repository.DoesRecordExistForUser(s.db,Id,user.User{},repository.Filter("id=?",Id))
	if !exists || err!=nil{
		return errors.NewValidationError("Invalid User Id")
	}
	return nil
}


func(s *UserService)GetUserById(userId uint)(interface{},error){
	err:=s.doesUserExist(userId)
	if err!=nil{
		return nil,err
	}

	uow:=repository.NewUnitOfWork(s.db,true)
	defer uow.Rollback()
	var user = []user.User{}
	errors:=s.repository.GetRecordForUser(uow,int(userId),&user)
	if errors!=nil{
		return nil,err
	}
	uow.Commit()
	return user,nil
}

func (s *UserService) GetUserByEmail(user *user.User)error{
	uow:=repository.NewUnitOfWork(s.db,true)

	defer uow.Rollback()

	res:=uow.DB.Where("Email = ?",user.Email).First(&user)
	if res.Error!=nil{
		return res.Error
	}
	uow.Commit()
	return nil
}

func (service *UserService) UpdateUser(userToUpdate *user.User) error {
	err := service.doesUserExist(userToUpdate.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	defer uow.Rollback()

	tempUser := user.User{}

	err = service.repository.GetRecordForUser(uow, int(userToUpdate.ID), &tempUser, repository.Select("`created_at`"),
		repository.Filter("`id` = ?", userToUpdate.ID))
	if err != nil {
		return err
	}

	userToUpdate.UpdatedAt = tempUser.CreatedAt

	err = service.repository.Save(uow, userToUpdate)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}


func (s *UserService) DeleteUser(userToDelete *user.User) error{
	err :=s.doesUserExist(userToDelete.ID)
	if err!=nil{
		return err
	}

	uow := repository.NewUnitOfWork(s.db,false)

	defer uow.Rollback()
	if err:=s.repository.UpdateWithMap(uow,userToDelete,map[string]interface{}{
		"DeletedAt":time.Now(),
	},repository.Filter("`id`=?",userToDelete.ID));err!=nil{
		uow.Rollback()
		return err
	}
	uow.Commit()
	return nil
}




