package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ngfenglong/food-randomizer-BE/pkg/config"
	"github.com/ngfenglong/food-randomizer-BE/pkg/database"
	"github.com/ngfenglong/food-randomizer-BE/pkg/http/router"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

const version = "1.0.0"

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err)
	}

	portStr := viper.GetString("PORT")
	defaultPort, err := strconv.Atoi(portStr)
	if err != nil {
		// Set a default port value if the environment variable is not set or cannot be converted to an integer.
		defaultPort = 4000
	}

	flag.IntVar(&cfg.Server.Port, "port", defaultPort, "Server port to listen on")
	flag.StringVar(&cfg.Database.DSN, "dsn", viper.GetString("DB_CONNECTIONSTRING"), "mySQL connection string")
	flag.StringVar(&cfg.SecretCode, "secretCode", viper.GetString("SECRET_CODE"), "registration secret code")
	flag.StringVar(&cfg.Env, "env", "development", "Application environment (development|production)")
	flag.StringVar(&cfg.JWT.Secret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "secret")
	flag.Parse()

	db, err := database.OpenDB(cfg.Database)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	r := router.NewRouter(cfg, db)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", cfg.Server.Port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}
