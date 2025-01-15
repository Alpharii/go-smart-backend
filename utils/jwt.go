package utils

import (
	"time"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GeneateToken(id uint, role string) (string, error){
	claims := jwt.MapClaims{
		"user_id": id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"role": role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}