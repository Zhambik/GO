package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Genre       string    `json:"genre"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float64   `json:"rating"`
	Director    string    `json:"director"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
