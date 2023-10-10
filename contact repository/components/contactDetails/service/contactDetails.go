package service

import (
	"github.com/jinzhu/gorm"
	"contactapp/repository"
	contactDetails "contactapp/models/contactDetails"
	"contactapp/utils"
	"contactapp/errors"
	"time"

)

type ContactDetailService struct {
	db *gorm.DB
	repository repository.Repository
	associations []string
}

func NewContactDetailService(db *gorm.DB,repository repository.Repository) *ContactDetailService {
	return &ContactDetailService{
		db:db,
		repository: repository,
		associations: []string{"Contact"},
	}
}

func (service *ContactDetailService) CreateContactDetail(newCd *contactDetails.ContactDetails) error {
	//  Creating unit of work.
	uow := repository.NewUnitOfWork(service.db, false)

	defer uow.Rollback()

	err := service.repository.Add(uow, newCd)
	if err != nil {
		uow.Rollback()
		return err
	}

	uow.Commit()
	return nil
}


func (s *ContactDetailService) GetAllContactDetails(userId, contactId uint, includes []string, limit, offset int) ([]contactDetails.ContactDetails, int, error) {
	uow := repository.NewUnitOfWork(s.db, true)

	defer uow.Rollback()

	matchedAssociations := utils.GetAssociations(includes, s.associations)

	contactsInfo := []contactDetails.ContactDetails{}
	var total int

	err := s.repository.GetAll(uow, &contactsInfo, &total, limit, offset,
		repository.Filter("`user_refer` = ? AND `contact_refer` = ?", userId, contactId),
		repository.Preload(matchedAssociations))

	if err != nil {
		return nil, -1, err
	}

	uow.Commit()
	return contactsInfo, total, nil
}


func (service *ContactDetailService) doesContactDetailExist(ID uint) error {
	exists, err := repository.DoesRecordExistForContactDetails(service.db, int(ID), contactDetails.ContactDetails{},
		repository.Filter("`id` = ?", ID))
	if !exists || err != nil {
		return errors.NewValidationError("Contact detail ID is Invalid")
	}
	return nil
}


func (service *ContactDetailService) GetContactDetails(userId, contactId, contactInfoId uint, includes []string) (interface{}, error) {
	err := service.doesContactDetailExist(contactInfoId)

	if err != nil {
		return nil, err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	defer uow.Rollback()

	matchedAssociations := utils.GetAssociations(includes, service.associations)

	contactInfo := contactDetails.ContactDetails{}

	err = service.repository.GetRecordForComponent(uow, contactInfoId, contactInfo,
		repository.Filter("`user_refer` = ? AND `contact_refer` = ?", userId, contactId),
		repository.Preload(matchedAssociations))

	if err != nil {
		return nil, err
	}

	uow.Commit()
	return contactInfo, nil
}


func (service *ContactDetailService) DeleteContactDetails(cdToDelete *contactDetails.ContactDetails) error {
	err := service.doesContactDetailExist(cdToDelete.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.Rollback()

	if err := service.repository.UpdateWithMap(uow, cdToDelete, map[string]interface{}{

		"DeletedAt": time.Now(),
	},
		repository.Filter("`id`=?", cdToDelete.ID)); err != nil {
		uow.Rollback()
		return err
	}
	uow.Commit()
	return nil
}


func (service *ContactDetailService) UpdateContactDetails(cdToUpdate *contactDetails.ContactDetails) error {
	err := service.doesContactDetailExist(cdToUpdate.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	defer uow.Rollback()

	tempCd := contactDetails.ContactDetails{}
	err = service.repository.GetRecordForComponent(uow, cdToUpdate.ID, &tempCd,
		repository.Filter("`user_refer` = ? AND `contact_refer` = ?", cdToUpdate.UserRefer, cdToUpdate.ContactRefer))
	if err != nil {
		return err
	}
	cdToUpdate.CreatedAt = tempCd.CreatedAt

	err = service.repository.Save(uow, cdToUpdate)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}
















