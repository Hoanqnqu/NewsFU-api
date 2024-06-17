package outport

import (
	"news-api/internal/db"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	AuthID    string
	Email     string
	Name      string
	Role      string
	ImageUrl  string
	CreatedAt time.Time
}

type Like struct {
	UserID string
	NewsID string
}

type UsersPort interface {
	GetAll() ([]db.User, error)
	GetAdmin(email string, password string) (user *User, err error)
	Insert(user User) error
	Update(user User) error
	GetByAuthID(authID string) (user User, err error)
	Like(like Like) error
	DisLike(like Like) error
	Save(like Like) error
	GetSavedNews(userID string) ([]NewsWithCategory, error)
	Search(keyword string) ([]db.User, error)
	Delete(id string) error
	View(userID string, newsID string)
}
