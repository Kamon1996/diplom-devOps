package handlers

import (
	"context"
	"habit-tracker/internal/auth"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db          *pgxpool.Pool
	templates   *template.Template
	authService *auth.Service
}

func NewHandler(db *pgxpool.Pool, tmpl *template.Template, authService *auth.Service) *Handler {
	return &Handler{
		db:          db,
		templates:   tmpl,
		authService: authService,
	}
}

// Регистрация
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.templates.ExecuteTemplate(w, "register.html", nil)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	ctx := r.Context()
	err := h.authService.Register(ctx, username, password)
	if err != nil {
		http.Error(w, "Username taken or error: "+err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

// Логин
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.templates.ExecuteTemplate(w, "login.html", nil)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	ctx := r.Context()
	userID, err := h.authService.Authenticate(ctx, username, password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_user",
		Value:    userID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	http.Redirect(w, r, "/habits", http.StatusFound)
}

// Логаут
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_user",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

func getUserIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_user")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

type Habit struct {
	ID            string
	Description   string
	Frequency     int
	TargetPercent int
}

// Загрузка привычек пользователя из БД
func (h *Handler) loadHabits(ctx context.Context, userID string) ([]Habit, error) {
	rows, err := h.db.Query(ctx, "SELECT id, description, frequency, target_percent FROM habits WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []Habit
	for rows.Next() {
		var habit Habit
		err := rows.Scan(&habit.ID, &habit.Description, &habit.Frequency, &habit.TargetPercent)
		if err != nil {
			return nil, err
		}
		habits = append(habits, habit)
	}
	return habits, nil
}

// Страница привычек
func (h *Handler) Habits(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	habits, err := h.loadHabits(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to load habits", http.StatusInternalServerError)
		return
	}

	// Добавляем Today в каждый элемент
	today := time.Now().Format("2006-01-02")
	type HabitWithToday struct {
		ID            string
		Description   string
		Frequency     int
		TargetPercent int
		Today         string
	}
	var habitsWithToday []HabitWithToday
	for _, hbt := range habits {
		habitsWithToday = append(habitsWithToday, HabitWithToday{
			ID:            hbt.ID,
			Description:   hbt.Description,
			Frequency:     hbt.Frequency,
			TargetPercent: hbt.TargetPercent,
			Today:         today,
		})
	}

	h.templates.ExecuteTemplate(w, "habits.html", habitsWithToday)
}


// Добавление привычки
func (h *Handler) AddHabit(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == http.MethodGet {
		h.templates.ExecuteTemplate(w, "habit_form.html", nil)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	desc := r.Form.Get("description")
	freqStr := r.Form.Get("frequency")
	targetStr := r.Form.Get("target_percent")

	freq, err := strconv.Atoi(freqStr)
	if err != nil {
		http.Error(w, "Invalid frequency", http.StatusBadRequest)
		return
	}
	target, err := strconv.Atoi(targetStr)
	if err != nil {
		http.Error(w, "Invalid target percent", http.StatusBadRequest)
		return
	}

	// Вставка привычки в базу
	_, err = h.db.Exec(r.Context(),
		"INSERT INTO habits (user_id, description, frequency, target_percent) VALUES ($1, $2, $3, $4)",
		userID, desc, freq, target)
	if err != nil {
		http.Error(w, "Failed to add habit: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/habits", http.StatusFound)
}

// Отметка выполнения привычки за день (маркировка)
func (h *Handler) MarkRecord(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	habitID := r.Form.Get("habit_id")
	dateStr := r.Form.Get("date") // строка даты в формате "YYYY-MM-DD"
	doneStr := r.Form.Get("done") // "true" или "false"

	done := doneStr == "true"

	// Проверить, что привычка принадлежит пользователю (упрощенно)
	var count int
	err = h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM habits WHERE id=$1 AND user_id=$2", habitID, userID).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, "Habit not found or access denied", http.StatusForbidden)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Вставить или обновить запись о выполнении привычки за день
	_, err = h.db.Exec(r.Context(), `
		INSERT INTO habit_records (habit_id, date, done)
		VALUES ($1, $2, $3)
		ON CONFLICT (habit_id, date) DO UPDATE SET done = EXCLUDED.done
	`, habitID, date, done)
	if err != nil {
		http.Error(w, "Failed to save record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/habits", http.StatusFound)
}

// Отчет (сводный за неделю/месяц/год)
// Для простоты предоставляет данные по привычкам и выполнению, можно расширить для отображения графика и PDF генерации
func (h *Handler) Report(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	period := r.URL.Query().Get("period") // "week", "month", "year"
	if period == "" {
		period = "week"
	}

	var start time.Time
	now := time.Now()
	switch period {
	case "week":
		start = now.AddDate(0, 0, -7)
	case "month":
		start = now.AddDate(0, -1, 0)
	case "year":
		start = now.AddDate(-1, 0, 0)
	default:
		start = now.AddDate(0, 0, -7)
	}

	// Получим данные: привычки и % выполнения за период
	type ReportData struct {
		Description string
		Frequency   int
		Target      int
		DonePercent float64
	}

	rows, err := h.db.Query(r.Context(), `
		SELECT h.description, h.frequency, h.target_percent,
		COALESCE(100.0 * SUM(CASE WHEN r.done THEN 1 ELSE 0 END) / NULLIF(COUNT(r.done),0), 0) as done_percent
		FROM habits h
		LEFT JOIN habit_records r ON r.habit_id = h.id AND r.date >= $1
		WHERE h.user_id = $2
		GROUP BY h.id
	`, start, userID)
	if err != nil {
		http.Error(w, "Failed to load report data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reports []ReportData
	for rows.Next() {
		var r ReportData
		if err := rows.Scan(&r.Description, &r.Frequency, &r.Target, &r.DonePercent); err != nil {
			http.Error(w, "Failed to read report row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		reports = append(reports, r)
	}

	// TODO: сюда можно добавить генерацию PDF (go-pdf или gofpdf)

	h.templates.ExecuteTemplate(w, "report.html", struct {
		Period  string
		Reports []ReportData
	}{
		Period:  period,
		Reports: reports,
	})
}
