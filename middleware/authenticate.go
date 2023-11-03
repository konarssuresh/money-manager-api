package middleware

import (
	"context"
	"fmt"
	"money-manager/helper"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func AuthenticateMiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "userId", "")
		if !strings.Contains(r.RequestURI, "login") && !strings.Contains(r.RequestURI, "signup") {
			fmt.Println(r.Method)
			for name, values := range r.Header {
				for _, value := range values {
					fmt.Printf("%s: %s", name, value)
					fmt.Println("")
				}
			}
			bearerTokenString := r.Header.Get("Authorization")
			clientToken := strings.Replace(bearerTokenString, "Bearer ", "", 1)
			token, err := jwt.ParseWithClaims(clientToken, &helper.SignedDetails{}, func(token *jwt.Token) (interface{}, error) { return []byte(os.Getenv("JWT_KEY")), nil })
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			claimData, ok := token.Claims.(*helper.SignedDetails)

			if !ok {
				fmt.Println("invalid token")
				http.Error(w, "token is invalid", http.StatusBadRequest)
				return
			}
			ctx = context.WithValue(r.Context(), "userId", claimData.UserId)
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
