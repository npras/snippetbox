package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/npras/snippetbox/internal/models"
)

type config struct {
	port      string
	staticDir string
	dsn       string
}

type application struct {
	config         *config
	logger         *slog.Logger
	snippet        *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func newLogger() *slog.Logger {
	logOpts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	return slog.New(slog.NewTextHandler(os.Stdout, logOpts))
}

func parseFlags(c *config) {
	flag.StringVar(&c.port, "port", ":4000", "port in which the server listens")
	flag.StringVar(&c.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&c.dsn, "dsn", "postgresql://postgres:golanger1234567@localhost:5432/snippetbox", "PostgreSQL data source name")
	flag.Parse()
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		return nil, err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func newSessionManager(pool *pgxpool.Pool) *scs.SessionManager {
	sm := scs.New()
	sm.Store = pgxstore.New(pool)
	sm.Lifetime = 12 * time.Hour
	return sm
}

//

func main() {
	logger := newLogger()

	config := &config{}
	parseFlags(config)

	pool, err := openDB(config.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:         logger,
		config:         config,
		snippet:        &models.SnippetModel{DbPool: pool},
		templateCache:  templateCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: newSessionManager(pool),
	}

	app.logger.Info("starting server on", slog.String("port", app.config.port))
	err = http.ListenAndServe(app.config.port, app.routes())
	app.logger.Error(err.Error())
	os.Exit(1)
}
