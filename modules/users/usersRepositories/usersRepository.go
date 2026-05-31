package usersrepositories

import (
	"context"
	"fmt"
	"time"

	"github.com/WhoaWicked/board-game-shop/modules/users"
	userspatterns "github.com/WhoaWicked/board-game-shop/modules/users/usersPatterns"
	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertOauth(req *users.UserPassport) error
	FindOneOauth(refreshToken string) (*users.Oauth, error)
	UpdateOauth(req *users.UserToken) error
	GetProfile(userId string) (*users.User, error)
	DeleteOauth(oauthId string) error
}

type usersRepository struct {
	db *sqlx.DB
}

func UsersRepository(db *sqlx.DB) IUsersRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	result := userspatterns.InsertUser(req, r.db, isAdmin)
	var err error
	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}
	user, err := result.Result()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *usersRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
	SELECT id, email, password, username, role_id
	FROM users
	WHERE email = $1;
	`
	user := new(users.UserCredentialCheck)
	if err := r.db.Get(user, query, email); err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (r *usersRepository) InsertOauth(req *users.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	query := `INSERT INTO oauth (user_id, access_token, refresh_token)
	VALUES ($1, $2, $3) RETURNING id;`
	if err := r.db.QueryRowContext(
		ctx,
		query,
		req.User.Id,
		req.Token.AccessToken,
		req.Token.RefreshToken,
	).Scan(&req.Token.Id); err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) FindOneOauth(refreshToken string) (*users.Oauth, error) {
	query := `SELECT id, user_id FROM oauth WHERE refresh_token = $1;`
	oauth := new(users.Oauth)
	if err := r.db.Get(oauth, query, refreshToken); err != nil {
		return nil, fmt.Errorf("oauth not found")
	}
	return oauth, nil
}

func (r *usersRepository) UpdateOauth(req *users.UserToken) error {
	query := `UPDATE oauth SET
	access_token = :access_token,
	refresh_token = :refresh_token
	WHERE id = :id;`
	if _, err := r.db.NamedExecContext(context.Background(), query, req); err != nil {
		return fmt.Errorf("update oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) GetProfile(userId string) (*users.User, error) {
	query := `SELECT id, username, email, role_id FROM users WHERE id = $1;`
	profile := new(users.User)
	if err := r.db.Get(profile, query, userId); err != nil {
		return nil, fmt.Errorf("get user profile failed: %v", err)
	}
	return profile, nil
}

func (r *usersRepository) DeleteOauth(oauthId string) error {
	query := `DELETE FROM oauth WHERE id = $1;`
	result, err := r.db.ExecContext(context.Background(), query, oauthId)
	if err != nil {
		return fmt.Errorf("delete oauth failed: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected failed: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("oauth not found")
	}
	return nil
}
