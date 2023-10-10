package service

import (
	"contactapp/models/contact"
	"contactapp/repository"
	"contactapp/utils"
	"contactapp/errors"
	"github.com/jinzhu/gorm"
	"time"
)

type ContactService struct {
	db *gorm.DB
	repository repository.Repository
	associations []string
}


func NewContactService(db *gorm.DB, repository repository.Repository) *ContactService {
	
	return &ContactService{db: db, repository: repository,associations: []string{"User"}}
}

func (s *ContactService) CreateContact(newContact *contact.Contact) error {

	uow:=repository.NewUnitOfWork(s.db,false)
	defer uow.Rollback()
	err:=s.repository.Add(uow,newContact)
	if(err!=nil){
		uow.Rollback()
		return err
	}
	uow.Commit()
	return nil
}

func (service *ContactService) GetAllContacts(userId uint, limit, offset int) ([]contact.Contact, int, error) {
	uow := repository.NewUnitOfWork(service.db, true)

	defer uow.Rollback()

	contacts := []contact.Contact{}
	var total int

	err:=service.repository.GetAllContacts(uow,&contacts,&total,limit,offset,repository.Filter("user_id=?",userId))

	if err != nil {
		return nil, -1, err
	}

	uow.Commit()
	return contacts, total, nil
}

func (service *ContactService) doesContactExist(ID uint) error {
	exists, err := repository.DoesRecordExistForContact(service.db, ID, contact.Contact{},
		repository.Filter("`id` = ?", ID))
	if !exists || err != nil {
		return errors.NewValidationError("Contact ID is Invalid")
	}
	return nil
}

func (service *ContactService) GetContact(userId uint, contactId int, includes []string) (interface{}, error) {
	err := service.doesContactExist(uint(contactId))

	if err != nil {
		return nil, err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	defer uow.Rollback()
	matchedAssociations := utils.GetAssociations(includes, service.associations)

	contact := contact.Contact{}

	err = service.repository.GetRecordForContact(uow, contactId, &contact,
		repository.Filter("`user_id` = ?", userId),
		repository.Preload(matchedAssociations))

	if err != nil {
		return nil, err
	}

	uow.Commit()
	return contact, nil
}

func (s *ContactService) UpdateContact(con *contact.Contact)error{
	err:=s.doesContactExist(con.ID)
	if err!=nil{
		return err
	}
	uow:=repository.NewUnitOfWork(s.db,false)

	defer uow.Rollback()
	tempContact:=contact.Contact{}

	err=s.repository.GetRecordForContact(uow,int(con.ID),tempContact,repository.Select("created_at"),repository.Filter("user_id=?",con.UserId))
	
	if err!=nil{
		return err
	}

	con.CreatedAt=tempContact.CreatedAt

	err = s.repository.Save(uow,con)
	if err!=nil{
		return err
	}
	uow.Commit()
	return nil
}

func(s *ContactService) DeleteContact(con *contact.Contact)error{
	err:=s.doesContactExist(con.ID)
	if err!=nil{
		return err
	}
	uow:=repository.NewUnitOfWork(s.db,false)

	defer uow.Rollback()
	if err:=s.repository.UpdateWithMap(uow,con,map[string]interface{}{
		"DeletedAt":time.Now(),
	},repository.Filter("id=?",con.ID));err!=nil{
		uow.Rollback()
		return err
	}
	uow.Commit()
	return nil
}

