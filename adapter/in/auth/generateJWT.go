package auth

import (
	"fmt"
	inport "news-api/application/port/in"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("my_secret_key")
var tokens []string

type Claims struct {
	ID       string `json:"ID"`
	AuthID   string `json:"auth_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	ImageUrl string `json:"image_url"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *inport.UpdateUserPayload) (string, error) {
	expirationTime := time.Now().Add(20000 * time.Minute)

	claims := &Claims{
		ID:       user.ID.String(),
		AuthID:   user.AuthID,
		Email:    user.Email,
		Name:     user.Name,
		Role:     user.Role,
		ImageUrl: user.ImageUrl,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)

}

func ExtractUser(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")

}
