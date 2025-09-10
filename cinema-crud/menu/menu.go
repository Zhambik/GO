package menu

import (
	"bufio"
	"cinema-crud/models"
	"cinema-crud/repository"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ShowMainMenu(repo *repository.CinemaRepository) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n=== ГЛАВНОЕ МЕНЮ КИНОТЕАТРА ===")
		fmt.Println("1. Управление фильмами")
		fmt.Println("2. Управление сеансами")
		fmt.Println("3. Продажа билетов")
		fmt.Println("4. Управление клиентами")
		fmt.Println("5. Просмотр всех данных")
		fmt.Println("6. Выйти из программы")
		fmt.Print("Выберите действие (1-6): ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			showMoviesMenu(repo, scanner)
		case "2":
			showScreeningsMenu(repo, scanner)
		case "3":
			showTicketsMenu(repo, scanner)
		case "4":
			showCustomersMenu(repo, scanner)
		case "5":
			showAllData(repo)
		case "6":
			fmt.Println("Выход из программы...")
			return
		default:
			fmt.Println("Неверный выбор. Пожалуйста, выберите от 1 до 6.")
		}

		fmt.Print("\nНажмите Enter для продолжения...")
		scanner.Scan()
	}
}

// ==================== МЕНЮ ФИЛЬМОВ ====================
func showMoviesMenu(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	for {
		fmt.Println("\n=== МЕНЮ ФИЛЬМОВ ===")
		fmt.Println("1. Показать все фильмы")
		fmt.Println("2. Добавить новый фильм")
		fmt.Println("3. Обновить фильм")
		fmt.Println("4. Удалить фильм")
		fmt.Println("5. Назад в главное меню")
		fmt.Print("Выберите действие (1-5): ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			showAllMovies(repo)
		case "2":
			addMovie(repo, scanner)
		case "3":
			updateMovie(repo, scanner)
		case "4":
			deleteMovie(repo, scanner)
		case "5":
			return
		default:
			fmt.Println("Неверный выбор.")
		}
	}
}

func showAllMovies(repo *repository.CinemaRepository) {
	movies, err := repo.GetAllMovies()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Println("\nСписок всех фильмов:")
	fmt.Println("----------------------------------")
	for _, movie := range movies {
		fmt.Printf("ID: %d | %s | %d мин | %s\n",
			movie.ID, movie.Title, movie.Duration, movie.AgeRating)
	}
	fmt.Printf("Всего фильмов: %d\n", len(movies))
}

func addMovie(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Println("\nДобавление нового фильма:")

	var movie models.Movie

	fmt.Print("Введите название фильма: ")
	scanner.Scan()
	movie.Title = strings.TrimSpace(scanner.Text())

	fmt.Print("Введите продолжительность (в минутах): ")
	scanner.Scan()
	durationStr := strings.TrimSpace(scanner.Text())
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		fmt.Println("Ошибка: продолжительность должна быть числом")
		return
	}
	movie.Duration = duration

	fmt.Print("Введите возрастной рейтинг (0+, 6+, 12+, 16+, 18+): ")
	scanner.Scan()
	movie.AgeRating = strings.TrimSpace(scanner.Text())

	if err := repo.CreateMovie(&movie); err != nil {
		fmt.Printf("Ошибка при добавлении фильма: %v\n", err)
		return
	}

	fmt.Printf("Фильм '%s' успешно добавлен с ID %d\n", movie.Title, movie.ID)
}

func updateMovie(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Print("Введите ID фильма для обновления: ")
	scanner.Scan()
	idStr := strings.TrimSpace(scanner.Text())
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		return
	}

	movie, err := repo.GetMovieByID(id)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("Текущие данные: %s, %d мин, %s\n", movie.Title, movie.Duration, movie.AgeRating)

	fmt.Printf("Введите новое название (Enter - оставить текущее): ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())
	if title != "" {
		movie.Title = title
	}

	fmt.Printf("Введите новую продолжительность (Enter - оставить текущую): ")
	scanner.Scan()
	durationStr := strings.TrimSpace(scanner.Text())
	if durationStr != "" {
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			fmt.Println("Ошибка: продолжительность должна быть числом")
			return
		}
		movie.Duration = duration
	}

	fmt.Printf("Введите новый возрастной рейтинг (Enter - оставить текущий): ")
	scanner.Scan()
	ageRating := strings.TrimSpace(scanner.Text())
	if ageRating != "" {
		movie.AgeRating = ageRating
	}

	if err := repo.UpdateMovie(&movie); err != nil {
		fmt.Printf("Ошибка при обновлении: %v\n", err)
		return
	}

	fmt.Printf("Фильм с ID %d успешно обновлен\n", movie.ID)
}

