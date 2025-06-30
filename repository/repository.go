// my-go-app/repository/repository.go
package repository

import (
	"fmt"
	"log"
	"my-go-app/model" // Adjust module path if different
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Repository struct holds the database connection
type Repository struct {
	DB *gorm.DB
}

// NewRepository initializes a new database connection and returns a Repository instance
func NewRepository() (*Repository, error) {
	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		return nil, fmt.Errorf("database environment variables (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT) must be set")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Yangon",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the User model to create the users table
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate database: %w", err)
	}
	log.Println("Database migration completed successfully.")

	return &Repository{DB: db}, nil
}

// CreateUser saves a new user to the database
func (r *Repository) CreateUser(user *model.User) error {
	result := r.DB.Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

// GetUserByUsername finds a user by their username
func (r *Repository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	result := r.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by username: %w", result.Error)
	}
	return &user, nil
}

// Close closes the database connection
func (r *Repository) Close() error {
	sqlDB, err := r.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}
