package main

import (
	"errors"
	"net/http"
	"os"
	// "strings"

	"github.com/golang-jwt/jwt"

	"github.com/moov-io/base/log"	
)

func withAuth(logger log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookie
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				logger.LogErrorf("No cookie found")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}			
			logger.LogErrorf("Error occured while reading cookie")
			http.Error(w, "Bad Request", http.StatusBadRequest)			
			return
		}
		tokenString := tokenCookie.Value
		logger.Logf("Token from cookie: %s", tokenString)

		// // Get the JWT from the "Authorization" header
		// authHeader := r.Header.Get("Authorization")
		// if authHeader == "" {
		// 	logger.LogErrorf("Request missing authorization from %s to %s", r.RemoteAddr, r.URL.Path)
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }
		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			logger.LogErrorf("Request missing token from %s to %s", r.RemoteAddr, r.URL.Path)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Parse the JWT
		secret := os.Getenv("MY_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm used to sign the JWT
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			// Return the secret
			return []byte(secret), nil
		})
		if err != nil {
			logger.LogErrorf("Request error from %s to %s: %v", r.RemoteAddr, r.URL.Path, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Check if the JWT is valid
		if !token.Valid {
			logger.LogErrorf("Request with invalid token from %s to %s", r.RemoteAddr, r.URL.Path)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
