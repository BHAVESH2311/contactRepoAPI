package repository

import (
	"contactapp/errors"

	"github.com/jinzhu/gorm"
)

type Repository interface{
	
	GetAll(uow *UnitOfWork,out interface{},total *int,limit int,offset int,queryprocessors ...QueryProcessor)error
	Add(uow *UnitOfWork,out interface{}) error
	Update(uow *UnitOfWork,out interface{}) error
	UpdateWithMap(uow *UnitOfWork, model interface{}, value map[string]interface{},queryprocessors ...QueryProcessor) error
	GetRecordForUser(uow *UnitOfWork,userID int,out interface{},queryprocessors ...QueryProcessor)error
	Save(uow *UnitOfWork,value interface{})error
	GetRecordForContact(uow *UnitOfWork,userID int,out interface{},queryprocessors ...QueryProcessor) error
	GetAllContacts(uow *UnitOfWork,out interface{},total *int,limit int,offset int,queryprocessors ...QueryProcessor) error
	GetAllContactsDetails(uow *UnitOfWork,out interface{},limit int,offset int,queryProcessors ...QueryProcessor)error
	GetRecordForContactDetails(uow *UnitOfWork,cdID uint,out interface{},queryProcesors ...QueryProcessor)error
	GetRecordForComponent(uow *UnitOfWork, userID uint, out interface{}, queryProcessors ...QueryProcessor) error
}

type GormRepository struct{}

func NewGormRepository() *GormRepository{
	return &GormRepository{}
}

type UnitOfWork struct{
	DB *gorm.DB
	Committed bool
	Readonly bool
}

func NewUnitOfWork(db *gorm.DB,readonly bool) *UnitOfWork{
	commit := false
	if readonly {

		return &UnitOfWork{

		DB: db.New(),
		Committed : commit,
		Readonly : readonly,

		}
		
	}

	return &UnitOfWork{
		DB: db.New().Begin(),
		Committed: commit,
		Readonly: readonly,
	}
}


func (uow *UnitOfWork) Commit(){

	if !uow.Readonly && !uow.Committed{
		uow.Committed=true
		uow.DB.Commit()
	}

}

func (uow *UnitOfWork) Rollback(){

	if !uow.Readonly && !uow.Committed{
		uow.DB.Rollback()
	}
}


func (repository *GormRepository) Save(uow *UnitOfWork, value interface{}) error{
	return uow.DB.Debug().Save(value).Error
}

func Filter(condition string, args ...interface{}) QueryProcessor{
	return func(db *gorm.DB,out interface{})(*gorm.DB,error){
		db=db.Debug().Where(condition,args...)
		return db,nil
	}
}

func Select(query interface{},args ...interface{})QueryProcessor{
	return func(db *gorm.DB,out interface{})(*gorm.DB,error){
		db=db.Select(query,args...)
		return db,nil
	}
}


func Preload(matchedAssociations []string) QueryProcessor{
	return func(db *gorm.DB,out interface{})(*gorm.DB,error){
		for _,association:=range matchedAssociations{
			db.Debug().Preload(association)
		}
		return db,nil
	}
}




func (repository *GormRepository) GetAllContacts(uow *UnitOfWork, out interface{}, total *int, limit int, offset int, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}

	db = db.Offset(offset)
	return db.Debug().Limit(limit).Find(out).Count(total).Error
}


func (repository *GormRepository) GetRecordForComponent(uow *UnitOfWork, userID uint, out interface{}, queryProcessors ...QueryProcessor) error {
	queryProcessors = append([]QueryProcessor{Filter("id = ?", userID)}, queryProcessors...)
	return repository.GetRecord(uow, out, queryProcessors...)
}


func (repository *GormRepository) GetAll(uow *UnitOfWork,out interface{},total *int,limit int,offset int,queryProcessors ...QueryProcessor)error{
	db:=uow.DB
	db,err:=executeQueryProcessors(db,out,queryProcessors...)
	if err!=nil{
		return err
	}
	db=db.Offset(offset)
	return db.Debug().Limit(limit).Find(out).Count(total).Error
}


