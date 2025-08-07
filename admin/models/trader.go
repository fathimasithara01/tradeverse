package models

import "gorm.io/gorm"

type Trader struct {
	gorm.Model
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"-"`
	IsBanned  bool   `json:"is_banned" gorm:"default:false"`
	Bio       string
	TotalPnL  float64
	Followers []Follower
	Signals   []Signal
}
