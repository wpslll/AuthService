package handlers

import (
	"AuthService/internal/domain"
	"AuthService/internal/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_Register(t *testing.T) {
	tests := []struct {
		Name             string
		ServiceErr       error
		ExpectServiceErr bool
		ExpectHandlerErr bool
		r                http.Request
		User             domain.User
	}{
		{
			Name:             "Success",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: false,
			r: *httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{
			"Username": "qwe",
			"Password": "111"
			}`)),
			User: domain.User{Username: "qwe", Password: "111"},
		},
		{
			Name:             "Internal service error",
			ServiceErr:       errors.New("Test service error"),
			ExpectServiceErr: true,
			ExpectHandlerErr: false,
			r: *httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{
			"Username": "qwe",
			"Password": "111"
			}`)),
			User: domain.User{Username: "qwe", Password: "111"},
		},
		{
			Name:             "User already exists",
			ServiceErr:       domain.ErrUserAlreadyExists,
			ExpectServiceErr: true,
			ExpectHandlerErr: false,
			r: *httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{
			"Username": "qwe",
			"Password": "111"
			}`)),
			User: domain.User{Username: "qwe", Password: "111"},
		},
		{
			Name:             "Empty json",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: true,
			r: *httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{
			"Username": "",
			"Password": ""
			}`)),
			User: domain.User{Username: "qwe", Password: "111"},
		},
		{
			Name:             "Absolutely empty json",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: true,
			r:                *httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(``)),
			User:             domain.User{},
		},
		{
			Name:             "Blank user data",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: true,
			r: *httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{
			"Username": "    ",
			"Password": "    "
			}`)),
			User: domain.User{Username: "    ", Password: "    "},
		},
		{
			Name:             "Border user data",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: false,
			r: *httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{
			"Username": "qwertyuiopasdfghjklz",
			"Password": "qwertyuiop"
			}`)),
			User: domain.User{Username: "qwertyuiopasdfghjklz", Password: "qwertyuiop"},
		},
		{
			Name:             "+1 max length user data",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: true,
			r: *httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{
			"Username": "qwertyuiopasdfghjklzx",
			"Password": "qwertyuiopa"
			}`)),
			User: domain.User{Username: "qwertyuiopasdfghjklzx", Password: "qwertyuiopa"},
		},
		{
			Name:             "-1 max length user data",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: false,
			r: *httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{
			"Username": "qwertyuiopasdfghjkl",
			"Password": "qwertyuio"
			}`)),
			User: domain.User{Username: "qwertyuiopasdfghjkl", Password: "qwertyuio"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			service := mocks.NewMockService(ctrl)
			if !tt.ExpectHandlerErr {
				service.EXPECT().Create(tt.r.Context(), tt.User).Return(tt.ServiceErr)
			}
			handler := NewHandler(service)
			rr := httptest.NewRecorder()
			handler.Register(rr, &tt.r)
			switch {
			case !tt.ExpectServiceErr && !tt.ExpectHandlerErr:
				require.Equal(t, http.StatusCreated, rr.Code)
			case tt.ExpectServiceErr && tt.ServiceErr.Error() == "Test service error":
				require.Equal(t, http.StatusInternalServerError, rr.Code)
				assert.Contains(t, rr.Body.String(), "Test service error")
			case tt.ExpectServiceErr && errors.Is(tt.ServiceErr, domain.ErrUserAlreadyExists):
				require.Equal(t, http.StatusConflict, rr.Code)
				assert.Contains(t, rr.Body.String(), domain.ErrUserAlreadyExists.Error())
			case tt.ExpectHandlerErr:
				require.Equal(t, http.StatusBadRequest, rr.Code)
			}
		})
	}
}

