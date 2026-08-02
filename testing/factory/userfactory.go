package factory

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/dakaii/vibegopher/internal/database"
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/repository/userrepo"
	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

// usernameSeq guarantees unique usernames across the process. faker.New() seeds
// from time.Now().Unix() (second resolution), so two users built within the same
// second would otherwise draw the same name and violate the username UNIQUE
// constraint.
var usernameSeq atomic.Uint64

func BuildUser() domain.User {
	fake := faker.New()
	return domain.User{
		ID:        uuid.New(),
		Username:  fmt.Sprintf("%s_%d", fake.Internet().User(), usernameSeq.Add(1)),
		Password:  "Password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func BuildUserWithPassword(password string) domain.User {
	user := BuildUser()
	user.Password = password
	return user
}

func BuildUsers(count int) []domain.User {
	users := make([]domain.User, count)
	for i := 0; i < count; i++ {
		users[i] = BuildUserWithPassword(fmt.Sprintf("Password%d", i+1))
	}
	return users
}

func CreateUser() domain.User {
	return createUserWithPassword("Password123")
}

func CreateUserWithPassword(password string) domain.User {
	return createUserWithPassword(password)
}

func CreateUsers(count int) []domain.User {
	users := make([]domain.User, count)
	for i := 0; i < count; i++ {
		users[i] = createUserWithPassword(fmt.Sprintf("Password%d", i+1))
	}
	return users
}

func createUserWithPassword(plainPassword string) domain.User {
	user := BuildUserWithPassword(plainPassword)
	repo := userrepo.NewUserRepo(database.GetDatabase(true))
	saved, err := repo.CreateUser(domain.User{
		Username: user.Username,
		Password: plainPassword,
	})
	if err != nil {
		panic(err)
	}
	saved.Password = plainPassword
	return *saved
}
