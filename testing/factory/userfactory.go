package factory

import (
	"fmt"
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/repository/userrepo"
	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

// BuildUser creates a user entity with faker data
func BuildUser() domain.User {
	fake := faker.New()
	return domain.User{
		ID:        uuid.New(),
		Username:  fake.Internet().User(),
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// BuildUserWithPassword creates a user entity with a specific password
func BuildUserWithPassword(password string) domain.User {
	user := BuildUser()
	user.Password = password
	return user
}

// BuildUsers creates multiple user entities with faker data
func BuildUsers(count int) []domain.User {
	users := make([]domain.User, count)
	for i := 0; i < count; i++ {
		users[i] = BuildUserWithPassword(fmt.Sprintf("password%d", i+1))
	}
	return users
}

// CreateUser creates and saves a single user to the database
func CreateUser() domain.User {
	return createUserWithPassword("password123")
}

// CreateUserWithPassword creates and saves a user with a specific password
func CreateUserWithPassword(password string) domain.User {
	return createUserWithPassword(password)
}

// CreateUsers creates and saves multiple users to the database
func CreateUsers(count int) []domain.User {
	users := make([]domain.User, count)
	for i := 0; i < count; i++ {
		users[i] = createUserWithPassword(fmt.Sprintf("password%d", i+1))
	}
	return users
}

// Helper function to create and save a user with password handling
func createUserWithPassword(plainPassword string) domain.User {
	user := BuildUserWithPassword(plainPassword)

	// Hash password for database
	hashedPassword, err := userrepo.HashPassword(user.Password)
	if err != nil {
		panic(err)
	}
	user.Password = hashedPassword

	// Save to database
	savedUser, err := Save(user)
	if err != nil {
		panic(err)
	}

	// Return with plain password for tests
	savedUser.Password = plainPassword
	return savedUser
}
