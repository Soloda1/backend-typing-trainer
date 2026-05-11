package users

import (
	"backend-typing-trainer/internal/domain/models"
	"time"
)

type userResponse struct {
	ID        string    `json:"id"`
	Login     string    `json:"login"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type userSingleResponse struct {
	User userResponse `json:"user"`
}

type userListResponse struct {
	Users []userResponse `json:"users"`
}

func fromModel(u *models.User) userResponse {
	return userResponse{
		ID:        u.ID.String(),
		Login:     u.Login,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt,
	}
}
