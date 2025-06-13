package factory

import (
	"fmt"
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

// BuildComment creates a comment entity with faker data
func BuildComment(userID, postID uuid.UUID) domain.Comment {
	fake := faker.New()
	return domain.Comment{
		ID:        uuid.New(),
		Content:   fake.Lorem().Sentence(5),
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// BuildCommentWithContent creates a comment entity with specific content
func BuildCommentWithContent(userID, postID uuid.UUID, content string) domain.Comment {
	comment := BuildComment(userID, postID)
	comment.Content = content
	return comment
}

// BuildComments creates multiple comment entities with faker data
func BuildComments(userID, postID uuid.UUID, count int) []domain.Comment {
	comments := make([]domain.Comment, count)
	for i := 0; i < count; i++ {
		comments[i] = BuildCommentWithContent(userID, postID, fmt.Sprintf("Test comment content %d", i+1))
	}
	return comments
}

// CreateComment creates and saves a single comment to the database
func CreateComment(userID, postID uuid.UUID) domain.Comment {
	return createCommentWithContent(userID, postID, "Test comment content")
}

// CreateCommentWithContent creates and saves a comment with specific content
func CreateCommentWithContent(userID, postID uuid.UUID, content string) domain.Comment {
	return createCommentWithContent(userID, postID, content)
}

// CreateComments creates and saves multiple comments to the database
func CreateComments(userID, postID uuid.UUID, count int) []domain.Comment {
	comments := make([]domain.Comment, count)
	for i := 0; i < count; i++ {
		comments[i] = createCommentWithContent(userID, postID, fmt.Sprintf("Test comment content %d", i+1))
	}
	return comments
}

// Helper function to create and save a comment
func createCommentWithContent(userID, postID uuid.UUID, content string) domain.Comment {
	comment := BuildCommentWithContent(userID, postID, content)

	// Save to database
	savedComment, err := Save(comment)
	if err != nil {
		panic(err)
	}

	return savedComment
}
