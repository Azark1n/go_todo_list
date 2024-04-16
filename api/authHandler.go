package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type SignInRequest struct {
	Password string `json:"password"`
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		sendErrorResponse(w, "Authentication not configured", http.StatusInternalServerError)
		return
	}

	if req.Password != expectedPassword {
		sendErrorResponse(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": generateSHA256Hash(expectedPassword),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(expectedPassword))
	if err != nil {
		sendErrorResponse(w, "Failed to sign token", http.StatusBadRequest)
		return
	}

	sendJSON(w, map[string]string{"token": tokenString}, http.StatusOK)
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass != "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				sendErrorResponse(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(pass), nil
			})

			if err != nil || !token.Valid {
				sendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				hashInToken, ok := claims["hash"].(string)
				if !ok || generateSHA256Hash(pass) != hashInToken {
					sendErrorResponse(w, "Invalid token hash", http.StatusUnauthorized)
					return
				}
			} else {
				sendErrorResponse(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

func generateSHA256Hash(input string) string {
	data := []byte(input)
	hash := sha256.Sum256(data)
	hashString := hex.EncodeToString(hash[:])
	return hashString
}
