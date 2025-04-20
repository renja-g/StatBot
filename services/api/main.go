package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renja-g/StatBot/internal/db/gen"
)

var (
	dbPool  *pgxpool.Pool
	queries *gen.Queries
)

func main() {
	ctx := context.Background()
	dbURL := "postgres://postgres:postgres@localhost:5432/postgres"

	if envURL := os.Getenv("DATABASE_URL"); envURL != "" {
		dbURL = envURL
	}

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Failed to parse database URL: %v", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2

	// Create the connection pool
	dbPool, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Failed to create db connection pool: %v", err)
	}
	defer dbPool.Close()

	// Test connection
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	queries = gen.New(dbPool)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)

	r.Get("/api/hello", helloHandler)
	r.Get("/api/status-changes/{guildId}/{userId}", getStatusChangesHandler)

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func getStatusChangesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	guildIdStr := chi.URLParam(r, "guildId")
	userIdStr := chi.URLParam(r, "userId")
	dateStr := r.URL.Query().Get("date")

	guildId, err := strconv.ParseInt(guildIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
		return
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var date time.Time
	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	statusChanges, err := queries.GetUserStatusTimeline(ctx, gen.GetUserStatusTimelineParams{
		Column1: guildId,
		Column2: userId,
		Column3: pgtype.Date{Time: date, Valid: true},
	})
	if err != nil {
		http.Error(w, "Failed to fetch status changes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type SimpleStatus struct {
		StartTime string `json:"StartTime"`
		EndTime   string `json:"EndTime"`
		Status    string `json:"Status"`
	}

	result := make([]SimpleStatus, 0, len(statusChanges))
	for _, s := range statusChanges {
		// Try to convert to time.Time for proper formatting
		var startTimeStr, endTimeStr string

		// Handle StartTime
		if t, ok := s.StartTime.(time.Time); ok {
			startTimeStr = t.Format(time.RFC3339)
		} else {
			startTimeStr = fmt.Sprintf("%v", s.StartTime)
		}

		// Handle EndTime
		if t, ok := s.EndTime.(time.Time); ok {
			endTimeStr = t.Format(time.RFC3339)
		} else {
			endTimeStr = fmt.Sprintf("%v", s.EndTime)
		}

		result = append(result, SimpleStatus{
			StartTime: startTimeStr,
			EndTime:   endTimeStr,
			Status:    string(s.Status.DiscordStatus),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
