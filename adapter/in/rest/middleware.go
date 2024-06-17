package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"news-api/adapter/in/auth"
	inport "news-api/application/port/in"

	"github.com/google/uuid"
)

type User struct {
	ID   string
	Role string
}

func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println(request.URL.Path)
		handler.ServeHTTP(writer, request)
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tokenString := request.Header.Get("Authorization")
		if tokenString == "" {
			writer.WriteHeader(401)
			json.NewEncoder(writer).Encode(APIResponse[any]{
				StatusCode: 401,
			})
			return
		}

		tokenString = tokenString[len("Bearer "):]
		claim, err := auth.ExtractUser(tokenString)

		if err != nil {
			writer.WriteHeader(401)
			json.NewEncoder(writer).Encode(APIResponse[any]{
				StatusCode: 401,
			})
			return
		}

		if claim["role"].(string) != "admin" {
			writer.WriteHeader(401)
			json.NewEncoder(writer).Encode(APIResponse[any]{
				StatusCode: 401,
			})
			return
		}

		ctx := context.WithValue(request.Context(), "user", inport.CreateUserPayload{
			AuthID:   claim["auth_id"].(string),
			Role:     claim["role"].(string),
			Email:    claim["email"].(string),
			Name:     claim["name"].(string),
			ImageUrl: claim["image_url"].(string),
		})

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tokenString := request.Header.Get("Authorization")
		if tokenString == "" {
			writer.WriteHeader(401)
			json.NewEncoder(writer).Encode(APIResponse[any]{
				StatusCode: 401,
			})
			return
		}
		tokenString = tokenString[len("Bearer "):]
		if tokenString == "" {
			writer.WriteHeader(401)
			json.NewEncoder(writer).Encode(APIResponse[any]{
				StatusCode: 401,
			})
			return
		}
		claim, err := auth.ExtractUser(tokenString)

		if err != nil {
			writer.WriteHeader(401)
			json.NewEncoder(writer).Encode(APIResponse[any]{
				StatusCode: 401,
			})
			return
		}

		if claim["role"].(string) != "user" {
			writer.WriteHeader(http.StatusForbidden)
			json.NewEncoder(writer).Encode(APIResponse[any]{
				StatusCode: 401,
			})
			return
		}
		ctx := context.WithValue(request.Context(), "user", inport.UpdateUserPayload{
			ID:       uuid.MustParse(claim["ID"].(string)),
			AuthID:   claim["auth_id"].(string),
			Role:     claim["role"].(string),
			Email:    claim["email"].(string),
			Name:     claim["name"].(string),
			ImageUrl: claim["image_url"].(string),
		})
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
