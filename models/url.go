package models

import (
	"time"
)

type URL struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ShortKey  string    `json:"short_key" gorm:"uniqueIndex;not null"`
	LongURL   string    `json:"long_url" gorm:"not null;type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	Clicks    int       `json:"clicks" gorm:"default:0"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	PasskeyHash string  `json:"-" gorm:"type:varchar(255)"`
}

type AuthToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}