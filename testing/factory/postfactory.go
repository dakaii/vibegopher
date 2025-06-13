package factory

import (
	"fmt"
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

// BuildPost creates a post entity with faker data
func BuildPost(userID uuid.UUID) domain.Post {
	fake := faker.New()
	return domain.Post{
		ID:        uuid.New(),
		Content:   fake.Lorem().Sentence(10),
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// BuildPostWithContent creates a post entity with specific content
func BuildPostWithContent(userID uuid.UUID, content string) domain.Post {
	post := BuildPost(userID)
	post.Content = content
	return post
}

// BuildPosts creates multiple post entities with faker data
func BuildPosts(userID uuid.UUID, count int) []domain.Post {
	posts := make([]domain.Post, count)
	for i := 0; i < count; i++ {
		posts[i] = BuildPostWithContent(userID, fmt.Sprintf("Test post content %d", i+1))
	}
	return posts
}

// CreatePost creates and saves a single post to the database
func CreatePost(userID uuid.UUID) domain.Post {
	return createPostWithContent(userID, "Test post content")
}

// CreatePostWithContent creates and saves a post with specific content
func CreatePostWithContent(userID uuid.UUID, content string) domain.Post {
	return createPostWithContent(userID, content)
}

// CreatePosts creates and saves multiple posts to the database
func CreatePosts(userID uuid.UUID, count int) []domain.Post {
	posts := make([]domain.Post, count)
	for i := 0; i < count; i++ {
		posts[i] = createPostWithContent(userID, fmt.Sprintf("Test post content %d", i+1))
	}
	return posts
}

// Helper function to create and save a post
func createPostWithContent(userID uuid.UUID, content string) domain.Post {
	post := BuildPostWithContent(userID, content)

	// Save to database
	savedPost, err := Save(post)
	if err != nil {
		panic(err)
	}

	return savedPost
}
