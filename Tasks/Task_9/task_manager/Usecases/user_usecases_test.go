package Usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"task_manager/Domain"
	"task_manager/Mocks"
	"task_manager/Usecases"
)

// Small hand-written mocks for the two service interfaces.
type passwordMock struct{ mock.Mock }

func (m *passwordMock) Hash(p string) (string, error) {
	a := m.Called(p)
	return a.String(0), a.Error(1)
}
func (m *passwordMock) Compare(hash, p string) bool {
	return m.Called(hash, p).Bool(0)
}

type tokenMock struct{ mock.Mock }

func (m *tokenMock) Generate(u Domain.User) (string, error) {
	a := m.Called(u)
	return a.String(0), a.Error(1)
}

type UserUsecaseSuite struct {
	suite.Suite
	repo      *Mocks.UserRepositoryMock
	passwords *passwordMock
	tokens    *tokenMock
	uc        Usecases.UserUsecase
	ctx       context.Context
}

func (s *UserUsecaseSuite) SetupTest() {
	s.repo = new(Mocks.UserRepositoryMock)
	s.passwords = new(passwordMock)
	s.tokens = new(tokenMock)
	s.uc = Usecases.NewUserUsecase(s.repo, s.passwords, s.tokens)
	s.ctx = context.Background()
}

func (s *UserUsecaseSuite) assertAll() {
	s.repo.AssertExpectations(s.T())
	s.passwords.AssertExpectations(s.T())
	s.tokens.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestRegister() {
	dbErr := errors.New("db down")
	hashErr := errors.New("hash failed")

	// Helper: the happy-path chain, parameterised by the user count.
	happy := func(count int64, role string) func() {
		return func() {
			s.repo.On("GetByUsername", s.ctx, "bob").Return(Domain.User{}, Domain.ErrNotFound).Once()
			s.passwords.On("Hash", "secret1").Return("hashed", nil).Once()
			s.repo.On("Count", s.ctx).Return(count, nil).Once()
			s.repo.On("Create", s.ctx, Domain.User{Username: "bob", Password: "hashed", Role: role}).
				Return(Domain.User{ID: "u1", Username: "bob", Password: "hashed", Role: role}, nil).Once()
		}
	}

	tests := []struct {
		name     string
		username string
		password string
		setup    func()
		wantErr  error
		wantRole string
	}{
		{"first user becomes admin", "bob", "secret1", happy(0, Domain.RoleAdmin), nil, Domain.RoleAdmin},
		{"later user is a normal user", "bob", "secret1", happy(3, Domain.RoleUser), nil, Domain.RoleUser},
		{"username too short", "bo", "secret1", func() {}, Domain.ErrInvalidInput, ""},
		{"password too short", "bob", "123", func() {}, Domain.ErrInvalidInput, ""},
		{"duplicate username", "bob", "secret1", func() {
			s.repo.On("GetByUsername", s.ctx, "bob").Return(Domain.User{Username: "bob"}, nil).Once()
		}, Domain.ErrUsernameTaken, ""},
		{"unexpected repo error on lookup", "bob", "secret1", func() {
			s.repo.On("GetByUsername", s.ctx, "bob").Return(Domain.User{}, dbErr).Once()
		}, dbErr, ""},
		{"hashing fails", "bob", "secret1", func() {
			s.repo.On("GetByUsername", s.ctx, "bob").Return(Domain.User{}, Domain.ErrNotFound).Once()
			s.passwords.On("Hash", "secret1").Return("", hashErr).Once()
		}, hashErr, ""},
		{"create fails", "bob", "secret1", func() {
			s.repo.On("GetByUsername", s.ctx, "bob").Return(Domain.User{}, Domain.ErrNotFound).Once()
			s.passwords.On("Hash", "secret1").Return("hashed", nil).Once()
			s.repo.On("Count", s.ctx).Return(int64(1), nil).Once()
			s.repo.On("Create", s.ctx, Domain.User{Username: "bob", Password: "hashed", Role: Domain.RoleUser}).
				Return(Domain.User{}, dbErr).Once()
		}, dbErr, ""},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			got, err := s.uc.Register(s.ctx, tc.username, tc.password)

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
			} else {
				s.NoError(err)
				s.Equal("bob", got.Username)
				s.EqualValues(tc.wantRole, got.Role)
				s.Empty(got.Password, "the hash must never be returned")
			}
			s.assertAll()
		})
	}
}

func (s *UserUsecaseSuite) TestLogin() {
	dbErr := errors.New("db down")
	tokenErr := errors.New("sign failed")
	stored := Domain.User{ID: "u1", Username: "bob", Password: "hashed", Role: Domain.RoleUser}

	tests := []struct {
		name      string
		setup     func()
		wantToken string
		wantErr   error
	}{
		{"success", func() {
			s.repo.On("GetByUsername", s.ctx, "bob").Return(stored, nil).Once()
			s.passwords.On("Compare", "hashed", "secret1").Return(true).Once()
			s.tokens.On("Generate", stored).Return("jwt-token", nil).Once()
		}, "jwt-token", nil},
		{"unknown user gives the generic credentials error", func() {
			s.repo.On("GetByUsername", s.ctx, "bob").Return(Domain.User{}, Domain.ErrNotFound).Once()
		}, "", Domain.ErrInvalidCredentials},
		{"wrong password gives the same error, no token issued", func() {
			s.repo.On("GetByUsername", s.ctx, "bob").Return(stored, nil).Once()
			s.passwords.On("Compare", "hashed", "secret1").Return(false).Once()
		}, "", Domain.ErrInvalidCredentials},
		{"repository error is propagated", func() {
			s.repo.On("GetByUsername", s.ctx, "bob").Return(Domain.User{}, dbErr).Once()
		}, "", dbErr},
		{"token generation error is propagated", func() {
			s.repo.On("GetByUsername", s.ctx, "bob").Return(stored, nil).Once()
			s.passwords.On("Compare", "hashed", "secret1").Return(true).Once()
			s.tokens.On("Generate", stored).Return("", tokenErr).Once()
		}, "", tokenErr},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			token, err := s.uc.Login(s.ctx, "bob", "secret1")

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
				s.Empty(token)
			} else {
				s.NoError(err)
				s.Equal(tc.wantToken, token)
			}
			s.assertAll()
		})
	}
}

func (s *UserUsecaseSuite) TestPromote() {
	dbErr := errors.New("update failed")

	tests := []struct {
		name    string
		setup   func()
		wantErr error
	}{
		{"success", func() {
			s.repo.On("GetByID", s.ctx, "u1").Return(Domain.User{ID: "u1"}, nil).Once()
			s.repo.On("UpdateRole", s.ctx, "u1", Domain.RoleAdmin).Return(nil).Once()
		}, nil},
		{"user not found, UpdateRole never called", func() {
			s.repo.On("GetByID", s.ctx, "u1").Return(Domain.User{}, Domain.ErrNotFound).Once()
		}, Domain.ErrNotFound},
		{"UpdateRole error is propagated", func() {
			s.repo.On("GetByID", s.ctx, "u1").Return(Domain.User{ID: "u1"}, nil).Once()
			s.repo.On("UpdateRole", s.ctx, "u1", Domain.RoleAdmin).Return(dbErr).Once()
		}, dbErr},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			err := s.uc.Promote(s.ctx, "u1")

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
			} else {
				s.NoError(err)
			}
			s.assertAll()
		})
	}
}

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuite))
}
