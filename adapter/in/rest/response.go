package rest

type APIResponse[T any] struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}
type APIResponseLogin struct {
	StatusCode  int    `json:"status"`
	Message     string `json:"message"`
	AccessToken string `json:"token"`
}
