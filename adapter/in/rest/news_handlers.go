package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"news-api/adapter/in/auth"
	inport "news-api/application/port/in"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type NewsHandlers struct {
	newsUseCase inport.NewsUseCase
}

func NewNewsHandlers(newsUseCase inport.NewsUseCase) *NewsHandlers {
	return &NewsHandlers{newsUseCase: newsUseCase}
}

func (u *NewsHandlers) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var newsList []*inport.News
	var err error
	keywords := request.URL.Query()["keyword"]
	if len(keywords) == 1 {
		fmt.Println("search keyword:", keywords[0])
		newsList, err = u.newsUseCase.SearchNews(keywords[0])
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})
			return
		}
	} else {
		fmt.Println("Get All")
		newsList, err = u.newsUseCase.GetAll()
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})
			return
		}
	}

	json.NewEncoder(response).Encode(APIResponse[[]*inport.News]{
		StatusCode: 200,
		Message:    "Ok",
		Data:       newsList,
	})
}

func (u *NewsHandlers) Insert(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var news inport.CreateNewsPayload
	err := json.NewDecoder(request.Body).Decode(&news)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request"})
		return
	}
	err = u.newsUseCase.Insert(&news)
	fmt.Println(err)
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

func (u *NewsHandlers) Update(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var news inport.UpdateNewsPayload
	err := json.NewDecoder(request.Body).Decode(&news)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request",
		})
		return
	}
	id := chi.URLParam(request, "id")
	news.ID, err = uuid.Parse(id)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request"})
		return
	}
	err = u.newsUseCase.Update(&news)
	fmt.Println(err)
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
func (u *NewsHandlers) GetNewsByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	newsId := chi.URLParam(request, "newsId")
	userId := uuid.Nil.String()
	tokenString := request.Header.Get("Authorization")
	if tokenString != "" {
		tokenString = tokenString[len("Bearer "):]
		if tokenString != "" {
			claim, err := auth.ExtractUser(tokenString)
			fmt.Println(err)
			if err == nil {
				userId = claim["ID"].(string)
			}
		}
	}

	news, err := u.newsUseCase.GetNewsByID(newsId, userId)
	fmt.Println(err)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 404,
			Message:    "Not Found"})
		return
	}

	json.NewEncoder(response).Encode(APIResponse[*inport.News]{
		StatusCode: 200,
		Message:    "Ok",
		Data:       news,
	})
}
func (u *NewsHandlers) GetLatest(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var err error
	count := 10
	offset := 0
	countQueryParams := request.URL.Query()["count"]
	if len(countQueryParams) == 1 {
		_count, err := strconv.Atoi(countQueryParams[0])
		if err == nil {
			count = _count
		}
	}
	offsetQueryParam := request.URL.Query()["offset"]
	if len(offsetQueryParam) == 1 {
		_offset, err := strconv.Atoi(offsetQueryParam[0])
		if err == nil {
			offset = _offset
		}
	}
	newsList, err := u.newsUseCase.GetLatestNews(count, offset)
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
func (u *NewsHandlers) GetPopular(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var err error
	count := 10
	offset := 0
	var categoryID string
	categoryQueryParams := request.URL.Query()["category"]
	if len(categoryQueryParams) == 1 {
		categoryID = categoryQueryParams[0]
	}
	countQueryParams := request.URL.Query()["count"]
	if len(countQueryParams) == 1 {
		_count, err := strconv.Atoi(countQueryParams[0])
		if err == nil {
			count = _count
		}
	}
	offsetQueryParam := request.URL.Query()["offset"]
	if len(offsetQueryParam) == 1 {
		_offset, err := strconv.Atoi(offsetQueryParam[0])
		if err == nil {
			offset = _offset
		}
	}
	newsList, err := u.newsUseCase.GetPopular(categoryID, count, offset)
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
func (u *NewsHandlers) GetRecommend(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var newsList []*inport.News
	var err error
	count := 10
	offset := 0
	var userID string
	tokenString := request.Header.Get("Authorization")
	if tokenString == "" {
		claim, err := auth.ExtractUser(tokenString)
		if err == nil {
			userID = claim["ID"].(string)
		}
	} else {
		tokenString = tokenString[len("Bearer "):]

		if tokenString != "" {
			claim, err := auth.ExtractUser(tokenString)
			if err == nil {
				userID = claim["ID"].(string)
			}
		}
	}

	countQueryParams := request.URL.Query()["count"]
	if len(countQueryParams) == 1 {
		_count, err := strconv.Atoi(countQueryParams[0])
		if err == nil {
			count = _count
		}
	}
	offsetQueryParam := request.URL.Query()["offset"]
	if len(offsetQueryParam) == 1 {
		_offset, err := strconv.Atoi(offsetQueryParam[0])
		if err == nil {
			offset = _offset
		}
	}
	if userID != "" {
		newsList, err = u.newsUseCase.GetRecommend(userID, count, offset)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})
			return
		}
	} else {
		newsList, err = u.newsUseCase.GetPopular("", count, offset)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})
			return
		}
	}
	json.NewEncoder(response).Encode(APIResponse[[]*inport.News]{
		StatusCode: 200,
		Message:    "Ok",
		Data:       newsList,
	})
}

func (u *NewsHandlers) Delete(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	newsID := chi.URLParam(request, "id")
	err := u.newsUseCase.Delete(newsID)
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
