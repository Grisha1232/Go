package services

import (
	"database/sql"
	"errors"
	"strings"
	"time"
	"todo-app/internal/repositories"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Секретный ключ для подписи JWT
var jwtSecret = []byte("supersecretkey")

// Claims структура для хранения ID пользователя в JWT
type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// Authenticate проверяет логин и пароль, возвращает JWT
func Authenticate(db *sql.DB, username, password string) (string, error) {
	// Получаем пользователя из БД
	user, err := repositories.GetUserByUsername(db, username)
	if err != nil {
		return "", errors.New("пользователь не найден")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("неверный пароль")
	}

	// Создаем JWT-токен
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken проверяет и декодирует JWT, возвращает userID
func ValidateToken(tokenString string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("недействительный токен")
	}

	return claims.UserID, nil
}

// RegisterUser создает нового пользователя
func RegisterUser(db *sql.DB, username, password string) error {
	// Проверяем, существует ли пользователь
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	if err != nil {
		return errors.New("ошибка проверки пользователя")
	}
	if exists {
		return errors.New("пользователь уже существует")
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("ошибка хеширования пароля")
	}

	// Записываем пользователя в БД
	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("пользователь уже существует")
		}
		return errors.New("ошибка при создании пользователя")
	}

	return nil
}
