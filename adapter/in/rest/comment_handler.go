package rest

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	inport "news-api/application/port/in"
)

type CommentHandler struct {
	commentUseCase inport.CommentUseCase
}

func NewCommentHandler(comment inport.CommentUseCase) *CommentHandler {
	return &CommentHandler{commentUseCase: comment}
}

type InsertCommentPayload struct {
	NewsID  uuid.UUID `json:"news_id"`
	Comment string    `json:"comment"`
}

func (c *CommentHandler) Insert(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var comment InsertCommentPayload
	err := json.NewDecoder(request.Body).Decode(&comment)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request"})
		return
	}
	user := request.Context().Value("user").(inport.UpdateUserPayload)

	err = c.commentUseCase.InsertComment(comment.NewsID, user.ID, comment.Comment)
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

func (c *CommentHandler) GetCommentsByNews(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	_newsID := chi.URLParam(request, "newsID")
	newsID, err := uuid.Parse(_newsID)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request"})
		return
	}
	comments, err := c.commentUseCase.GetCommentsByNewsID(newsID)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Bad request"})
		return
	}
	json.NewEncoder(response).Encode(APIResponse[[]*inport.Comment]{
		StatusCode: 200,
		Message:    "Ok",
		Data:       comments,
	})

}
