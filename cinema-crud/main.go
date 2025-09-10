package main

import (
	"cinema-crud/database"
	"cinema-crud/menu"
	"cinema-crud/repository"
	"log"
)

func main() {
	// Инициализация базы данных
	db := database.InitDB()
	defer db.Close()

	// Создание репозитория
	repo := &repository.CinemaRepository{DB: db}

	// Запуск главного меню
	menu.ShowMainMenu(repo)

	log.Println("Программа завершена")
}
