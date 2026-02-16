package auth

import (
    "net/http"
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "golang.org/x/crypto/bcrypt"
)

type Service struct {
    db *pgxpool.Pool
}

// Создание сервиса аутентификации
func NewService(db *pgxpool.Pool) *Service {
    return &Service{db: db}
}

func (s *Service) HashPassword(password string) ([]byte, error) {
    return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (s *Service) CheckPassword(hash []byte, password string) error {
    return bcrypt.CompareHashAndPassword(hash, []byte(password))
}

// Регистрация пользователя
func (s *Service) Register(ctx context.Context, username, password string) error {
    hash, err := s.HashPassword(password)
    if err != nil {
        return err
    }
    _, err = s.db.Exec(ctx,
        "INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, string(hash))
    return err
}

// Проверка данных и возврат user_id
func (s *Service) Authenticate(ctx context.Context, username, password string) (string, error) {
    var hash string
    var userID string
    row := s.db.QueryRow(ctx, "SELECT id, password_hash FROM users WHERE username = $1", username)
    err := row.Scan(&userID, &hash)
    if err != nil {
        return "", err
    }
    err = s.CheckPassword([]byte(hash), password)
    if err != nil {
        return "", err
    }
    return userID, nil
}

// Middleware для проверки сессии (упрощенно - cookie с user_id)
func (s *Service) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("session_user")
        if err != nil || cookie.Value == "" {
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }
        // Можно проверить, что userID есть в базе если нужно
        next.ServeHTTP(w, r)
    }
}
