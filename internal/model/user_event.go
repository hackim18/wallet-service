package model

import (
	"time"

	"github.com/google/uuid"
)

type UserEvent struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (u *UserEvent) GetId() uuid.UUID {
	return u.ID
}
