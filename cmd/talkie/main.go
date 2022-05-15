package main

import (
	"fmt"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"

	"github.com/quantonganh/talkie"
	"github.com/quantonganh/talkie/http"
	"github.com/quantonganh/talkie/postgres"
)

func main() {
	_ = godotenv.Load("./cmd/talkie/.env")

	var cfg talkie.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Database: cfg.DBName,
	})
	defer db.Close()

	userSvc := postgres.NewUserService(db)
	commentSvc := postgres.NewCommentService(db)
	s := http.NewServer()
	s.UserService = userSvc
	s.CommentService = commentSvc
	if err := s.Run(cfg.Port); err != nil {
		log.Fatal(err)
	}
}
