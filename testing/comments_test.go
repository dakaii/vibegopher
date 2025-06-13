package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dakaii/vibegopher/internal/api"
	"github.com/dakaii/vibegopher/internal/controller"
	"github.com/dakaii/vibegopher/internal/database"
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/repository"
	"github.com/dakaii/vibegopher/testing/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommentTestSuite struct {
	suite.Suite
	router http.Handler
}

func (suite *CommentTestSuite) SetupTest() {
	db := database.GetDatabase(true)
	repos := repository.InitRepositories(db)
	controllers := controller.InitControllers(repos)
	suite.router = api.SetupRouter(controllers)
}

func (suite *CommentTestSuite) TearDownTest() {
	TruncateAllTables()
}

func (suite *CommentTestSuite) TestCreateComment() {
	// Create user and post first
	user := factory.CreateUser()
	post := factory.CreatePostWithContent(user.ID, "Test post content")

	// Prepare request
	commentReq := map[string]interface{}{
		"content": "This is a test comment",
		"post_id": post.ID,
	}

	reqBody, err := json.Marshal(commentReq)
	suite.NoError(err)

	// Make request with authorization
	req, err := http.NewRequest("POST", "/api/comments", bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", suite.getBearerToken(user))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusCreated, rr.Code)

	var comment domain.Comment
	err = json.Unmarshal(rr.Body.Bytes(), &comment)
	suite.NoError(err)

	assert.Equal(suite.T(), "This is a test comment", comment.Content)
	assert.Equal(suite.T(), user.ID, comment.UserID)
	assert.Equal(suite.T(), post.ID, comment.PostID)
	assert.Equal(suite.T(), user.Username, comment.User.Username)
	assert.NotEmpty(suite.T(), comment.ID)
}

func (suite *CommentTestSuite) TestCreateCommentUnauthorized() {
	user := factory.CreateUser()
	post := factory.CreatePostWithContent(user.ID, "Test post content")

	commentReq := map[string]interface{}{
		"content": "This is a test comment",
		"post_id": post.ID,
	}

	reqBody, err := json.Marshal(commentReq)
	suite.NoError(err)

	req, err := http.NewRequest("POST", "/api/comments", bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
}

func (suite *CommentTestSuite) TestCreateCommentInvalidContent() {
	user := factory.CreateUser()
	post := factory.CreatePostWithContent(user.ID, "Test post content")

	// Test empty content
	commentReq := map[string]interface{}{
		"content": "",
		"post_id": post.ID,
	}

	reqBody, err := json.Marshal(commentReq)
	suite.NoError(err)

	req, err := http.NewRequest("POST", "/api/comments", bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", suite.getBearerToken(user))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func (suite *CommentTestSuite) TestGetCommentsByPostID() {
	// Create users and post
	user1 := factory.CreateUser()
	user2 := factory.CreateUserWithPassword("password456")
	post := factory.CreatePostWithContent(user1.ID, "Test post content")

	// Create comments
	factory.CreateCommentWithContent(user1.ID, post.ID, "Comment 1 by user 1")
	factory.CreateCommentWithContent(user2.ID, post.ID, "Comment 2 by user 2")
	factory.CreateCommentWithContent(user1.ID, post.ID, "Comment 3 by user 1")

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/comments/post/%s", post.ID), nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var comments []domain.Comment
	err = json.Unmarshal(rr.Body.Bytes(), &comments)
	suite.NoError(err)

	assert.Len(suite.T(), comments, 3)
	for _, comment := range comments {
		assert.Equal(suite.T(), post.ID, comment.PostID)
	}
}

func (suite *CommentTestSuite) TestGetCommentsByUserID() {
	user1 := factory.CreateUser()
	user2 := factory.CreateUserWithPassword("password456")
	post1 := factory.CreatePostWithContent(user1.ID, "Post 1")
	post2 := factory.CreatePostWithContent(user2.ID, "Post 2")

	factory.CreateCommentWithContent(user1.ID, post1.ID, "Comment 1 by user 1")
	factory.CreateCommentWithContent(user1.ID, post2.ID, "Comment 2 by user 1")
	factory.CreateCommentWithContent(user2.ID, post1.ID, "Comment by user 2")

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/comments/user/%s", user1.ID), nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var comments []domain.Comment
	err = json.Unmarshal(rr.Body.Bytes(), &comments)
	suite.NoError(err)

	assert.Len(suite.T(), comments, 2)
	for _, comment := range comments {
		assert.Equal(suite.T(), user1.ID, comment.UserID)
	}
}

func (suite *CommentTestSuite) TestGetCommentByID() {
	user := factory.CreateUser()
	post := factory.CreatePostWithContent(user.ID, "Test post content")
	comment := factory.CreateCommentWithContent(user.ID, post.ID, "Test comment content")

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/comments/%s", comment.ID), nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var retrievedComment domain.Comment
	err = json.Unmarshal(rr.Body.Bytes(), &retrievedComment)
	suite.NoError(err)

	assert.Equal(suite.T(), comment.ID, retrievedComment.ID)
	assert.Equal(suite.T(), "Test comment content", retrievedComment.Content)
	assert.Equal(suite.T(), user.ID, retrievedComment.UserID)
	assert.Equal(suite.T(), post.ID, retrievedComment.PostID)
}

func (suite *CommentTestSuite) TestUpdateComment() {
	user := factory.CreateUser()
	post := factory.CreatePostWithContent(user.ID, "Test post content")
	comment := factory.CreateCommentWithContent(user.ID, post.ID, "Original comment")

	updateReq := map[string]string{
		"content": "Updated comment",
	}

	reqBody, err := json.Marshal(updateReq)
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/api/comments/%s", comment.ID), bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", suite.getBearerToken(user))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var updatedComment domain.Comment
	err = json.Unmarshal(rr.Body.Bytes(), &updatedComment)
	suite.NoError(err)

	assert.Equal(suite.T(), "Updated comment", updatedComment.Content)
	assert.Equal(suite.T(), comment.ID, updatedComment.ID)
}

func (suite *CommentTestSuite) TestUpdateCommentUnauthorized() {
	user1 := factory.CreateUser()
	user2 := factory.CreateUserWithPassword("password456")
	post := factory.CreatePostWithContent(user1.ID, "Test post content")
	comment := factory.CreateCommentWithContent(user1.ID, post.ID, "Original comment")

	updateReq := map[string]string{
		"content": "Updated comment",
	}

	reqBody, err := json.Marshal(updateReq)
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/api/comments/%s", comment.ID), bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", suite.getBearerToken(user2))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusForbidden, rr.Code)
}

