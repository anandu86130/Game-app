package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID   uint   `gorm:"primaryKey" json:"user_id"`
	Name     string `json:"name" gorm:"not null"`
	Email    string `gorm:"unique;not null" json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password" gorm:"not null"`
}

type OTP struct {
	gorm.Model
	Email string `json:"email"`
	Otp   string `json:"otp"`
	Exp   time.Time
}

type League struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	PrizePool float64   `json:"prize_pool"`
	Teams     []Team    `json:"teams"`
	StartTime time.Time `json:"start_time"`
}

type Team struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	PlayerID uint    `json:"player_id"`
	Score    float64 `json:"score"`
	LeagueID uint    `json:"league_id"` // Foreign key reference to League
}

type TeamA struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	Name         string  `json:"name" gorm:"not null"`
	PlayerID     uint    `json:"player_id"`
	Score        float64 `json:"score"`
	TournamentID uint    `json:"tournament_id"`
}

// TeamBID struct represents Team B in the tournament.
type TeamB struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	Name         string  `json:"name" gorm:"not null"`
	PlayerID     uint    `json:"player_id"`
	Score        float64 `json:"score"`
	TournamentID uint    `json:"tournament_id"`
}

// Tournament struct represents a tournament between Team A and Team B.
type Tournament struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	PrizePool float64   `json:"prize_pool"`
	TeamA     []TeamA   `json:"TeamA"`
	TeamB     []TeamB   `json:"TeamB"`
	StartTime time.Time `json:"start_time"`
}
