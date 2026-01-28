package services

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/llorenzinho/goauth/internal"
	"github.com/llorenzinho/goauth/internal/database"
	"github.com/llorenzinho/goauth/internal/rest/dtos"
	"github.com/llorenzinho/goauth/pkg/log"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(dtos.CreateUserParams) (database.User, error)
	GetUserByEmail(string) (database.User, error)
	GetUserByID(id string) (database.User, error)
	VerifyUserEmail(id string) error
	UpdateUserPassword(id string, params dtos.UpdateUserPasswordHashParams) error
	UpdateEmail(id string, params dtos.UpdateUserEmailParams) error
}

type userService struct {
	db *pgxpool.Pool
	q  *database.Queries
	l  *zap.Logger
}

func NewUserService(db *pgxpool.Pool) UserService {
	return &userService{
		db: db,
		q:  database.New(),
		l:  log.Get(),
	}
}

func (s *userService) CreateUser(arg dtos.CreateUserParams) (database.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	s.l.Info("Creating user", zap.String("email", arg.Email))

	err := validator.New().Struct(arg)
	if err != nil {
		s.l.Error("Failed to validate user creation parameters", zap.Error(err))
		return database.User{}, err
	}

	return s.q.CreateUser(ctx, s.db, database.CreateUserParams{
		Email:        arg.Email,
		PasswordHash: arg.Password,
	})
}

func (s *userService) GetUserByEmail(email string) (database.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	user, err := s.q.GetUserByEmail(ctx, s.db, email)
	if err != nil {
		s.l.Error("Failed to get user by email", zap.String("email", email), zap.Error(err))
		return database.User{}, err
	}

	return user, nil
}

func (s *userService) GetUserByID(id string) (database.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	parsed, err := uuid.Parse(id)
	if err != nil {
		s.l.Error("Failed to parse uuid string", zap.String("id", id), zap.Error(err))
		return database.User{}, internal.ErrInvalidUUID
	}

	user, err := s.q.GetUserByID(ctx, s.db, parsed)
	if err != nil {
		s.l.Error("Failed to get user by ID", zap.String("id", id), zap.Error(err))
		return database.User{}, err
	}

	return user, nil
}

func (s *userService) VerifyUserEmail(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	parsed, err := uuid.Parse(id)
	if err != nil {
		s.l.Error("Failed to parse uuid string", zap.String("id", id), zap.Error(err))
		return internal.ErrInvalidUUID
	}

	_, err = s.q.VerifyUserEmail(ctx, s.db, parsed)
	if err != nil {
		s.l.Error("Failed to verify user email", zap.String("id", id), zap.Error(err))
		return err
	}

	return nil
}

func (s *userService) UpdateUserPassword(id string, params dtos.UpdateUserPasswordHashParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	s.l.Info("Updating user password", zap.String("id", id))

	err := validator.New().Struct(params)
	if err != nil {
		s.l.Error("Failed to validate update user password parameters", zap.String("id", id), zap.Error(err))
		return err
	}

	parsed, err := uuid.Parse(id)
	if err != nil {
		s.l.Error("Failed to parse uuid string", zap.String("id", id), zap.Error(err))
		return internal.ErrInvalidUUID
	}

	_, err = s.q.UpdateUserPasswordHash(ctx, s.db, database.UpdateUserPasswordHashParams{
		ID:           parsed,
		PasswordHash: params.Password,
	})
	if err != nil {
		s.l.Error("Failed to update user password", zap.String("id", id), zap.Error(err))
		return err
	}

	return nil
}

func (s *userService) UpdateEmail(id string, params dtos.UpdateUserEmailParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	s.l.Info("Updating user email", zap.String("id", id), zap.String("newEmail", params.Email))

	err := validator.New().Struct(params)
	if err != nil {
		s.l.Error("Failed to validate update user email parameters", zap.String("id", id), zap.Error(err))
		return err
	}

	parsed, err := uuid.Parse(id)
	if err != nil {
		s.l.Error("Failed to parse uuid string", zap.String("id", id), zap.Error(err))
		return err
	}

	_, err = s.q.UpdateUserEmail(ctx, s.db, database.UpdateUserEmailParams{
		ID:    parsed,
		Email: params.Email,
	})
	if err != nil {
		s.l.Error("Failed to update user email", zap.String("id", id), zap.Error(err))
		return err
	}

	return nil
}
