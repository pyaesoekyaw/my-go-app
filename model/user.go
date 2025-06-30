// my-go-app/model/user.go
package model

import "time"

// User struct defines the user model for the database
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"` // Storing hashed password
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
