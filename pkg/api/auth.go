package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SignInRequest структура для запроса входа
type SignInRequest struct {
	Password string `json:"password"`
}

// SignInResponse структура для ответа с токеном
type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// hashPassword создает SHA256 хэш пароля
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// SignInHandler обрабатывает POST /api/signin
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

	envPassword := os.Getenv("TODO_PASSWORD")
	if envPassword == "" || req.Password != envPassword {
		writeJSONError(w, http.StatusUnauthorized, "Неверный пароль")
		return
	}

	// Создаем JWT токен
	passHash := hashPassword(envPassword)
	claims := jwt.MapClaims{
		"pass_hash": passHash,
		"exp":       time.Now().Add(8 * time.Hour).Unix(), // Токен живет 8 часов
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// В качестве секретного ключа используем хэш пароля.
	// Если пароль в ENV изменится, старые токены станут невалидными.
	tokenString, err := token.SignedString([]byte(passHash))
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "ошибка генерации токена")
		return
	}

	writeJSON(w, SignInResponse{Token: tokenString})
}

// AuthMiddleware проверяет наличие и валидность JWT токена
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		envPassword := os.Getenv("TODO_PASSWORD")

		// Если пароль не задан, аутентификация не требуется
		if envPassword == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		passHash := hashPassword(envPassword)

		// Парсим и валидируем токен
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte(passHash), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Дополнительная проверка: совпадает ли хэш внутри токена с текущим хэшем пароля
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["pass_hash"] != passHash {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Если все ок, передаем управление следующему обработчику
		next(w, r)
	}
}
