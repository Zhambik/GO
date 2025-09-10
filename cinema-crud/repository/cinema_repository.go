package repository

import (
	"cinema-crud/models"
	"database/sql"
	"fmt"
	"time"
)

type CinemaRepository struct {
	DB *sql.DB
}

// ==================== MOVIES CRUD ====================

func (r *CinemaRepository) CreateMovie(movie *models.Movie) error {
	result, err := r.DB.Exec("INSERT INTO movies (title, duration, age_rating) VALUES (?, ?, ?)",
		movie.Title, movie.Duration, movie.AgeRating)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	movie.ID = int(id)
	return nil
}

func (r *CinemaRepository) GetAllMovies() ([]models.Movie, error) {
	rows, err := r.DB.Query("SELECT id, title, duration, age_rating FROM movies ORDER BY title")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Duration, &m.AgeRating); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *CinemaRepository) GetMovieByID(id int) (models.Movie, error) {
	var movie models.Movie
	err := r.DB.QueryRow("SELECT id, title, duration, age_rating FROM movies WHERE id = ?", id).
		Scan(&movie.ID, &movie.Title, &movie.Duration, &movie.AgeRating)
	return movie, err
}

func (r *CinemaRepository) UpdateMovie(movie *models.Movie) error {
	result, err := r.DB.Exec("UPDATE movies SET title = ?, duration = ?, age_rating = ? WHERE id = ?",
		movie.Title, movie.Duration, movie.AgeRating, movie.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("фильм с ID %d не найден", movie.ID)
	}
	return nil
}

func (r *CinemaRepository) DeleteMovie(id int) error {
	result, err := r.DB.Exec("DELETE FROM movies WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("фильм с ID %d не найден", id)
	}
	return nil
}

// ==================== SCREENINGS CRUD ====================

func (r *CinemaRepository) CreateScreening(screening *models.Screening) error {
	result, err := r.DB.Exec("INSERT INTO screenings (movie_id, screen_time, hall_number, price) VALUES (?, ?, ?, ?)",
		screening.MovieID, screening.ScreenTime, screening.HallNumber, screening.Price)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	screening.ID = int(id)
	return nil
}

func (r *CinemaRepository) GetAllScreenings() ([]models.Screening, error) {
	rows, err := r.DB.Query(`
        SELECT s.id, s.movie_id, s.screen_time, s.hall_number, s.price 
        FROM screenings s 
        ORDER BY s.screen_time
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var screenings []models.Screening
	for rows.Next() {
		var s models.Screening
		if err := rows.Scan(&s.ID, &s.MovieID, &s.ScreenTime, &s.HallNumber, &s.Price); err != nil {
			return nil, err
		}
		screenings = append(screenings, s)
	}
	return screenings, nil
}

func (r *CinemaRepository) GetScreeningsByMovie(movieID int) ([]models.Screening, error) {
	rows, err := r.DB.Query(`
        SELECT s.id, s.movie_id, s.screen_time, s.hall_number, s.price 
        FROM screenings s 
        WHERE s.movie_id = ? 
        ORDER BY s.screen_time
    `, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var screenings []models.Screening
	for rows.Next() {
		var s models.Screening
		if err := rows.Scan(&s.ID, &s.MovieID, &s.ScreenTime, &s.HallNumber, &s.Price); err != nil {
			return nil, err
		}
		screenings = append(screenings, s)
	}
	return screenings, nil
}

func (r *CinemaRepository) DeleteScreening(id int) error {
	result, err := r.DB.Exec("DELETE FROM screenings WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("сеанс с ID %d не найден", id)
	}
	return nil
}

// ==================== TICKETS CRUD ====================

func (r *CinemaRepository) BuyTicket(ticket *models.Ticket) error {
	// Проверяем, свободно ли место
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tickets WHERE screening_id = ? AND seat_number = ?)",
		ticket.ScreeningID, ticket.SeatNumber).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("место %d уже занято", ticket.SeatNumber)
	}

	result, err := r.DB.Exec("INSERT INTO tickets (screening_id, seat_number, customer_name) VALUES (?, ?, ?)",
		ticket.ScreeningID, ticket.SeatNumber, ticket.CustomerName)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	ticket.ID = int(id)
	ticket.PurchaseTime = time.Now()
	return nil
}

func (r *CinemaRepository) GetTicketsByScreening(screeningID int) ([]models.Ticket, error) {
	rows, err := r.DB.Query(`
        SELECT id, screening_id, seat_number, purchase_time, customer_name 
        FROM tickets 
        WHERE screening_id = ? 
        ORDER BY seat_number
    `, screeningID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.ScreeningID, &t.SeatNumber, &t.PurchaseTime, &t.CustomerName); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

func (r *CinemaRepository) ReturnTicket(id int) error {
	result, err := r.DB.Exec("DELETE FROM tickets WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("билет с ID %d не найден", id)
	}
	return nil
}

// ==================== CUSTOMERS CRUD ====================

func (r *CinemaRepository) CreateCustomer(customer *models.Customer) error {
	result, err := r.DB.Exec("INSERT INTO customers (name, email, phone) VALUES (?, ?, ?)",
		customer.Name, customer.Email, customer.Phone)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	customer.ID = int(id)
	return nil
}

func (r *CinemaRepository) GetAllCustomers() ([]models.Customer, error) {
	rows, err := r.DB.Query("SELECT id, name, email, phone, loyalty_points FROM customers ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.LoyaltyPoints); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (r *CinemaRepository) GetCustomerByEmail(email string) (models.Customer, error) {
	var customer models.Customer
	err := r.DB.QueryRow("SELECT id, name, email, phone, loyalty_points FROM customers WHERE email = ?", email).
		Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.LoyaltyPoints)
	return customer, err
}

func (r *CinemaRepository) UpdateCustomer(customer *models.Customer) error {
	result, err := r.DB.Exec("UPDATE customers SET name = ?, phone = ?, loyalty_points = ? WHERE id = ?",
		customer.Name, customer.Phone, customer.LoyaltyPoints, customer.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("клиент с ID %d не найден", customer.ID)
	}
	return nil
}

func (r *CinemaRepository) DeleteCustomer(id int) error {
	result, err := r.DB.Exec("DELETE FROM customers WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("клиент с ID %d не найден", id)
	}
	return nil
}

// ==================== UTILITY METHODS ====================

func (r *CinemaRepository) GetAvailableSeats(screeningID int) ([]int, error) {
	// Предположим, что в зале 50 мест
	allSeats := make([]int, 50)
	for i := 0; i < 50; i++ {
		allSeats[i] = i + 1
	}

	// Получаем занятые места
	rows, err := r.DB.Query("SELECT seat_number FROM tickets WHERE screening_id = ?", screeningID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	takenSeats := make(map[int]bool)
	for rows.Next() {
		var seat int
		if err := rows.Scan(&seat); err != nil {
			return nil, err
		}
		takenSeats[seat] = true
	}

	// Фильтруем свободные места
	var availableSeats []int
	for _, seat := range allSeats {
		if !takenSeats[seat] {
			availableSeats = append(availableSeats, seat)
		}
	}

	return availableSeats, nil
}
