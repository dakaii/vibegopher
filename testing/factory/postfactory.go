package factory

import (
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

// Example Post entity (this would typically be in your domain package)
type Post struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  uuid.UUID `json:"author_id"`
}

// BuildPost creates a post entity with faker data without saving to database
func BuildPost() Post {
	fake := faker.New()

	return Post{
		ID:        uuid.New(),
		Title:     fake.Lorem().Sentence(5),
		Content:   fake.Lorem().Paragraph(3),
		AuthorID:  uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// BuildPostWithAuthor creates a post entity with a specific author ID
func BuildPostWithAuthor(authorID uuid.UUID) Post {
	fake := faker.New()

	return Post{
		ID:        uuid.New(),
		Title:     fake.Lorem().Sentence(5),
		Content:   fake.Lorem().Paragraph(3),
		AuthorID:  authorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// BuildPosts creates multiple post entities with faker data without saving to database
func BuildPosts(count int) []Post {
	posts := make([]Post, count)
	fake := faker.New()

	for i := 0; i < count; i++ {
		posts[i] = Post{
			ID:        uuid.New(),
			Title:     fake.Lorem().Sentence(5),
			Content:   fake.Lorem().Paragraph(3),
			AuthorID:  uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	return posts
}

// CreatePost creates and saves a single post to the database
func CreatePost() Post {
	post := BuildPost()

	savedPost, err := Save(post)
	if err != nil {
		panic(err)
	}

	return savedPost
}

// CreatePostWithAuthor creates and saves a post with a specific author
func CreatePostWithAuthor(authorID uuid.UUID) Post {
	post := BuildPostWithAuthor(authorID)

	savedPost, err := Save(post)
	if err != nil {
		panic(err)
	}

	return savedPost
}

// CreatePosts creates and saves multiple posts to the database
func CreatePosts(count int) []Post {
	posts := BuildPosts(count)

	savedPosts, err := SaveMany(posts)
	if err != nil {
		panic(err)
	}

	return savedPosts
}
