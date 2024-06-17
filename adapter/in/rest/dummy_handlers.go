package rest

import (
	"encoding/json"
	"net/http"
	inport "news-api/application/port/in"
)

type DummyHandler struct {
	dummyUseCase inport.Dummy
}

type Test struct {
	Name string
	Age  int
}

func NewDummyHandler(dummyUseCase inport.Dummy) *DummyHandler {
	return &DummyHandler{dummyUseCase: dummyUseCase}
}

func (d *DummyHandler) Dummy(response http.ResponseWriter, request *http.Request) {
	//dummyOutput := d.dummyUseCase.Get()
	response.WriteHeader(404)
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(APIResponse[any]{
		StatusCode: 404,
		Message:    "Not Found",
	})
}
