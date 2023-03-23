package middlewares

import (
	"assignment2/api/responses"
	"assignment2/api/utils"
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func SetContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func AuthJwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{"status": "failed", "message": "Missing authorization token"}
		header := r.Header.Get("Authorization")

		if header == "" {
			responses.JSON(rw, http.StatusForbidden, resp)
			return
		}

		token, err := jwt.Parse(header, func(t *jwt.Token) (interface{}, error) {
			return utils.GetSecretKey(), nil
		})

		if err != nil {
			resp["message"] = "Invalid token"
			responses.JSON(rw, http.StatusForbidden, resp)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		userID := claims[utils.KEY_USER_ID]
		ctx := context.WithValue(r.Context(), utils.KEY_USER_ID, userID)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
