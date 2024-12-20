package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key used for signing the JWT
var secretKey = os.Getenv("JWT_ACCESS_SECRET_KEY") // Ensure this is set correctly

type contextKey string

const UserIDKey contextKey = "userID"

// JSONResponse writes a JSON response
func JSONResponse(w http.ResponseWriter, statusCode int, message string, success bool) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
		"success": success,
	})
}

// JWT Middleware to validate user tokens
func IsValidUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			JSONResponse(w, http.StatusUnauthorized, "Authorization cookie is missing", false)
			return
		}
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil // Return the secret key as a byte slice
		})
		if err != nil {
			JSONResponse(w, http.StatusUnauthorized, "Invalid token: "+err.Error(), false)
			return
		}
		// Check if the token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId, ok := claims["userId"].(string)
			if !ok {
				JSONResponse(w, http.StatusUnauthorized, "Invalid token: userId claim not found", false)
				return
			}
			// Set userId in the context using the correct key
			ctx := context.WithValue(r.Context(), UserIDKey, userId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			JSONResponse(w, http.StatusUnauthorized, "Invalid token", false)
		}
	})
}

// Checks if valid user but won't return if not
func IsValidUserWithoutReturn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the cookie
		cookie, err := r.Cookie("access_token")
		if err != nil {
			ctx := context.WithValue(r.Context(), UserIDKey, "")
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		// Parse the JWT token
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				ctx := context.WithValue(r.Context(), UserIDKey, "")
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(secretKey), nil
		})
		if err != nil {
			ctx := context.WithValue(r.Context(), UserIDKey, "")
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		// Validate token claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId, ok := claims["userId"].(string)
			if !ok {
				ctx := context.WithValue(r.Context(), UserIDKey, "")
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
			// Set userId in the context
			ctx := context.WithValue(r.Context(), UserIDKey, userId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		// Default case: invalid token
		ctx := context.WithValue(r.Context(), UserIDKey, "")
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		return
	})
}

// JWT Middleware to validate admin tokens
func IsValidAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			JSONResponse(w, http.StatusUnauthorized, "Authorization cookie is missing", false)
			return
		}
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil // Return the secret key as a byte slice
		})
		if err != nil {
			JSONResponse(w, http.StatusUnauthorized, "Invalid token: "+err.Error(), false)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["role"] == "admin" {
			userId, ok := claims["userId"].(string)
			if !ok {
				JSONResponse(w, http.StatusUnauthorized, "Invalid token: userId claim not found", false)
				return
			}
			// Set userId in the context using the correct key
			ctx := context.WithValue(r.Context(), UserIDKey, userId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			JSONResponse(w, http.StatusUnauthorized, "Invalid token", false)
		}
	})
}
