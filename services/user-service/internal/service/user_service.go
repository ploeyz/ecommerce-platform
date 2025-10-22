package service

import (
	"errors"

	"github.com/ploezy/ecommerce-platform/user-service/internal/model"
	"github.com/ploezy/ecommerce-platform/user-service/internal/repository"
	"github.com/ploezy/ecommerce-platform/user-service/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(email, password, firstName, lastName string) (*model.User, error)
	Login(email,password,jwtSecret string)(string,*model.User,error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService{
	return &userService{repo: repo}
}

func (s *userService) Register(email,password,firstName,lastName string) (*model.User,error){
	// Check if user already exists
	existingUser, _ := s.repo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	//Create User
	user := &model.User{
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		Role:      "customer",
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil,err
	}
	return user, nil
}

func (s *userService) Login(email,password,jwtSecret string)(string,*model.User,error){
	// Find user
	user, err := s.repo.FindByEmail(email)
	if err != nil{
		return "", nil, errors.New("invalid email or password")
	}
	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID,user.Email,user.Role,jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil

}