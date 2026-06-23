package service

import (
	"AuthService/internal/domain"
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Repository interface {
	Create(context.Context, domain.User) error
	Auth(context.Context, domain.User) error
}

type AuthService struct {
	repo Repository
}

func NewService(repo Repository) AuthService {
	return AuthService{repo: repo}
}

func (a *AuthService) Create(ctx context.Context, user domain.User) error {
	err := a.repo.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) Auth(ctx context.Context, user domain.User) (string, error) {
	if err := a.repo.Auth(ctx, user); err != nil {
		return "", err
	}
	secretString := []byte(os.Getenv("SECRET"))
	accesToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256, 
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	signedToken, err := accesToken.SignedString(secretString)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (a *AuthService) Validate(tokenString string) error {
	token, err := jwt.Parse(tokenString, func (token *jwt.Token) (any, error)  {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Wrong signing method")
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("Token isn't valid")
	}
	return nil
}