func (suite *CommentTestSuite) TestDeleteComment() {
	user := factory.CreateUser()
	post := factory.CreatePostWithContent(user.ID, "Test post content")
	comment := factory.CreateCommentWithContent(user.ID, post.ID, "Comment to delete")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/comments/%s", comment.ID), nil)
	suite.NoError(err)
	req.Header.Set("Authorization", suite.getBearerToken(user))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusNoContent, rr.Code)

	// Verify comment is deleted
	getReq, err := http.NewRequest("GET", fmt.Sprintf("/api/comments/%s", comment.ID), nil)
	suite.NoError(err)

	getRr := httptest.NewRecorder()
	suite.router.ServeHTTP(getRr, getReq)

	assert.Equal(suite.T(), http.StatusNotFound, getRr.Code)
}

func (suite *CommentTestSuite) TestDeleteCommentUnauthorized() {
	user1 := factory.CreateUser()
	user2 := factory.CreateUserWithPassword("password456")
	post := factory.CreatePostWithContent(user1.ID, "Test post content")
	comment := factory.CreateCommentWithContent(user1.ID, post.ID, "Comment to delete")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/comments/%s", comment.ID), nil)
	suite.NoError(err)
	req.Header.Set("Authorization", suite.getBearerToken(user2))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusForbidden, rr.Code)
}

// Helper method to get bearer token
func (suite *CommentTestSuite) getBearerToken(user domain.User) string {
	signupReq := map[string]string{
		"username": user.Username,
		"password": "password123",
	}

	reqBody, _ := json.Marshal(signupReq)
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	var authToken domain.AuthToken
	json.Unmarshal(rr.Body.Bytes(), &authToken)

	return fmt.Sprintf("Bearer %s", authToken.Token)
}

func TestCommentTestSuite(t *testing.T) {
	suite.Run(t, new(CommentTestSuite))
}
