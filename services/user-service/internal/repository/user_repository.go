package repository

import (
	"errors"

	"github.com/ploezy/ecommerce-platform/user-service/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindbyId(id uint) (*model.User,error)
}

type userRepository  struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository{
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error{
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error){
	var user model.User
	err := r.db.Where("email = ?",email).First(&user).Error
	if err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil,errors.New("user not found")
		}
		return nil,err
	}
	return &user,nil
}

func (r * userRepository) FindbyId(id uint) (*model.User, error){

	var user model.User
	err := r.db.First(&user,id).Error
	if err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, errors.New("user not found")
		}
		return nil,err
	}
	return &user,nil
}
