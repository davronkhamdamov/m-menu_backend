package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/davronkhamdamov/restaraunt_backend/models"
	"github.com/davronkhamdamov/restaraunt_backend/utils"
	"github.com/golang-jwt/jwt/v5"
)

const your_secret_key string = "your_secret_key"

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type contextKey string

const (
	UserIDKey = contextKey("user_id")
	RoleKey   = contextKey("role")
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header is missing", nil)
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(your_secret_key), nil
		})

		if err != nil || !token.Valid {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token", err.Error())
			return
		}
		user := models.User{}
		if dbResult := models.DB.Where("ID = ?", claims.UserID).First(&user); dbResult.Error != nil {
			utils.RespondWithError(w, http.StatusNotFound, "User not found", dbResult.Error)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func CorsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