func deleteMovie(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Print("Введите ID фильма для удаления: ")
	scanner.Scan()
	idStr := strings.TrimSpace(scanner.Text())
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		return
	}

	movie, err := repo.GetMovieByID(id)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("Удалить фильм '%s'? (y/n): ", movie.Title)
	scanner.Scan()
	confirmation := strings.TrimSpace(strings.ToLower(scanner.Text()))

	if confirmation != "y" && confirmation != "yes" {
		fmt.Println("Удаление отменено")
		return
	}

	if err := repo.DeleteMovie(id); err != nil {
		fmt.Printf("Ошибка при удалении: %v\n", err)
		return
	}

	fmt.Printf("Фильм '%s' успешно удален\n", movie.Title)
}

// ==================== МЕНЮ СЕАНСОВ ====================
func showScreeningsMenu(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	for {
		fmt.Println("\n=== МЕНЮ СЕАНСОВ ===")
		fmt.Println("1. Показать все сеансы")
		fmt.Println("2. Показать сеансы по фильму")
		fmt.Println("3. Добавить новый сеанс")
		fmt.Println("4. Удалить сеанс")
		fmt.Println("5. Назад в главное меню")
		fmt.Print("Выберите действие (1-5): ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			showAllScreenings(repo)
		case "2":
			showScreeningsByMovie(repo, scanner)
		case "3":
			addScreening(repo, scanner)
		case "4":
			deleteScreening(repo, scanner)
		case "5":
			return
		default:
			fmt.Println("Неверный выбор.")
		}
	}
}

func showAllScreenings(repo *repository.CinemaRepository) {
	screenings, err := repo.GetAllScreenings()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Println("\nВсе сеансы:")
	fmt.Println("----------------------------------")
	for _, screening := range screenings {
		fmt.Printf("ID: %d | Фильм ID: %d | %s | Зал: %d | %.2f руб.\n",
			screening.ID, screening.MovieID, screening.ScreenTime.Format("02.01.2006 15:04"),
			screening.HallNumber, screening.Price)
	}
}

func showScreeningsByMovie(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Print("Введите ID фильма: ")
	scanner.Scan()
	movieIDStr := strings.TrimSpace(scanner.Text())
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		return
	}

	screenings, err := repo.GetScreeningsByMovie(movieID)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("\nСеансы для фильма ID %d:\n", movieID)
	fmt.Println("----------------------------------")
	for _, screening := range screenings {
		fmt.Printf("ID: %d | %s | Зал: %d | %.2f руб.\n",
			screening.ID, screening.ScreenTime.Format("02.01.2006 15:04"),
			screening.HallNumber, screening.Price)
	}
}

func addScreening(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Println("\nДобавление нового сеанса:")

	var screening models.Screening

	fmt.Print("Введите ID фильма: ")
	scanner.Scan()
	movieIDStr := strings.TrimSpace(scanner.Text())
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		return
	}
	screening.MovieID = movieID

	fmt.Print("Введите дату и время (формат: 2006-01-02 15:04): ")
	scanner.Scan()
	timeStr := strings.TrimSpace(scanner.Text())
	screenTime, err := time.Parse("2006-01-02 15:04", timeStr)
	if err != nil {
		fmt.Println("Ошибка: неверный формат даты")
		return
	}
	screening.ScreenTime = screenTime

	fmt.Print("Введите номер зала: ")
	scanner.Scan()
	hallStr := strings.TrimSpace(scanner.Text())
	hall, err := strconv.Atoi(hallStr)
	if err != nil {
		fmt.Println("Ошибка: номер зала должен быть числом")
		return
	}
	screening.HallNumber = hall

	fmt.Print("Введите цену билета: ")
	scanner.Scan()
	priceStr := strings.TrimSpace(scanner.Text())
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		fmt.Println("Ошибка: цена должна быть числом")
		return
	}
	screening.Price = price

	if err := repo.CreateScreening(&screening); err != nil {
		fmt.Printf("Ошибка при добавлении сеанса: %v\n", err)
		return
	}

	fmt.Printf("Сеанс успешно добавлен с ID %d\n", screening.ID)
}

func deleteScreening(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Print("Введите ID сеанса для удаления: ")
	scanner.Scan()
	idStr := strings.TrimSpace(scanner.Text())
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		return
	}

	fmt.Print("Удалить сеанс? (y/n): ")
	scanner.Scan()
	confirmation := strings.TrimSpace(strings.ToLower(scanner.Text()))

	if confirmation != "y" && confirmation != "yes" {
		fmt.Println("Удаление отменено")
		return
	}

	if err := repo.DeleteScreening(id); err != nil {
		fmt.Printf("Ошибка при удалении: %v\n", err)
		return
	}

	fmt.Printf("Сеанс с ID %d успешно удален\n", id)
}

