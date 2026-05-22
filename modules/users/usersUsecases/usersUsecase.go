package usersusecases

import (
	"fmt"

	"github.com/WhoaWicked/board-game-shop/config"
	"github.com/WhoaWicked/board-game-shop/modules/users"
	usersrepositories "github.com/WhoaWicked/board-game-shop/modules/users/usersRepositories"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersrepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepository usersrepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}
	result, err := u.usersRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *usersUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}
	result, err := u.usersRepository.InsertUser(req, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}
	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			Role_id:  user.Role_id,
		},
		Token: nil,
	}
	return passport, nil
}
