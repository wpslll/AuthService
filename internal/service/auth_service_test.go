package service

import (
	"AuthService/internal/domain"
	"AuthService/internal/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_Create(t *testing.T) {
	tests := []struct {
		Name      string
		RepoErr   error
		ExpectErr bool
		Ctx       context.Context
		User      domain.User
	}{
		{
			Name:      "Success",
			RepoErr:   nil,
			ExpectErr: false,
			Ctx:       context.Background(),
			User:      domain.User{Username: "qwe", Password: "111"},
		},
		{
			Name:      "Repository error",
			RepoErr:   errors.New("Repository error"),
			ExpectErr: true,
			Ctx:       context.Background(),
			User:      domain.User{Username: "qwe", Password: "111"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mocks.NewMockRepository(ctrl)
			repo.EXPECT().Create(tt.Ctx, tt.User).Return(tt.RepoErr)
			service := NewService(repo)
			err := service.Create(tt.Ctx, tt.User)
			if !tt.ExpectErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestService_Auth(t *testing.T) {
	t.Setenv("SECRET", "test-secret")
	tests := []struct {
		Name      string
		RepoErr   error
		ExpectErr bool
		Ctx       context.Context
		User      domain.User
	}{
		{
			Name:      "Success",
			RepoErr:   nil,
			ExpectErr: false,
			Ctx:       context.Background(),
			User:      domain.User{Username: "qwe", Password: "111"},
		},
		{
			Name:      "Repository error",
			RepoErr:   errors.New("Repository error"),
			ExpectErr: true,
			Ctx:       context.Background(),
			User:      domain.User{Username: "qwe", Password: "111"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mocks.NewMockRepository(ctrl)
			repo.EXPECT().Auth(tt.Ctx, tt.User).Return(tt.RepoErr)
			service := NewService(repo)
			token, err := service.Auth(tt.Ctx, tt.User)
			if !tt.ExpectErr {
				require.NotEmpty(t, token)
				require.NoError(t, err)
			} else {
				require.Empty(t, token)
				require.Error(t, err)
			}
		})
	}
}

func TestService_Validate(t *testing.T) {
	t.Setenv("SECRET", "test-secret")
	tests := []struct {
		Name      string
		RepoErr   error
		ExpectErr bool
		Ctx       context.Context
		User      domain.User
	}{
		{
			Name:      "Success",
			RepoErr:   nil,
			ExpectErr: false,
			Ctx:       context.Background(),
			User:      domain.User{Username: "qwe", Password: "111"},
		},
		{
			Name:      "Validate error",
			RepoErr:   nil,
			ExpectErr: true,
			Ctx:       context.Background(),
			User:      domain.User{Username: "qwe", Password: "111"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mocks.NewMockRepository(ctrl)
			repo.EXPECT().Auth(tt.Ctx, tt.User).Return(tt.RepoErr)
			service := NewService(repo)
			token, err := service.Auth(tt.Ctx, tt.User)
			require.NotEmpty(t, token)
			require.NoError(t, err)
			if !tt.ExpectErr {
				err = service.Validate(token)
				require.NoError(t, err)
			} else {
				err = service.Validate("")
				require.Error(t, err)
				err = service.Validate("qwertyuiop")
				require.Error(t, err)
			}
		})
	}
}
