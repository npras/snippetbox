package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/npras/snippetbox/internal/models"
)

type config struct {
	port      string
	staticDir string
	dsn       string
}

type application struct {
	config  *config
	logger  *slog.Logger
	snippet *models.SnippetModel
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
	flag.StringVar(&c.dsn, "dsn", "postgresql://web:golanger1234567@localhost:5432/snippetbox", "PostgreSQL data source name")
	flag.Parse()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

//

func main() {
	logger := newLogger()

	config := &config{}
	parseFlags(config)

	db, err := openDB(config.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger:  logger,
		config:  config,
		snippet: &models.SnippetModel{DB: db},
	}

	var greeting string
	err = db.QueryRow("select content from snippets limit 1").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)

	app.logger.Info("starting server on", slog.String("port", app.config.port))
	err = http.ListenAndServe(app.config.port, app.routes())
	app.logger.Error(err.Error())
	os.Exit(1)
}
