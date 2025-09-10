package models

import "time"

type Movie struct {
	ID        int
	Title     string
	Duration  int
	AgeRating string
}

type Screening struct {
	ID         int
	MovieID    int
	ScreenTime time.Time
	HallNumber int
	Price      float64
}

type Ticket struct {
	ID           int
	ScreeningID  int
	SeatNumber   int
	PurchaseTime time.Time
	CustomerName string
}

type Customer struct {
	ID            int
	Name          string
	Email         string
	Phone         string
	LoyaltyPoints int
}
