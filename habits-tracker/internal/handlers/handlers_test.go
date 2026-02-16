package handlers

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"habit-tracker/internal/auth"
	"habit-tracker/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testHandler struct {
	*Handler
	authService *auth.Service
}

func setupTestHandler(t *testing.T) *testHandler {
	db, cleanup := testutils.SetupTestDB()
	t.Cleanup(cleanup)

	tmpl := template.Must(template.ParseGlob("../../internal/templates/*.html"))
	authService := auth.NewService(db)
	handler := NewHandler(db, tmpl, authService)

	return &testHandler{
		Handler:     handler,
		authService: authService,
	}
}

func (th *testHandler) createTestUser(t *testing.T, username, password string) string {
	ctx := context.Background()
	err := th.authService.Register(ctx, username, password)
	require.NoError(t, err)

	userID, err := th.authService.Authenticate(ctx, username, password)
	require.NoError(t, err)

	return userID
}

func (th *testHandler) setAuthCookie(r *http.Request, userID string) {
	r.AddCookie(&http.Cookie{
		Name:  "session_user",
		Value: userID,
	})
}

func TestRegisterHandler(t *testing.T) {
	th := setupTestHandler(t)

	t.Run("GET Register", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/register", nil)
		rr := httptest.NewRecorder()

		th.Register(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "register")
	})

	t.Run("POST Register Success", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "newuser")
		form.Add("password", "password123")

		req := httptest.NewRequest("POST", "/register", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		th.Register(rr, req)

		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Equal(t, "/login", rr.Header().Get("Location"))
	})

	t.Run("POST Register Duplicate User", func(t *testing.T) {
		form1 := url.Values{}
		form1.Add("username", "duplicate")
		form1.Add("password", "password123")

		req1 := httptest.NewRequest("POST", "/register", strings.NewReader(form1.Encode()))
		req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr1 := httptest.NewRecorder()
		th.Register(rr1, req1)

		form2 := url.Values{}
		form2.Add("username", "duplicate")
		form2.Add("password", "password123")

		req2 := httptest.NewRequest("POST", "/register", strings.NewReader(form2.Encode()))
		req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr2 := httptest.NewRecorder()

		th.Register(rr2, req2)

		assert.Equal(t, http.StatusBadRequest, rr2.Code)
		assert.Contains(t, rr2.Body.String(), "Username taken")
	})
}

func TestLoginHandler(t *testing.T) {
	th := setupTestHandler(t)
	userID := th.createTestUser(t, "loginuser", "loginpass")

	t.Run("GET Login", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/login", nil)
		rr := httptest.NewRecorder()

		th.Login(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "login")
	})

	t.Run("POST Login Success", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "loginuser")
		form.Add("password", "loginpass")

		req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		th.Login(rr, req)

		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Equal(t, "/habits", rr.Header().Get("Location"))

		cookies := rr.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session_user" {
				sessionCookie = cookie
				break
			}
		}
		require.NotNil(t, sessionCookie)
		assert.Equal(t, userID, sessionCookie.Value)
	})

	t.Run("POST Login Invalid Credentials", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "loginuser")
		form.Add("password", "wrongpassword")

		req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		th.Login(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid credentials")
	})
}

func TestHabitsHandler(t *testing.T) {
	th := setupTestHandler(t)
	userID := th.createTestUser(t, "habitsuser", "habitpass")

	t.Run("Habits Without Auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/habits", nil)
		rr := httptest.NewRecorder()

		th.Habits(rr, req)

		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Equal(t, "/login", rr.Header().Get("Location"))
	})

	t.Run("Habits With Auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/habits", nil)
		th.setAuthCookie(req, userID)
		rr := httptest.NewRecorder()

		th.Habits(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "habits")
	})
}

func TestAddHabitHandler(t *testing.T) {
	th := setupTestHandler(t)
	userID := th.createTestUser(t, "addhabituser", "addhabitpass")

	t.Run("GET Add Habit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/habits/add", nil)
		th.setAuthCookie(req, userID)
		rr := httptest.NewRecorder()

		th.AddHabit(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "habit_form")
	})

	t.Run("POST Add Habit", func(t *testing.T) {
		form := url.Values{}
		form.Add("description", "Test Habit")
		form.Add("frequency", "5")
		form.Add("target_percent", "80")

		req := httptest.NewRequest("POST", "/habits/add", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		th.setAuthCookie(req, userID)
		rr := httptest.NewRecorder()

		th.AddHabit(rr, req)

		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Equal(t, "/habits", rr.Header().Get("Location"))
	})

	t.Run("POST Add Habit Invalid Data", func(t *testing.T) {
		form := url.Values{}
		form.Add("description", "Test Habit")
		form.Add("frequency", "invalid")
		form.Add("target_percent", "80")

		req := httptest.NewRequest("POST", "/habits/add", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		th.setAuthCookie(req, userID)
		rr := httptest.NewRecorder()

		th.AddHabit(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid frequency")
	})
}

func TestMarkRecordHandler(t *testing.T) {
	th := setupTestHandler(t)
	userID := th.createTestUser(t, "markuser", "markpass")

	habitID := th.createTestHabit(t, userID, "Test Habit", 5, 80)

	t.Run("Mark Record Success", func(t *testing.T) {
		form := url.Values{}
		form.Add("habit_id", habitID)
		form.Add("date", time.Now().Format("2006-01-02"))
		form.Add("done", "true")

		req := httptest.NewRequest("POST", "/records/mark", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		th.setAuthCookie(req, userID)
		rr := httptest.NewRecorder()

		th.MarkRecord(rr, req)

		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Equal(t, "/habits", rr.Header().Get("Location"))
	})

	t.Run("Mark Record Invalid Habit", func(t *testing.T) {
		form := url.Values{}
		form.Add("habit_id", "invalid-uuid")
		form.Add("date", time.Now().Format("2006-01-02"))
		form.Add("done", "true")

		req := httptest.NewRequest("POST", "/records/mark", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		th.setAuthCookie(req, userID)
		rr := httptest.NewRecorder()

		th.MarkRecord(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})
}

func TestReportHandler(t *testing.T) {
	th := setupTestHandler(t)
	userID := th.createTestUser(t, "reportuser", "reportpass")

	t.Run("Report Without Auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/report", nil)
		rr := httptest.NewRecorder()

		th.Report(rr, req)

		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Equal(t, "/login", rr.Header().Get("Location"))
	})

	t.Run("Report With Auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/report?period=week", nil)
		th.setAuthCookie(req, userID)
		rr := httptest.NewRecorder()

		th.Report(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "report")
	})
}

func (th *testHandler) createTestHabit(t *testing.T, userID, description string, frequency, targetPercent int) string {
	var habitID string
	err := th.db.QueryRow(context.Background(),
		"INSERT INTO habits (user_id, description, frequency, target_percent) VALUES ($1, $2, $3, $4) RETURNING id",
		userID, description, frequency, targetPercent,
	).Scan(&habitID)
	require.NoError(t, err)
	return habitID
}
