package auth

import (
	"context"
	"time"

	"github.com/anousoneFS/clean-architecture/helper"
	"github.com/anousoneFS/clean-architecture/internal/user"
	"github.com/o1egl/paseto/v2"
)

type AuthUsecase struct {
	userUsecase user.UserUsecase
	pasetoKey   []byte
}

func NewAuthUsecase(uc user.UserUsecase, key []byte) AuthUsecase {
	return AuthUsecase{userUsecase: uc, pasetoKey: key}
}

func (u AuthUsecase) Login(ctx context.Context, req LoginRequest) (res LoginResponse, err error) {
	user, err := u.userUsecase.GetByUsername(ctx, req.Username)
	if err != nil {
		return LoginResponse{}, err
	}
	if err := helper.ComparePassword(req.Password, user.Password); err != nil {
		return LoginResponse{}, err
	}
	return generateToken(u.pasetoKey, user)
}

var now = time.Now

func generateToken(secret []byte, u *user.User) (LoginResponse, error) {
	issAt := now()
	claims := paseto.JSONToken{
		Subject:    u.Name,
		IssuedAt:   issAt,
		Expiration: issAt.Add(5 * time.Hour),
		NotBefore:  issAt,
	}
	userClaims := UserClaims{
		ID: u.Name,
	}
	claims.Set("user", userClaims)
	accessKey, err := paseto.Encrypt(secret, claims, nil)
	if err != nil {
		return LoginResponse{}, err
	}
	// TODO: defines the values used to refresh the token.
	claims.Set("renewable", true)
	claims.Expiration = claims.Expiration.Add(48 * time.Hour)
	refreshKey, err := paseto.Encrypt(secret, claims, nil)
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{
		AccessToken:  accessKey,
		RefreshToken: refreshKey,
	}, nil
}

type UserClaims struct {
	ID       string `json:"id"`
	BranchID string `json:"branchID"`
	VendorID string `json:"vendorID"`
	Role     string `json:"role"`
	RoleID   string `json:"roleID"`
}
