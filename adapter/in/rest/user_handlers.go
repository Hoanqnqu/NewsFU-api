package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"news-api/adapter/in/auth"
	inport "news-api/application/port/in"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandlers struct {
	userUseCase      inport.UsersUseCase
	recommendUseCase inport.RecommendUseCase
}

func NewUserHandlers(userUseCase inport.UsersUseCase, recommendUseCase inport.RecommendUseCase) *UserHandlers {
	return &UserHandlers{userUseCase: userUseCase, recommendUseCase: recommendUseCase}
}
func (u *UserHandlers) GetSavedNews(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	user := request.Context().Value("user").(inport.UpdateUserPayload)
	newsList, err := u.userUseCase.GetSavedNews(user.ID.String())
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 500,
			Message:    "Unknown err",
		})
		return
	}
	json.NewEncoder(response).Encode(APIResponse[[]*inport.News]{
		StatusCode: 200,
		Message:    "Ok",
		Data:       newsList,
	})
}
func (u *UserHandlers) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var usersList []*inport.User
	var err error
	keywords := request.URL.Query()["keyword"]
	if len(keywords) == 1 {
		usersList, err = u.userUseCase.Search(keywords[0])
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})
			return
		}
	} else {
		usersList, err = u.userUseCase.GetAll()
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})
		}
	}
	json.NewEncoder(response).Encode(APIResponse[[]*inport.User]{
		StatusCode: 200,
		Message:    "Ok",
		Data:       usersList,
	})
}

func (u *UserHandlers) Insert(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user inport.CreateUserPayload
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request"})
		return
	}
	err = u.userUseCase.Insert(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 500,
			Message:    "Unknown err",
		})
		return
	}
	existUser, err := u.userUseCase.GetUserByAuthID(user.AuthID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 500,
			Message:    "Unknown err",
		})
		return
	}
	if err = u.recommendUseCase.InsertUser(context.Background(), existUser.ID); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 500,
			Message:    "Unknown err",
		})
		return
	}

	json.NewEncoder(response).Encode(APIResponse[any]{
		StatusCode: 200,
		Message:    "Ok",
		Data:       user,
	})
}

func (u *UserHandlers) Update(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user inport.UpdateUserPayload
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request",
		})
		return
	}
	id := chi.URLParam(request, "id")
	user.ID, err = uuid.Parse(id)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request"})
		return
	}
	err = u.userUseCase.Update(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 500,
			Message:    "Unknown err",
		})
		return
	}
	json.NewEncoder(response).Encode(APIResponse[any]{
		StatusCode: 200,
		Message:    "Ok",
	})
}

func (u *UserHandlers) Login(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user inport.CreateUserPayload
	err := json.NewDecoder(request.Body).Decode(&user)
	var accessToken string
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request"})
		return
	}
	if existUser, _ := u.userUseCase.GetUserByAuthID(user.AuthID); existUser == nil {
		err = u.userUseCase.Insert(&user)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})
			return
		}
		existUser, err = u.userUseCase.GetUserByAuthID(user.AuthID)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})
			return
		}
		if err = u.recommendUseCase.InsertUser(context.Background(), existUser.ID); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})
			return
		}
		accessToken, err = auth.GenerateJWT(existUser)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})
			return
		}
	} else {
		accessToken, err = auth.GenerateJWT(existUser)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})
			return
		}
	}
	json.NewEncoder(response).Encode(APIResponseLogin{
		StatusCode:  200,
		Message:     "Ok",
		AccessToken: accessToken,
	})
}

func (u *UserHandlers) AdminLogin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user inport.AdminLoginPayload
	err := json.NewDecoder(request.Body).Decode(&user)
	var accessToken string
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request"})
		return
	}
	existUser, err := u.userUseCase.GetAdmin(user.Email, user.Password)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 500,
			Message:    "Unknown err",
		})
		return
	}
	if existUser == nil {
		response.WriteHeader(400)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request",
		})
		return
	}
	accessToken, err = auth.GenerateJWT(existUser)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 500,
			Message:    "Unknown err",
		})
		return
	}

	json.NewEncoder(response).Encode(APIResponseLogin{
		StatusCode:  200,
		Message:     "Ok",
		AccessToken: accessToken,
	})
}

func (u *UserHandlers) Like(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	newsId := chi.URLParam(request, "newsId")
	user := request.Context().Value("user").(inport.UpdateUserPayload)
	err := u.userUseCase.Like(&inport.Like{
		UserId: user.ID.String(),
		NewsId: newsId,
	})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 500,
			Message:    "Unknown err",
		})
		return
	}
	json.NewEncoder(response).Encode(APIResponse[any]{
		StatusCode: 201,
		Message:    "Created",
	})
}

func (u *UserHandlers) Dislike(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	newsId := chi.URLParam(request, "newsId")
	user := request.Context().Value("user").(inport.UpdateUserPayload)
	err := u.userUseCase.DisLike(&inport.Like{
		UserId: user.ID.String(),
		NewsId: newsId,
	})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 500,
			Message:    "Unknown err",
		})
		return
	}
	json.NewEncoder(response).Encode(APIResponse[any]{
		StatusCode: 201,
		Message:    "Ok",
	})
}

func (u *UserHandlers) Save(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	newsId := chi.URLParam(request, "newsId")
	user := request.Context().Value("user").(inport.UpdateUserPayload)
	err := u.userUseCase.Save(&inport.Like{
		UserId: user.ID.String(),
		NewsId: newsId,
	})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 500,
			Message:    "Unknown err",
		})
		return
	}
	json.NewEncoder(response).Encode(APIResponse[any]{
		StatusCode: 201,
		Message:    "Ok",
	})
}

func (u *UserHandlers) Delete(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(request, "id")
	err := u.userUseCase.Delete(userID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 500,
			Message:    "Unknown err",
		})
		return
	}
	json.NewEncoder(response).Encode(APIResponse[any]{
		StatusCode: 201,
		Message:    "Ok",
	})

}

func (u *UserHandlers) View(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	newsID := chi.URLParam(request, "newsID")
	tokenString := request.Header.Get("Authorization")
	var userID string
	if tokenString != "" {
		tokenString = tokenString[len("Bearer "):]
		if tokenString != "" {
			claim, err := auth.ExtractUser(tokenString)
			if err == nil {
				userID = claim["ID"].(string)
			}
		}
	}
	if userID == "" {
		return
	}
	u.userUseCase.View(userID, newsID)
	json.NewEncoder(response).Encode(APIResponse[any]{
		StatusCode: 201,
		Message:    "Ok",
	})
}
