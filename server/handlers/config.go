package handlers

import (
	"VideoHub/server/utils"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	jwtSecret = os.Getenv("jwt_secret")
	jwtDuration = os.Getenv("jwt_duration")
)

type Config struct {
	JwtManager *utils.JWTManager
}

func NewConfig() *Config {
	if len(jwtSecret) == 0 {
		jwtSecret = "debug_secret"
	}
	if len(jwtDuration) == 0 {
		jwtDuration = "8600000"
	}

	duration, err := time.ParseDuration(fmt.Sprintf("%sms", jwtDuration))
	if err != nil {
		log.Fatal(err)
	}

	return &Config{ utils.NewJWTManager(jwtSecret, duration) }
}