func TestHandler_Auth(t *testing.T) {
	tests := []struct {
		Name             string
		ServiceErr       error
		ExpectServiceErr bool
		ExpectHandlerErr bool
		r                http.Request
		User             domain.User
	}{
		{
			Name:             "Success",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: false,
			r: *httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(`{
			"Username": "qwe",
			"Password": "111"
			}`)),
			User: domain.User{Username: "qwe", Password: "111"},
		},
		{
			Name:             "Internal service error",
			ServiceErr:       errors.New("Test service error"),
			ExpectServiceErr: true,
			ExpectHandlerErr: false,
			r: *httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(`{
			"Username": "qwe",
			"Password": "111"
			}`)),
			User: domain.User{Username: "qwe", Password: "111"},
		},
		{
			Name:             "User not found",
			ServiceErr:       domain.ErrUserNotFound,
			ExpectServiceErr: true,
			ExpectHandlerErr: false,
			r: *httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(`{
			"Username": "qwe",
			"Password": "111"
			}`)),
			User: domain.User{Username: "qwe", Password: "111"},
		},
		{
			Name:             "Empty json",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: true,
			r: *httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(`{
			"Username": "",
			"Password": ""
			}`)),
			User: domain.User{Username: "qwe", Password: "111"},
		},
		{
			Name:             "Absolutely empty json",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: true,
			r:                *httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(``)),
			User:             domain.User{},
		},
		{
			Name:             "Blank user data",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: true,
			r: *httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(`{
			"Username": "    ",
			"Password": "    "
			}`)),
			User: domain.User{Username: "    ", Password: "    "},
		},
		{
			Name:             "Border user data",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: false,
			r: *httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(`{
			"Username": "qwertyuiopasdfghjklz",
			"Password": "qwertyuiop"
			}`)),
			User: domain.User{Username: "qwertyuiopasdfghjklz", Password: "qwertyuiop"},
		},
		{
			Name:             "+1 max length user data",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: true,
			r: *httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(`{
			"Username": "qwertyuiopasdfghjklzx",
			"Password": "qwertyuiopa"
			}`)),
			User: domain.User{Username: "qwertyuiopasdfghjklzx", Password: "qwertyuiopa"},
		},
		{
			Name:             "-1 max length user data",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: false,
			r: *httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(`{
			"Username": "qwertyuiopasdfghjkl",
			"Password": "qwertyuio"
			}`)),
			User: domain.User{Username: "qwertyuiopasdfghjkl", Password: "qwertyuio"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			service := mocks.NewMockService(ctrl)
			if !tt.ExpectHandlerErr {
				service.EXPECT().Auth(tt.r.Context(), tt.User).Return("fake.jwt.token", tt.ServiceErr)
			}
			handler := NewHandler(service)
			rr := httptest.NewRecorder()
			handler.Auth(rr, &tt.r)
			switch {
			case !tt.ExpectServiceErr && !tt.ExpectHandlerErr:
				require.Equal(t, http.StatusOK, rr.Code)
				res := rr.Result()
				cookies := res.Cookies()
				cookie := cookies[0]
				require.Equal(t, "accessToken", cookie.Name)
				require.NotEmpty(t, cookie.Value)
				require.Equal(t, "/", cookie.Path)
				require.True(t, cookie.HttpOnly)
			case tt.ExpectServiceErr && tt.ServiceErr.Error() == "Test service error":
				require.Equal(t, http.StatusInternalServerError, rr.Code)
				assert.Contains(t, rr.Body.String(), "Test service error")
			case tt.ExpectServiceErr && errors.Is(tt.ServiceErr, domain.ErrUserNotFound):
				require.Equal(t, http.StatusNotFound, rr.Code)
				assert.Contains(t, rr.Body.String(), domain.ErrUserNotFound.Error())
			case tt.ExpectHandlerErr:
				require.Equal(t, http.StatusBadRequest, rr.Code)
			}
		})
	}
}

func TestHandler_Validate(t *testing.T) {
	tests := []struct {
		Name             string
		ServiceErr       error
		ExpectServiceErr bool
		ExpectHandlerErr bool
		r                http.Request
		tCookie          http.Cookie
	}{
		{
			Name:             "Success",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: false,
			r:                *httptest.NewRequest(http.MethodPost, "/validate", strings.NewReader(``)),
			tCookie: http.Cookie{Name: "accessToken",
				Value:    "fake.jwt.token",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Path:     "/",
				MaxAge:   15 * 60 * 60},
		},
		{
			Name:             "Internal service error",
			ServiceErr:       errors.New("Test service error"),
			ExpectServiceErr: true,
			ExpectHandlerErr: false,
			r:                *httptest.NewRequest(http.MethodPost, "/validate", strings.NewReader(``)),
			tCookie: http.Cookie{Name: "accessToken",
				Value:    "fake.jwt.token",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Path:     "/",
				MaxAge:   15 * 60 * 60},
		},
		{
			Name:             "No cookie value",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: true,
			r:                *httptest.NewRequest(http.MethodPost, "/validate", strings.NewReader(``)),
			tCookie: http.Cookie{Name: "accessToken",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Path:     "/",
				MaxAge:   15 * 60 * 60},
		},
		{
			Name:             "Empty cookie value",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: true,
			r:                *httptest.NewRequest(http.MethodPost, "/validate", strings.NewReader(``)),
			tCookie: http.Cookie{Name: "accessToken",
				Value:    "",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Path:     "/",
				MaxAge:   15 * 60 * 60},
		},
		{
			Name:             "Absolutely empty cookie",
			ServiceErr:       nil,
			ExpectServiceErr: false,
			ExpectHandlerErr: true,
			r:                *httptest.NewRequest(http.MethodPost, "/validate", strings.NewReader(``)),
			tCookie:          http.Cookie{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			service := mocks.NewMockService(ctrl)
			tt.r.AddCookie(&tt.tCookie)
			if !tt.ExpectHandlerErr {
				service.EXPECT().Validate(tt.tCookie.Value).Return(tt.ServiceErr)
			}
			handler := NewHandler(service)
			rr := httptest.NewRecorder()
			handler.Validate(rr, &tt.r)
			switch {
			case !tt.ExpectServiceErr && !tt.ExpectHandlerErr:
				require.Equal(t, http.StatusOK, rr.Code)
			case tt.ExpectServiceErr && tt.ServiceErr.Error() == "Test service error":
				require.Equal(t, http.StatusUnauthorized, rr.Code)
				assert.Contains(t, rr.Body.String(), "Test service error")
			case tt.ExpectHandlerErr:
				require.Equal(t, http.StatusUnauthorized, rr.Code)
			}
		})
	}
}
