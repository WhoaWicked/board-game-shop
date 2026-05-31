package usersusecases

import (
	"fmt"

	"github.com/WhoaWicked/board-game-shop/config"
	"github.com/WhoaWicked/board-game-shop/modules/users"
	usersrepositories "github.com/WhoaWicked/board-game-shop/modules/users/usersRepositories"
	"github.com/WhoaWicked/board-game-shop/pkg/shopauth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
	GetUserProfile(userId string) (*users.User, error)
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

	userClaims := &users.UserClaims{
		Id:      user.Id,
		Role_id: user.Role_id,
	}

	accessToken, err := shopauth.NewShopAuth(shopauth.Access, u.cfg.Jwt(), userClaims)
	refreshToken, err := shopauth.NewShopAuth(shopauth.Refresh, u.cfg.Jwt(), userClaims)

	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			Role_id:  user.Role_id,
		},
		Token: &users.UserToken{
			// Id:           user.Id, ในฟังก์ชัน InsertOauth จะกำหนดให้เอง
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}
	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil
}

func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	// Parse token
	claims, err := shopauth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// find oauth
	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// get profile
	profile, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	userClaims := &users.UserClaims{
		Id:      profile.Id,
		Role_id: profile.Role_id,
	}

	// new access token
	accessToken, err := shopauth.NewShopAuth(
		shopauth.Access,
		u.cfg.Jwt(),
		userClaims,
	)
	if err != nil {
		return nil, err
	}

	// new refresh token
	refreshToken := shopauth.RepeatToken(
		u.cfg.Jwt(),
		userClaims,
		claims.ExpiresAt.Unix(),
	)

	// set passport
	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *usersUsecase) DeleteOauth(oauthId string) error {
	if err := u.usersRepository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil
}

func (h *usersUsecase) GetUserProfile(userId string) (*users.User, error) {
	profile, err := h.usersRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
