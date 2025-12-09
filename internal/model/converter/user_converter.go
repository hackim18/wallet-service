package converter

import (
	"wallet-service/internal/entity"
	"wallet-service/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}
}

func UserToTokenResponse(user *entity.User, login *model.UserLogin) *model.UserResponse {
	return &model.UserResponse{
		Token: login.AccessToken,
	}
}

func UserToEvent(user *entity.User) *model.UserEvent {
	return &model.UserEvent{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
