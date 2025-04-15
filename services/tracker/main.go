package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/renja-g/StatBot/internal/db/gen"
)

var (
	token   string
	guildID snowflake.ID
	dbConn  *pgx.Conn
	queries *gen.Queries
)

func init() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	token = os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN environment variable is required")
	}

	guildIDStr := os.Getenv("GUILD_ID")
	if guildIDStr == "" {
		log.Fatal("GUILD_ID environment variable is required")
	}

	var parseErr error
	guildID, parseErr = snowflake.Parse(guildIDStr)
	if parseErr != nil {
		log.Fatalf("Failed to parse GUILD_ID: %v", parseErr)
	}
}

func main() {
	slog.Info("starting example...")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	// Connect to database
	ctx := context.Background()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/postgres"
	}

	var err error
	dbConn, err = pgx.Connect(ctx, dbURL)
	if err != nil {
		slog.Error("failed to connect to database", slog.Any("err", err))
		return
	}
	defer dbConn.Close(ctx)

	// Initialize queries
	queries = gen.New(dbConn)

	client, err := disgo.New(token,
		bot.WithDefaultGateway(),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnPresenceUpdate: onPresenceUpdate,
		}),
		bot.WithGatewayConfigOpts(gateway.WithIntents(
			gateway.IntentGuildPresences,
			gateway.IntentGuildMessages,
			gateway.IntentDirectMessages,
		)),
	)
	if err != nil {
		slog.Error("error while building disgo instance", slog.Any("err", err))
		return
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to gateway", slog.Any("err", err))
	}

	slog.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onPresenceUpdate(event *events.PresenceUpdate) {
	if event.GuildID != guildID {
		return
	}
	slog.Info("Update detected")

	slog.Info(fmt.Sprintf("%s", event.Presence.ClientStatus.Desktop))
	slog.Info(fmt.Sprintf("%s", event.Presence.ClientStatus.Mobile))
	slog.Info(fmt.Sprintf("%s", event.Presence.ClientStatus.Web))

	// Serialize activities to JSON
	activitiesJSON, err := json.Marshal(event.Presence.Activities)
	if err != nil {
		slog.Error("failed to marshal activities", slog.Any("err", err))
		activitiesJSON = []byte("[]") // Use empty array as fallback
	}

	// Save presence update to database
	err = queries.CreatePresenceUpdate(context.Background(), gen.CreatePresenceUpdateParams{
		Timestamp: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		UserID:  int64(event.PresenceUser.ID),
		GuildID: int64(event.GuildID),
		ClientStatusDesktop: gen.NullDiscordStatus{
			DiscordStatus: gen.DiscordStatus(event.Presence.ClientStatus.Desktop),
			Valid:         event.Presence.ClientStatus.Desktop != "",
		},
		ClientStatusMobile: gen.NullDiscordStatus{
			DiscordStatus: gen.DiscordStatus(event.Presence.ClientStatus.Mobile),
			Valid:         event.Presence.ClientStatus.Mobile != "",
		},
		ClientStatusWeb: gen.NullDiscordStatus{
			DiscordStatus: gen.DiscordStatus(event.Presence.ClientStatus.Web),
			Valid:         event.Presence.ClientStatus.Web != "",
		},
		Activities: activitiesJSON,
	})

	if err != nil {
		slog.Error("failed to save presence update", slog.Any("err", err))
	} else {
		slog.Info("presence update saved to database")
	}
}
