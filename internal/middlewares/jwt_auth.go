package middlewares

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// Secret key used for signing the JWT
var secretKey = os.Getenv("MONGO_URI") // Ensure this is set correctly

// JWT Middleware
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the cookie
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Authorization cookie is missing", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil // Return the secret key as a byte slice
		})

		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Check if the token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// You can access claims here if needed
			fmt.Println("User  ID:", claims["userId"]) // Example of accessing a claim
			next.ServeHTTP(w, r)                       // Proceed to the next handler
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}
	})
}
