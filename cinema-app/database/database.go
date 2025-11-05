package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./cinema.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
	seedData()
}

func createTables() {
	usersTable := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
        )`

	moviesTable := `
        CREATE TABLE IF NOT EXISTS movies (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            genre TEXT NOT NULL,
            release_date DATETIME,
            rating REAL,
            director TEXT,
			UNIQUE(title, release_date)
        )`

	if _, err := DB.Exec(usersTable); err != nil {
		log.Fatal(err)
	}
	if _, err := DB.Exec(moviesTable); err != nil {
		log.Fatal(err)
	}
}

func seedData() {
	// Добавим тестового пользователя
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	DB.Exec(`INSERT OR IGNORE INTO users (username, password) VALUES (?, ?)`,
		"testuser", string(hashedPassword))

	// Добавляем фильмы напрямую в базу, без использования структуры models.Movie
	// movies := []struct {
	// 	title       string
	// 	genre       string
	// 	releaseDate string
	// 	rating      float64
	// 	director    string
	// }{
	// 	{"Inception", "Sci-Fi", "2010-07-16", 8.8, "Christopher Nolan"},
	// 	{"The Shawshank Redemption", "Drama", "1994-09-23", 9.3, "Frank Darabont"},
	// 	{"The Dark Knight", "Action", "2008-07-18", 9.0, "Christopher Nolan"},
	// 	{"Pulp Fiction", "Crime", "1994-10-14", 8.9, "Quentin Tarantino"},
	// 	{"Forrest Gump", "Drama", "1994-07-06", 8.8, "Robert Zemeckis"},
	// }

	// for _, movie := range movies {
	// 	// Преобразуем строку даты в time.Time для базы данных
	// 	releaseTime, err := time.Parse("2006-01-02", movie.releaseDate)
	// 	if err != nil {
	// 		log.Printf("Error parsing date %s: %v", movie.releaseDate, err)
	// 		continue
	// 	}

	// 	DB.Exec(`INSERT OR IGNORE INTO movies (title, genre, release_date, rating, director)
	//              VALUES (?, ?, ?, ?, ?)`,
	// 		movie.title, movie.genre, releaseTime, movie.rating, movie.director)
	// }

	// log.Println("Database seeded successfully")
}
