package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type config struct {
	port      string
	staticDir string
	dsn       string
}

type application struct {
	config config
	logger *slog.Logger
}

func newLogger() *slog.Logger {
	logOpts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	return slog.New(slog.NewTextHandler(os.Stdout, logOpts))
}

func parseFlags(app *application) {
	flag.StringVar(&app.config.port, "port", ":4000", "port in which the server listens")
	flag.StringVar(&app.config.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&app.config.dsn, "dsn", "postgresql://web:golanger1234567@localhost:5432/snippetbox", "PostgreSQL data source name")
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
	app := &application{logger: newLogger()}
	parseFlags(app)

	db, err := openDB(app.config.dsn)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

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