// ==================== МЕНЮ БИЛЕТОВ ====================
func showTicketsMenu(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	for {
		fmt.Println("\n=== МЕНЮ БИЛЕТОВ ===")
		fmt.Println("1. Купить билет")
		fmt.Println("2. Показать билеты на сеанс")
		fmt.Println("3. Показать свободные места")
		fmt.Println("4. Вернуть билет")
		fmt.Println("5. Назад в главное меню")
		fmt.Print("Выберите действие (1-5): ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			buyTicket(repo, scanner)
		case "2":
			showTicketsForScreening(repo, scanner)
		case "3":
			showAvailableSeats(repo, scanner)
		case "4":
			returnTicket(repo, scanner)
		case "5":
			return
		default:
			fmt.Println("Неверный выбор.")
		}
	}
}

func buyTicket(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Println("\nПокупка билета:")

	var ticket models.Ticket

	fmt.Print("Введите ID сеанса: ")
	scanner.Scan()
	screeningIDStr := strings.TrimSpace(scanner.Text())
	screeningID, err := strconv.Atoi(screeningIDStr)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		return
	}
	ticket.ScreeningID = screeningID

	availableSeats, err := repo.GetAvailableSeats(screeningID)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("Свободные места: %v\n", availableSeats)

	fmt.Print("Введите номер места: ")
	scanner.Scan()
	seatStr := strings.TrimSpace(scanner.Text())
	seat, err := strconv.Atoi(seatStr)
	if err != nil {
		fmt.Println("Ошибка: номер места должен быть числом")
		return
	}
	ticket.SeatNumber = seat

	fmt.Print("Введите имя покупателя: ")
	scanner.Scan()
	ticket.CustomerName = strings.TrimSpace(scanner.Text())

	if err := repo.BuyTicket(&ticket); err != nil {
		fmt.Printf("Ошибка при покупке билета: %v\n", err)
		return
	}

	fmt.Printf("Билет успешно куплен с ID %d\n", ticket.ID)
}

func showTicketsForScreening(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Print("Введите ID сеанса: ")
	scanner.Scan()
	screeningIDStr := strings.TrimSpace(scanner.Text())
	screeningID, err := strconv.Atoi(screeningIDStr)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		return
	}

	tickets, err := repo.GetTicketsByScreening(screeningID)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("\nБилеты на сеанс %d:\n", screeningID)
	fmt.Println("----------------------------------")
	for _, ticket := range tickets {
		fmt.Printf("ID: %d | Место: %d | Покупатель: %s\n",
			ticket.ID, ticket.SeatNumber, ticket.CustomerName)
	}
	fmt.Printf("Всего билетов: %d\n", len(tickets))
}

func showAvailableSeats(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Print("Введите ID сеанса: ")
	scanner.Scan()
	screeningIDStr := strings.TrimSpace(scanner.Text())
	screeningID, err := strconv.Atoi(screeningIDStr)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		return
	}

	availableSeats, err := repo.GetAvailableSeats(screeningID)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("\nСвободные места на сеанс %d:\n", screeningID)
	fmt.Println("----------------------------------")
	fmt.Printf("Места: %v\n", availableSeats)
	fmt.Printf("Всего свободных мест: %d\n", len(availableSeats))
}

func returnTicket(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Print("Введите ID билета для возврата: ")
	scanner.Scan()
	idStr := strings.TrimSpace(scanner.Text())
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		return
	}

	fmt.Print("Вернуть билет? (y/n): ")
	scanner.Scan()
	confirmation := strings.TrimSpace(strings.ToLower(scanner.Text()))

	if confirmation != "y" && confirmation != "yes" {
		fmt.Println("Возврат отменен")
		return
	}

	if err := repo.ReturnTicket(id); err != nil {
		fmt.Printf("Ошибка при возврате: %v\n", err)
		return
	}

	fmt.Printf("Билет с ID %d успешно возвращен\n", id)
}

// ==================== МЕНЮ КЛИЕНТОВ ====================
func showCustomersMenu(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	for {
		fmt.Println("\n=== МЕНЮ КЛИЕНТОВ ===")
		fmt.Println("1. Показать всех клиентов")
		fmt.Println("2. Добавить нового клиента")
		fmt.Println("3. Обновить клиента")
		fmt.Println("4. Удалить клиента")
		fmt.Println("5. Назад в главное меню")
		fmt.Print("Выберите действие (1-5): ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			showAllCustomers(repo)
		case "2":
			addCustomer(repo, scanner)
		case "3":
			updateCustomer(repo, scanner)
		case "4":
			deleteCustomer(repo, scanner)
		case "5":
			return
		default:
			fmt.Println("Неверный выбор.")
		}
	}
}

