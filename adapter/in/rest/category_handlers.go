package rest

import (
	"encoding/json"
	"net/http"
	inport "news-api/application/port/in"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CategoryHandlers struct {
	categoryUseCase inport.CategoriesUseCase
}

func NewCategoryHandlers(categoryUseCase inport.CategoriesUseCase) *CategoryHandlers {
	return &CategoryHandlers{categoryUseCase: categoryUseCase}
}

func (u *CategoryHandlers) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var categoriesList []*inport.Category
	var err error
	keywords := request.URL.Query()["keyword"]
	if len(keywords) == 1 {
		categoriesList, err = u.categoryUseCase.Search(keywords[0])
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})

		}
	} else {
		categoriesList, err = u.categoryUseCase.GetAll()
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(APIResponse[any]{
				StatusCode: 500,
				Message:    "Unknown err",
			})
		}
	}
	json.NewEncoder(response).Encode(APIResponse[[]*inport.Category]{
		StatusCode: 200,
		Message:    "Ok",
		Data:       categoriesList,
	})
}

func (u *CategoryHandlers) Insert(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var category inport.CreateCategoryPayload
	err := json.NewDecoder(request.Body).Decode(&category)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request"})
		return
	}
	err = u.categoryUseCase.Insert(&category)

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

func (u *CategoryHandlers) Update(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var category inport.UpdateCategoryPayload
	err := json.NewDecoder(request.Body).Decode(&category)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request"})
		return
	}
	id := chi.URLParam(request, "id")
	category.ID, err = uuid.Parse(id)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request"})
		return
	}

	err = u.categoryUseCase.Update(&category)
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

func (u *CategoryHandlers) Delete(response http.ResponseWriter, request *http.Request) {
	return

}
