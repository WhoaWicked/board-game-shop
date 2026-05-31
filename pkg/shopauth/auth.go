package shopauth

import (
	"errors"
	"fmt"
	"time"

	"github.com/WhoaWicked/board-game-shop/config"
	"github.com/WhoaWicked/board-game-shop/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
)

type shopAuth struct {
	mapClaims *shopMapClaims
	cfg       config.IJwtConfig
}

type shopMapClaims struct {
	Claims *users.UserClaims `json:"Claims"`
	jwt.RegisteredClaims
}

type shopAdmin struct {
	*shopAuth
}

type IShopAuth interface {
	SignToken() string
}

// type IShopAdmin interface {
// 	SignToken() string
// }

func (a *shopAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func (a *shopAdmin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(t) * time.Second))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*shopMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &shopMapClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signin method is invalid")
		}
		return cfg.SecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}
	if claims, ok := token.Claims.(*shopMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid: %v", err)
	}
}

func ParseAdmin(cfg config.IJwtConfig, tokenString string) (*shopMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &shopMapClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signin method is invalid")
		}
		return cfg.AdminKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*shopMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid: %v", err)
	}

}

func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &shopAuth{
		cfg: cfg,
		mapClaims: &shopMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "boardgameshop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "staff", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func NewShopAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IShopAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil

	default:
		return nil, fmt.Errorf("unknow token type")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IShopAuth {
	return &shopAuth{
		cfg: cfg,
		mapClaims: &shopMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "boardgameshop-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "staff", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IShopAuth {
	return &shopAuth{
		cfg: cfg,
		mapClaims: &shopMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "boardgameshop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "staff", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newAdminToken(cfg config.IJwtConfig) IShopAuth {
	return &shopAdmin{
		shopAuth: &shopAuth{
			cfg: cfg,
			mapClaims: &shopMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "boardgameshop-api",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
