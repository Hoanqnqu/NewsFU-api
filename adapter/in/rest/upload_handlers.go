package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	inport "news-api/application/port/in"

	"github.com/google/uuid"
)

type UploadHandlers struct {
	uploadUseCases inport.UploadUseCase
}

func NewUploadHandlers(uploadUseCases inport.UploadUseCase) *UploadHandlers {
	return &UploadHandlers{uploadUseCases}
}

func (u *UploadHandlers) Upload(response http.ResponseWriter, request *http.Request) {
	var url string
	response.Header().Set("Content-Type", "application/json")
	file, fileHeader, err := request.FormFile("image")
	if err != nil {
		response.WriteHeader(400)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 400,
			Message:    "Invalid payload",
		})
		return
	}
	defer file.Close()
	objectKey := fmt.Sprintf("readFU/%s/%s", uuid.New().String(), fileHeader.Filename)
	url, err = u.uploadUseCases.PutObject(context.Background(), objectKey, file)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(APIResponse[any]{
			StatusCode: 500,
			Message:    "Error during upload file",
		})
		return
	}
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(APIResponse[any]{
		StatusCode: 201,
		Message:    "Uploaded successfully",
		Data:       url,
	})
}
