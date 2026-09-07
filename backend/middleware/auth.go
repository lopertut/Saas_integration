package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend/service"
	"log"
)

type AuthMiddleware struct {
	auth *service.AuthService
}

func NewAuthMiddleware(auth *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		auth: auth,
	}
}

func (m *AuthMiddleware) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		userId, cartId, err := m.auth.DataFromJWT(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("JWT RESULT: userId=%v cartId=%v", userId, cartId)

		ctx := context.WithValue(r.Context(), "user_id", userId)
		ctx = context.WithValue(ctx, "cart_id", cartId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
