package main

import (
    "context"
    "habit-tracker/internal/auth"
    "habit-tracker/internal/handlers"
    "html/template"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    dbUrl := os.Getenv("DATABASE_URL")
    if dbUrl == "" {
        dbUrl = "postgres://postgres:postgres@db:5432/habitdb?sslmode=disable"
    }
    // Заменяем undefined ctx на context.Background()
    dbpool, err := pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        log.Fatalf("Unable to connect to DB: %v", err)
    }
    defer dbpool.Close()

    tmpl := template.Must(template.ParseGlob("internal/templates/*.html"))

    authService := auth.NewService(dbpool)
    handler := handlers.NewHandler(dbpool, tmpl, authService)

    mux := http.NewServeMux()

    mux.HandleFunc("/register", handler.Register)
    mux.HandleFunc("/login", handler.Login)
    mux.HandleFunc("/logout", handler.Logout)

    mux.HandleFunc("/habits", authService.RequireAuth(handler.Habits))
    mux.HandleFunc("/habits/add", authService.RequireAuth(handler.AddHabit))

    // mux.HandleFunc("/habits/edit", authService.RequireAuth(handler.EditHabit))
    // mux.HandleFunc("/habits/delete", authService.RequireAuth(handler.DeleteHabit))

    mux.HandleFunc("/records/mark", authService.RequireAuth(handler.MarkRecord))

    mux.HandleFunc("/report", authService.RequireAuth(handler.Report))

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    log.Println("Starting server on :8080")
    log.Fatal(srv.ListenAndServe())
}
