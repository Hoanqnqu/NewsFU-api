package inport

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string `json:"id"`
	AuthID    string `json:"auth_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	ImageUrl  string `json:"image_url"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type AdminLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserPayload struct {
	AuthID   string `json:"auth_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	ImageUrl string `json:"image_url"`
}

type UpdateUserPayload struct {
	ID       uuid.UUID
	AuthID   string
	Email    string
	Name     string
	Role     string
	ImageUrl string
}

type Like struct {
	UserId string
	NewsId string
}
type UsersUseCase interface {
	GetAll() ([]*User, error)
	Insert(user *CreateUserPayload) error
	Update(user *UpdateUserPayload) error
	GetAdmin(email string, password string) (user *UpdateUserPayload, err error)
	GetUserByAuthID(authID string) (user *UpdateUserPayload, err error)
	Like(like *Like) error
	DisLike(like *Like) error
	Save(like *Like) error
	GetSavedNews(userID string) ([]*News, error)
	Search(keyword string) ([]*User, error)
	Delete(userID string) error
	View(userID string, newsID string)
}
