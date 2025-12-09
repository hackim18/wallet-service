package usecase

import (
	"context"
	"net/http"
	"wallet-service/internal/constants"
	"wallet-service/internal/entity"
	"wallet-service/internal/model"
	"wallet-service/internal/model/converter"
	"wallet-service/internal/repository"
	"wallet-service/internal/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	JWT            *utils.JWTHelper
	UserRepository *repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, jwt *utils.JWTHelper,
	userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		JWT:            jwt,
		UserRepository: userRepository,
	}
}

func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if request.Token == "" {
		return nil, model.ErrUnauthorized
	}

	tokenID, err := uuid.Parse(request.Token)
	if err != nil {
		c.Log.Warnf("Failed parse token to UUID : %+v", err)
		return nil, model.ErrUnauthorized
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, tokenID); err != nil {
		c.Log.Warnf("Failed find user by token : %+v", err)
		return nil, model.ErrUnauthorized
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, model.ErrInternalServerError
	}

	return &model.Auth{ID: user.ID}, nil
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	total, err := c.UserRepository.CountByCondition(tx, "email = ?", request.Email)
	if err != nil {
		c.Log.Warnf("Failed to check existing user : %+v", err)
		return nil, utils.Error(constants.ErrCheckUser, http.StatusInternalServerError, err)
	}

	if total > 0 {
		return nil, utils.Error(constants.ErrUserAlreadyExists, http.StatusConflict, nil)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
		return nil, utils.Error(constants.ErrProcessPassword, http.StatusInternalServerError, err)
	}

	user := &entity.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(password),
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed to insert user : %+v", err)
		return nil, utils.Error(constants.ErrCreateUser, http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, utils.Error(constants.ErrCommitTransaction, http.StatusInternalServerError, err)
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindByCondition(tx, user, "email = ?", request.Email); err != nil {
		c.Log.Warnf("Failed to find user by email : %+v", err)
		return nil, utils.Error(constants.ErrInvalidEmailOrPassword, http.StatusUnauthorized, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		c.Log.Warnf("Invalid password : %+v", err)
		return nil, utils.Error(constants.ErrInvalidEmailOrPassword, http.StatusUnauthorized, err)
	}

	if c.JWT == nil {
		c.Log.Warn("JWT helper not configured")
		return nil, utils.Error(constants.ErrGenerateAccessToken, http.StatusInternalServerError, nil)
	}

	accessToken, err := c.JWT.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		c.Log.Warnf("Failed to generate access token : %+v", err)
		return nil, utils.Error(constants.ErrGenerateAccessToken, http.StatusInternalServerError, err)
	}

	login := &model.UserLogin{
		AccessToken: accessToken,
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, model.ErrInternalServerError
	}

	return converter.UserToTokenResponse(user, login), nil
}
