package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "cinema.db")
	if err != nil {
		log.Fatal(err)
	}

	// Создаем таблицы
	if err := createTables(db); err != nil {
		log.Fatal(err)
	}

	// Наполняем тестовыми данными
	if err := seedData(db); err != nil {
		log.Fatal(err)
	}

	//log.Println("База данных инициализирована")
	return db
}

func createTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS movies (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            duration INTEGER NOT NULL,
            age_rating TEXT NOT NULL CHECK(age_rating IN ('0+', '6+', '12+', '16+', '18+'))
        )`,

		`CREATE TABLE IF NOT EXISTS screenings (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            movie_id INTEGER NOT NULL,
            screen_time DATETIME NOT NULL,
            hall_number INTEGER NOT NULL,
            price DECIMAL(10, 2) NOT NULL,
            FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE
        )`,

		`CREATE TABLE IF NOT EXISTS tickets (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            screening_id INTEGER NOT NULL,
            seat_number INTEGER NOT NULL,
            purchase_time DATETIME DEFAULT CURRENT_TIMESTAMP,
            customer_name TEXT NOT NULL,
            FOREIGN KEY (screening_id) REFERENCES screenings (id) ON DELETE CASCADE,
            UNIQUE(screening_id, seat_number)
        )`,

		`CREATE TABLE IF NOT EXISTS customers (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL,
            phone TEXT,
            loyalty_points INTEGER DEFAULT 0
        )`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return err
		}
	}
	return nil
}

func seedData(db *sql.DB) error {
	// Проверяем, есть ли уже данные
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM movies").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Данные уже есть
	}

	// Добавляем тестовые данные фильмов
	movies := []struct {
		title     string
		duration  int
		ageRating string
	}{
		{"Дюна: Часть вторая", 166, "12+"},
		{"Оппенгеймер", 180, "16+"},
		{"Человек-паук: Паутина вселенных", 140, "6+"},
		{"Барби", 114, "12+"},
		{"Оставь мир позади", 138, "16+"},
	}

	for _, m := range movies {
		_, err := db.Exec("INSERT INTO movies (title, duration, age_rating) VALUES (?, ?, ?)",
			m.title, m.duration, m.ageRating)
		if err != nil {
			return err
		}
	}

	// Добавляем тестовые данные сеансов
	screenings := []struct {
		movieID    int
		screenTime string
		hallNumber int
		price      float64
	}{
		{1, "2024-03-20 18:00:00", 1, 450.00},
		{1, "2024-03-20 21:00:00", 2, 500.00},
		{2, "2024-03-20 17:30:00", 3, 400.00},
		{3, "2024-03-20 15:00:00", 1, 350.00},
		{4, "2024-03-20 19:30:00", 2, 380.00},
	}

	for _, s := range screenings {
		_, err := db.Exec("INSERT INTO screenings (movie_id, screen_time, hall_number, price) VALUES (?, ?, ?, ?)",
			s.movieID, s.screenTime, s.hallNumber, s.price)
		if err != nil {
			return err
		}
	}

	// Добавляем тестовые данные клиентов
	customers := []struct {
		name  string
		email string
		phone string
	}{
		{"Аюров Жамбал", "zhambal@mail.com", "+7-900-111-22-33"},
	}

	for _, c := range customers {
		_, err := db.Exec("INSERT INTO customers (name, email, phone) VALUES (?, ?, ?)",
			c.name, c.email, c.phone)
		if err != nil {
			return err
		}
	}

	log.Println("Тестовые данные добавлены")
	return nil
}
