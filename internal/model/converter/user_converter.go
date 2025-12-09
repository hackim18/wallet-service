package converter

import (
	"wallet-service/internal/entity"
	"wallet-service/internal/model"
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

func UserToEvent(user *entity.User) *model.UserEvent {
	return &model.UserEvent{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
