package main

import (
	"backend/models"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

const version = "1.0.0"

type config struct {
	port       int
	env        string
	secretCode string
	db         struct {
		dsn string
	}
	jwt struct {
		secret string
	}
}

type AppStatus struct {
	Environment string `json:"environment"`
	Status      string `json:"status"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {
	var cfg config
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	err := LoadConfig()
	if err != nil {
		logger.Fatal(err)
	}

	portStr := viper.GetString("PORT")
	defaultPort, err := strconv.Atoi(portStr)
	if err != nil {
		// Set a default port value if the environment variable is not set or cannot be converted to an integer.
		defaultPort = 4000
	}

	flag.IntVar(&cfg.port, "port", defaultPort, "Server port to listen on")
	flag.StringVar(&cfg.db.dsn, "dsn", viper.GetString("DB_CONNECTIONSTRING"), "mySQL connection string")
	flag.StringVar(&cfg.secretCode, "secretCode", viper.GetString("SECRET_CODE"), "registration secret code")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production)")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "secret")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", cfg.port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func LoadConfig() error {
	viper.AddConfigPath(".")

	viper.SetConfigName("app")

	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
