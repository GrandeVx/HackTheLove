package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"backend/config"
	"backend/internal/database"
	"backend/internal/logger"
	"backend/internal/server/routes"
	"backend/internal/services"
)

type Server struct {
	port int
	db   database.Database
	http *http.Server
}

var log = logger.GetLogger()

func NewServer() *Server {
	port, err := strconv.Atoi(config.AppPort)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse APP_PORT environment variable")
	}

	db := database.New()
	log.Info().Msg("Initializing database tables")
	if err := db.InitTables(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database tables")
	}

	router := routes.InitRoutes(*db)

	server := &Server{
		port: port,
		db:   *db,
		http: &http.Server{
			Addr:         fmt.Sprintf("0.0.0.0:%d", port),
			Handler:      router,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}

	log.Info().Msgf("Server is starting on port %d", port)
	return server
}

func scheduleMatching(runAt time.Time, task func()) {
	delay := time.Until(runAt)
	if delay <= 0 {
		log.Warn().Msg("Task already expired")
		return
	}

	log.Info().Msgf("Task scheduled for: %v", runAt)
	time.AfterFunc(delay, task)
}

func (s *Server) Run() error {
	log.Info().Msg("Starting HTTP server")

	surveyService := services.NewSurveyService(s.db.GetDB())
	matchService := services.NewMatchService(s.db.GetDB())

	scheduleTime, err := time.Parse(time.RFC3339, config.ScheduleTime)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing schedule time. The format should be RFC3339 (e.g., 2021-09-01T16:00:00Z)")
	}

	scheduleMatching(scheduleTime, func() {
		if err := surveyService.StartMatching(matchService); err != nil {
			log.Error().Err(err).Msg("Error starting matching")
		}
	})

	log.Info().Msg("Server is now running and listening for requests")
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down server...")

	if err := s.db.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing database connection")
	}

	return s.http.Shutdown(ctx)
}
