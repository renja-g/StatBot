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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/renja-g/StatBot/internal/db/gen"
)

var (
	dbConn  *pgx.Conn
	queries *gen.Queries
)

func main() {
	ctx := context.Background()
	dbURL := "postgres://postgres:postgres@localhost:5432/postgres"

	if envURL := os.Getenv("DATABASE_URL"); envURL != "" {
		dbURL = envURL
	}

	var err error
	dbConn, err = pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close(ctx)

	queries = gen.New(dbConn)

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

	statusChanges, err := queries.GetStatusChangesForDay(r.Context(), gen.GetStatusChangesForDayParams{
		UserID:  userId,
		Column2: pgtype.Date{Time: date, Valid: true},
		GuildID: guildId,
	})
	if err != nil {
		http.Error(w, "Failed to fetch status changes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type StatusPeriod struct {
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
		Status    string    `json:"status"`
	}

	cleanResponse := make([]StatusPeriod, 0, len(statusChanges))
	for _, change := range statusChanges {
		if !change.Status.Valid {
			continue
		}

		var endTime time.Time
		switch t := change.EndTime.(type) {
		case time.Time:
			endTime = t
		case pgtype.Timestamptz:
			endTime = t.Time
		default:
			endTime = change.StartTime.Time.Add(time.Hour)
		}

		cleanResponse = append(cleanResponse, StatusPeriod{
			StartTime: change.StartTime.Time,
			EndTime:   endTime,
			Status:    string(change.Status.DiscordStatus),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cleanResponse)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
