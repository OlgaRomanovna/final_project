package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Настройки аутентификации, заполняются один раз при старте через Init
var (
	authPassword string
	authPassHash string
)

// Init сохраняет конфигурацию аутентификации.
func Init(password string) {
	authPassword = password
	authPassHash = hashPassword(password)
}

type SignInRequest struct {
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "метод не разрешен")
		return
	}

	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}

	if authPassword == "" || req.Password != authPassword {
		writeJSONError(w, http.StatusUnauthorized, "Неверный пароль")
		return
	}

	claims := jwt.MapClaims{
		"pass_hash": authPassHash,
		"exp":       time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(authPassHash))
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "ошибка генерации токена")
		return
	}

	writeJSON(w, SignInResponse{Token: tokenString})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Если пароль не задан, аутентификация отключена
		if authPassword == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte(authPassHash), nil
		}, jwt.WithValidMethods([]string{"HS256"}))

		if err != nil || !token.Valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["pass_hash"] != authPassHash {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