func showAllCustomers(repo *repository.CinemaRepository) {
	customers, err := repo.GetAllCustomers()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Println("\nСписок всех клиентов:")
	fmt.Println("----------------------------------")
	for _, customer := range customers {
		fmt.Printf("ID: %d | %s | %s | %s | Баллы: %d\n",
			customer.ID, customer.Name, customer.Email, customer.Phone, customer.LoyaltyPoints)
	}
	fmt.Printf("Всего клиентов: %d\n", len(customers))
}

func addCustomer(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Println("\nДобавление нового клиента:")

	var customer models.Customer

	fmt.Print("Введите имя клиента: ")
	scanner.Scan()
	customer.Name = strings.TrimSpace(scanner.Text())

	fmt.Print("Введите email: ")
	scanner.Scan()
	customer.Email = strings.TrimSpace(scanner.Text())

	fmt.Print("Введите телефон: ")
	scanner.Scan()
	customer.Phone = strings.TrimSpace(scanner.Text())

	if err := repo.CreateCustomer(&customer); err != nil {
		fmt.Printf("Ошибка при добавлении клиента: %v\n", err)
		return
	}

	fmt.Printf("Клиент '%s' успешно добавлен с ID %d\n", customer.Name, customer.ID)
}

func updateCustomer(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Print("Введите email клиента для обновления: ")
	scanner.Scan()
	email := strings.TrimSpace(scanner.Text())

	customer, err := repo.GetCustomerByEmail(email)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("Текущие данные: %s, %s, баллы: %d\n",
		customer.Name, customer.Phone, customer.LoyaltyPoints)

	fmt.Printf("Введите новое имя (Enter - оставить текущее): ")
	scanner.Scan()
	name := strings.TrimSpace(scanner.Text())
	if name != "" {
		customer.Name = name
	}

	fmt.Printf("Введите новый телефон (Enter - оставить текущий): ")
	scanner.Scan()
	phone := strings.TrimSpace(scanner.Text())
	if phone != "" {
		customer.Phone = phone
	}

	fmt.Printf("Введите новые баллы лояльности (Enter - оставить текущие): ")
	scanner.Scan()
	pointsStr := strings.TrimSpace(scanner.Text())
	if pointsStr != "" {
		points, err := strconv.Atoi(pointsStr)
		if err != nil {
			fmt.Println("Ошибка: баллы должны быть числом")
			return
		}
		customer.LoyaltyPoints = points
	}

	if err := repo.UpdateCustomer(&customer); err != nil {
		fmt.Printf("Ошибка при обновлении: %v\n", err)
		return
	}

	fmt.Printf("Клиент с email %s успешно обновлен\n", customer.Email)
}

func deleteCustomer(repo *repository.CinemaRepository, scanner *bufio.Scanner) {
	fmt.Print("Введите email клиента для удаления: ")
	scanner.Scan()
	email := strings.TrimSpace(scanner.Text())

	customer, err := repo.GetCustomerByEmail(email)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("Удалить клиента %s? (y/n): ", customer.Name)
	scanner.Scan()
	confirmation := strings.TrimSpace(strings.ToLower(scanner.Text()))

	if confirmation != "y" && confirmation != "yes" {
		fmt.Println("Удаление отменено")
		return
	}

	if err := repo.DeleteCustomer(customer.ID); err != nil {
		fmt.Printf("Ошибка при удалении: %v\n", err)
		return
	}

	fmt.Printf("Клиент '%s' успешно удален\n", customer.Name)
}

// ==================== ПРОСМОТР ВСЕХ ДАННЫХ ====================
func showAllData(repo *repository.CinemaRepository) {
	fmt.Println("\n=== ВСЕ ДАННЫЕ КИНОТЕАТРА ===")

	movies, _ := repo.GetAllMovies()
	fmt.Println("\nФИЛЬМЫ:")
	fmt.Println("----------------------------------")
	for _, m := range movies {
		fmt.Printf("ID: %d | %s | %d мин | %s\n", m.ID, m.Title, m.Duration, m.AgeRating)
	}

	screenings, _ := repo.GetAllScreenings()
	fmt.Println("\nСЕАНСЫ:")
	fmt.Println("----------------------------------")
	for _, s := range screenings {
		fmt.Printf("ID: %d | Фильм ID: %d | %s | Зал: %d | %.2f руб.\n",
			s.ID, s.MovieID, s.ScreenTime.Format("02.01.2006 15:04"),
			s.HallNumber, s.Price)
	}

	customers, _ := repo.GetAllCustomers()
	fmt.Println("\nКЛИЕНТЫ:")
	fmt.Println("----------------------------------")
	for _, c := range customers {
		fmt.Printf("ID: %d | %s | %s | Баллы: %d\n",
			c.ID, c.Name, c.Email, c.LoyaltyPoints)
	}

	fmt.Printf("\nИтого: %d фильмов, %d сеансов, %d клиентов\n",
		len(movies), len(screenings), len(customers))
}
