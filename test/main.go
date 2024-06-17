package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/zhenghaoz/gorse/client"
	"log"
	"net/http"
	"news-api/adapter/in/rest"
	outAdapter "news-api/adapter/out"
	"news-api/application/domain/service"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		"postgres",
		"password",
		"localhost",
		"5432",
		"postgres",
	)
	gorse := client.NewGorseClient("http://127.0.0.1:8087", "")
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Fatalln("Can not connect to sql")
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalln("Can not connect to sql")
	}
	defer pool.Close()
	fmt.Println(os.Getenv("AWS_ACCESS_KEY"))
	//init adapter
	s3Adapter := outAdapter.NewS3Adapter(os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_BUCKET_REGION"),
		os.Getenv("AWS_BUCKET_NAME"))
	dummyAdapter := outAdapter.NewDummyAdapter(pool)
	userAdapter := outAdapter.NewUserAdapter(pool)
	categoryAdapter := outAdapter.NewCategoryAdapter(pool)
	newsAdapter := outAdapter.NewNewsAdapter(pool)
	gorseAdapter := outAdapter.NewGorseAdapter(gorse)
	commentAdapter := outAdapter.NewCommentAdapter(pool)
	//init Use case
	s3UseCase := service.NewUploadService(s3Adapter)
	dummyUseCase := service.NewDummyService(dummyAdapter)
	userUseCase := service.NewUsersService(userAdapter, gorseAdapter)
	categoryUseCase := service.NewCategoriesService(categoryAdapter)
	newsUseCase := service.NewNewsService(newsAdapter, gorseAdapter)
	recommendUseCase := service.NewRecommendService(gorseAdapter)
	commentUseCase := service.NewCommentService(commentAdapter)
	//init handler
	dummyHandler := rest.NewDummyHandler(dummyUseCase)
	userHandler := rest.NewUserHandlers(userUseCase, recommendUseCase)
	categoryHandler := rest.NewCategoryHandlers(categoryUseCase)
	newsHandler := rest.NewNewsHandlers(newsUseCase)
	s3Handler := rest.NewUploadHandlers(s3UseCase)
	commentHandler := rest.NewCommentHandler(commentUseCase)

	router := rest.AppRouter(dummyHandler, userHandler, categoryHandler, newsHandler, s3Handler, commentHandler)
	http.ListenAndServe(":3000", router)
}
