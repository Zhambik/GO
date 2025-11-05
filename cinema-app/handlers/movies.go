package handlers

import (
	"cinema-app/database"
	"cinema-app/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetMoviesPublic - получение фильмов без авторизации
func GetMoviesPublic(c *gin.Context) {
	log.Println("GetMoviesPublic called")

	query := "SELECT id, title, genre, release_date, rating, director FROM movies WHERE 1=1"
	args := []interface{}{}

	// Фильтрация
	if title := c.Query("title"); title != "" {
		query += " AND title LIKE ?"
		args = append(args, "%"+title+"%")
	}
	if genre := c.Query("genre"); genre != "" {
		query += " AND genre = ?"
		args = append(args, genre)
	}
	if director := c.Query("director"); director != "" {
		query += " AND director LIKE ?"
		args = append(args, "%"+director+"%")
	}

	// Сортировка
	if sort := c.Query("sort"); sort != "" {
		order := "ASC"
		if c.Query("order") == "desc" {
			order = "DESC"
		}
		query += " ORDER BY " + sort + " " + order
	}

	log.Printf("Executing query: %s with args: %v", query, args)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Genre, &m.ReleaseDate, &m.Rating, &m.Director)
		if err != nil {
			log.Printf("Scan error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		movies = append(movies, m)
	}

	log.Printf("Returning %d movies", len(movies))
	c.JSON(http.StatusOK, movies)
}

// GetMovies - получение фильмов с авторизацией
func GetMovies(c *gin.Context) {
	log.Println("GetMovies (protected) called")

	query := "SELECT id, title, genre, release_date, rating, director FROM movies WHERE 1=1"
	args := []interface{}{}

	// Фильтрация
	if title := c.Query("title"); title != "" {
		query += " AND title LIKE ?"
		args = append(args, "%"+title+"%")
	}
	if genre := c.Query("genre"); genre != "" {
		query += " AND genre = ?"
		args = append(args, genre)
	}
	if director := c.Query("director"); director != "" {
		query += " AND director LIKE ?"
		args = append(args, "%"+director+"%")
	}

	// Сортировка
	if sort := c.Query("sort"); sort != "" {
		order := "ASC"
		if c.Query("order") == "desc" {
			order = "DESC"
		}
		query += " ORDER BY " + sort + " " + order
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Genre, &m.ReleaseDate, &m.Rating, &m.Director)
		if err != nil {
			log.Printf("Scan error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		movies = append(movies, m)
	}

	c.JSON(http.StatusOK, movies)
}

// MovieRequest - структура для входящих запросов (с датой как строка)
type MovieRequest struct {
	Title       string  `json:"title"`
	Genre       string  `json:"genre"`
	ReleaseDate string  `json:"release_date"`
	Rating      float64 `json:"rating"`
	Director    string  `json:"director"`
}

func CreateMovie(c *gin.Context) {
	log.Println("CreateMovie called")

	var movieReq MovieRequest
	if err := c.ShouldBindJSON(&movieReq); err != nil {
		log.Printf("Bind JSON error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем строку даты в time.Time
	releaseDate, err := time.Parse("2006-01-02", movieReq.ReleaseDate)
	if err != nil {
		log.Printf("Date parsing error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	log.Printf("Creating movie: %+v", movieReq)

	result, err := database.DB.Exec(
		"INSERT INTO movies (title, genre, release_date, rating, director) VALUES (?, ?, ?, ?, ?)",
		movieReq.Title, movieReq.Genre, releaseDate, movieReq.Rating, movieReq.Director)
	if err != nil {
		log.Printf("Error creating movie: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	// Создаем объект Movie для ответа
	movie := models.Movie{
		ID:          int(id),
		Title:       movieReq.Title,
		Genre:       movieReq.Genre,
		ReleaseDate: releaseDate,
		Rating:      movieReq.Rating,
		Director:    movieReq.Director,
	}

	log.Printf("Movie created with ID: %d", movie.ID)
	c.JSON(http.StatusCreated, movie)
}

func UpdateMovie(c *gin.Context) {
	id := c.Param("id")
	log.Printf("UpdateMovie called for ID: %s", id)

	var movieReq MovieRequest
	if err := c.ShouldBindJSON(&movieReq); err != nil {
		log.Printf("Bind JSON error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем строку даты в time.Time
	releaseDate, err := time.Parse("2006-01-02", movieReq.ReleaseDate)
	if err != nil {
		log.Printf("Date parsing error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	log.Printf("Updating movie ID %s with data: %+v", id, movieReq)

	result, err := database.DB.Exec(
		"UPDATE movies SET title=?, genre=?, release_date=?, rating=?, director=? WHERE id=?",
		movieReq.Title, movieReq.Genre, releaseDate, movieReq.Rating, movieReq.Director, id)
	if err != nil {
		log.Printf("Error updating movie: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, была ли обновлена хотя бы одна строка
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
	} else {
		log.Printf("Rows affected: %d", rowsAffected)
		if rowsAffected == 0 {
			log.Printf("No movie found with ID: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
	}

	movieID, _ := strconv.Atoi(id)
	movie := models.Movie{
		ID:          movieID,
		Title:       movieReq.Title,
		Genre:       movieReq.Genre,
		ReleaseDate: releaseDate,
		Rating:      movieReq.Rating,
		Director:    movieReq.Director,
	}

	log.Printf("Movie updated successfully: %+v", movie)
	c.JSON(http.StatusOK, movie)
}

func DeleteMovie(c *gin.Context) {
	id := c.Param("id")
	log.Printf("DeleteMovie called for ID: %s", id)

	_, err := database.DB.Exec("DELETE FROM movies WHERE id=?", id)
	if err != nil {
		log.Printf("Error deleting movie: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Movie deleted successfully: ID %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted"})
}
