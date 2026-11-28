package converter

import (
	"github.com/hackim18/wallet-service/internal/entity"
	"github.com/hackim18/wallet-service/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	id := user.ID
	return &model.UserResponse{
		ID:    &id,
		Name:  user.Name,
		Email: user.Email,
	}
}

func UserToTokenResponse(user *entity.User, login *model.UserLogin) *model.UserResponse {
	return &model.UserResponse{
		AccessToken: login.AccessToken,
	}
}
