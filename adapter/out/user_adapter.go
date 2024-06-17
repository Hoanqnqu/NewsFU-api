package outAdapter

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	outport "news-api/application/port/out"
	db "news-api/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserAdapter struct {
	pool *pgxpool.Pool
}

func NewUserAdapter(pool *pgxpool.Pool) *UserAdapter {
	return &UserAdapter{pool: pool}
}

func (u *UserAdapter) GetAll() ([]db.User, error) {
	query := db.New(u.pool)

	return query.GetAllUsers(context.Background())

}

func (u *UserAdapter) Insert(user outport.User) error {
	query := db.New(u.pool)
	err := query.InsertUser(context.Background(), db.InsertUserParams{
		ID: pgtype.UUID{
			Bytes: user.ID,
			Valid: true,
		},
		AuthID: user.AuthID,
		Email: pgtype.Text{
			String: user.Email,
			Valid:  true,
		},
		Name: pgtype.Text{
			String: user.Name,
			Valid:  true,
		},
		Role: pgtype.Text{
			String: "user",
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: user.ImageUrl,
			Valid:  true,
		},
	})
	return err
}

func (u *UserAdapter) View(userID string, newsID string) {
	query := db.New(u.pool)
	query.InsertView(context.Background(), db.InsertViewParams{
		NewsID: pgtype.UUID{
			Bytes: uuid.MustParse(newsID),
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: uuid.MustParse(userID),
			Valid: true,
		},
	})
}

func (u *UserAdapter) Update(user outport.User) error {
	query := db.New(u.pool)
	err := query.UpdateUser(context.Background(), db.UpdateUserParams{
		Name: pgtype.Text{
			String: user.Name,
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: user.ImageUrl,
			Valid:  true,
		},
		ID: pgtype.UUID{
			Bytes: user.ID,
			Valid: true,
		},
	})
	return err
}

func (u *UserAdapter) GetByAuthID(authID string) (outUser outport.User, err error) {
	query := db.New(u.pool)
	dbUser, err := query.GetUserByAuthID(context.Background(), authID)
	if err != nil {
		return outport.User{}, err
	}
	outUser = outport.User{
		ID:        dbUser.ID.Bytes,
		AuthID:    dbUser.AuthID,
		Email:     dbUser.Email.String,
		CreatedAt: dbUser.CreatedAt.Time,
		Name:      dbUser.Name.String,
		Role:      dbUser.Role.String,
		ImageUrl:  dbUser.ImageUrl.String,
	}
	return outUser, nil
}
func (u *UserAdapter) GetAdmin(email string, password string) (user *outport.User, err error) {
	query := db.New(u.pool)
	dbUsers, err := query.GetAdmin(context.Background(), db.GetAdminParams{
		Email: pgtype.Text{
			String: email,
			Valid:  true,
		},
		Password: pgtype.Text{
			String: password,
			Valid:  true,
		},
	})
	if len(dbUsers) == 0 {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	user = &outport.User{
		ID:        dbUsers[0].ID.Bytes,
		AuthID:    dbUsers[0].AuthID,
		Email:     dbUsers[0].Email.String,
		CreatedAt: dbUsers[0].CreatedAt.Time,
		Name:      dbUsers[0].Name.String,
		Role:      dbUsers[0].Role.String,
		ImageUrl:  dbUsers[0].ImageUrl.String,
	}
	return user, nil
}

func (u *UserAdapter) Like(like outport.Like) error {
	query := db.New(u.pool)
	err := query.InsertLike(context.Background(), db.InsertLikeParams{
		NewsID: pgtype.UUID{
			Bytes: uuid.MustParse(like.NewsID),
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: uuid.MustParse(like.UserID),
			Valid: true,
		},
	})
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			err = query.DeleteLike(context.Background(), db.DeleteLikeParams{
				NewsID: pgtype.UUID{
					Bytes: uuid.MustParse(like.NewsID),
					Valid: true,
				},
				UserID: pgtype.UUID{
					Bytes: uuid.MustParse(like.UserID),
					Valid: true,
				},
			})
		}
	} else {
		query.DeleteDisLike(context.Background(), db.DeleteDisLikeParams{
			NewsID: pgtype.UUID{
				Bytes: uuid.MustParse(like.NewsID),
				Valid: true,
			},
			UserID: pgtype.UUID{
				Bytes: uuid.MustParse(like.UserID),
				Valid: true,
			},
		})
	}
	return err
}

func (u *UserAdapter) DisLike(like outport.Like) error {
	query := db.New(u.pool)
	err := query.InsertDisLike(context.Background(), db.InsertDisLikeParams{
		NewsID: pgtype.UUID{
			Bytes: uuid.MustParse(like.NewsID),
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: uuid.MustParse(like.UserID),
			Valid: true,
		},
	})
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			err = query.DeleteDisLike(context.Background(), db.DeleteDisLikeParams{
				NewsID: pgtype.UUID{
					Bytes: uuid.MustParse(like.NewsID),
					Valid: true,
				},
				UserID: pgtype.UUID{
					Bytes: uuid.MustParse(like.UserID),
					Valid: true,
				},
			})
		}
	} else {
		query.DeleteLike(context.Background(), db.DeleteLikeParams{
			NewsID: pgtype.UUID{
				Bytes: uuid.MustParse(like.NewsID),
				Valid: true,
			},
			UserID: pgtype.UUID{
				Bytes: uuid.MustParse(like.UserID),
				Valid: true,
			},
		})
	}
	return err
}

func (u *UserAdapter) Save(like outport.Like) error {
	query := db.New(u.pool)
	err := query.InsertSave(context.Background(), db.InsertSaveParams{
		NewsID: pgtype.UUID{
			Bytes: uuid.MustParse(like.NewsID),
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: uuid.MustParse(like.UserID),
			Valid: true,
		},
	})
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			err = query.DeleteSave(context.Background(), db.DeleteSaveParams{
				NewsID: pgtype.UUID{
					Bytes: uuid.MustParse(like.NewsID),
					Valid: true,
				},
				UserID: pgtype.UUID{
					Bytes: uuid.MustParse(like.UserID),
					Valid: true,
				},
			})
		}
	}
	return err
}

func (u *UserAdapter) GetSavedNews(userID string) ([]outport.NewsWithCategory, error) {
	query := db.New(u.pool)
	response, err := query.GetSaves(context.Background(), pgtype.UUID{
		Bytes: uuid.MustParse(userID),
		Valid: true,
	})
	if err != nil {
		return nil, err
	}
	sl := make([]outport.NewsWithCategory, len(response))
	if err != nil {
		return nil, err
	}
	var categoryIds []pgtype.UUID
	for i, v := range response {
		sl[i].Author = v.Author
		sl[i].Content = v.Content
		sl[i].Description = v.Description
		sl[i].Title = v.Title
		sl[i].Url = v.Url
		sl[i].ImageUrl = v.ImageUrl
		sl[i].PublishAt = v.PublishAt
		sl[i].ID = v.ID
		err = json.Unmarshal(v.CategoryIds, &categoryIds)
		if err != nil {
			return nil, err
		}
		sl[i].Categories = categoryIds
	}
	return sl, nil

}

func (u *UserAdapter) Search(keyword string) ([]db.User, error) {
	query := db.New(u.pool)

	return query.SearchUsers(context.Background(), pgtype.Text{
		String: keyword,
		Valid:  true,
	})
}

func (u *UserAdapter) Delete(id string) error {
	query := db.New(u.pool)
	return query.DeleteUser(context.Background(), pgtype.UUID{
		Bytes: uuid.MustParse(id),
		Valid: true,
	})
}