func (repository *GormRepository) GetAllContactsDetails(uow *UnitOfWork,out interface{},limit int,offset int,queryProcessors ...QueryProcessor)error{
	db:=uow.DB
	db,err:=executeQueryProcessors(db,out,queryProcessors...)
	if err!=nil{
		return err
	}

	db=db.Offset(offset)

	return db.Debug().Limit(limit).Find(out).Error


}

func executeQueryProcessors(db *gorm.DB,out interface{},queryProcessors ...QueryProcessor)(*gorm.DB,error){


	var err error
	for _,query := range queryProcessors{
		if query != nil{
			db,err = query(db,out)
			if err!=nil{
				return db,err
			}
		}
		
	}
	return db,nil
}

// Add adds record to table.
func (repository *GormRepository) Add(uow *UnitOfWork, out interface{}) error {
	return uow.DB.Create(out).Error
}

func (repository *GormRepository) Update(uow *UnitOfWork, out interface{}) error {
	return uow.DB.Model(out).Update(out).Error
}

func DoesRecordExistForUser(db *gorm.DB,userID uint,out interface{},queryProcessors ...QueryProcessor) (bool,error){
	if userID==0{
		return false,errors.NewValidationError("DoesRecordExistForUser : Invalid user ID")
	}
	count:=0
	db,err:=executeQueryProcessors(db,out,queryProcessors...)
	if err!=nil{
		return false,err
	}
	if err:=db.Debug().Model(out).Where("id = ?",userID).Count(&count).Error;err!=nil{
		return false,err
	}
	if(count>0){
		return true,nil
	}
	return false,nil
}

func DoesRecordExistForContact(db *gorm.DB,conID uint,out interface{},queryProcessors ...QueryProcessor)(bool,error){
	if(conID==0){
		return false,errors.NewValidationError("DoesRecordExistForContact:Invalid Contact ID")
	}
	count:=0
	db,err:=executeQueryProcessors(db,out,queryProcessors...)
	if err!=nil{
		return false,err
	}
	if err:=db.Debug().Model(out).Where("id = ?",conID).Count(&count).Error;err!=nil{
		return false,err
	}
	if(count>0){
		return true,nil
	}
	return false,nil
}

func DoesRecordExistForContactDetails(db *gorm.DB,conID int,out interface{},queryProcessors ...QueryProcessor)(bool,error){
	if(conID==0){
		return false,errors.NewValidationError("DoesRecordExistForContactDetail:Invalid contact ID")
	}
	db,err:=executeQueryProcessors(db,out,queryProcessors...)
	if err!=nil{
		return false,err
	}
	count:=0
	if err:=db.Debug().Model(out).Where("id=?",conID).Error;err!=nil{
		return false,err
	}
	if(count>0){
		return true,nil
	}
	return false,nil

}

func (repository *GormRepository)GetRecordForUser(uow *UnitOfWork,userID int,out interface{},queryProcessors ...QueryProcessor)error{
	queryProcessors=append([]QueryProcessor{Filter("id=?",userID)},queryProcessors...)
	return repository.GetRecord(uow,out,queryProcessors...)
}

func(repository *GormRepository)GetRecordForContact(uow *UnitOfWork,conID int,out interface{},queryProcessors ...QueryProcessor)error{
	queryProcessors=append([]QueryProcessor{Filter("id=?",conID)},queryProcessors...)
	return repository.GetRecord(uow,out,queryProcessors...)
}

func(repository *GormRepository)GetRecordForContactDetails(uow *UnitOfWork,conID uint,out interface{},queryProcessors ...QueryProcessor)error{
	queryProcessors=append([]QueryProcessor{Filter("id=?",conID)},queryProcessors...)
	return repository.GetRecord(uow,out,queryProcessors...)
}


func (repository *GormRepository) GetRecord(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().First(out).Error
}
func(repository *GormRepository)UpdateWithMap(uow *UnitOfWork,model interface{},value map[string]interface{},queryProcessors ...QueryProcessor)error{
	db:=uow.DB
	db,err:=executeQueryProcessors(db,value,queryProcessors...)
	if(err!=nil){
		return err
	}
	return db.Debug().Model(model).Update(value).Error
}