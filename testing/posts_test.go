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

type PostTestSuite struct {
	suite.Suite
	router http.Handler
}

func (suite *PostTestSuite) SetupTest() {
	db := database.GetDatabase(true)
	repos := repository.InitRepositories(db)
	controllers := controller.InitControllers(repos)
	suite.router = api.SetupRouter(controllers)
}

func (suite *PostTestSuite) TearDownTest() {
	TruncateAllTables()
}

func (suite *PostTestSuite) TestCreatePost() {
	// Create a user first
	user := factory.CreateUser()

	// Prepare request
	postReq := map[string]string{
		"content": "This is a test post",
	}

	reqBody, err := json.Marshal(postReq)
	suite.NoError(err)

	// Make request with authorization
	req, err := http.NewRequest("POST", "/api/posts", bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", suite.getBearerToken(user))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusCreated, rr.Code)

	var post domain.Post
	err = json.Unmarshal(rr.Body.Bytes(), &post)
	suite.NoError(err)

	assert.Equal(suite.T(), "This is a test post", post.Content)
	assert.Equal(suite.T(), user.ID, post.UserID)
	assert.Equal(suite.T(), user.Username, post.User.Username)
	assert.NotEmpty(suite.T(), post.ID)
}

func (suite *PostTestSuite) TestCreatePostUnauthorized() {
	postReq := map[string]string{
		"content": "This is a test post",
	}

	reqBody, err := json.Marshal(postReq)
	suite.NoError(err)

	req, err := http.NewRequest("POST", "/api/posts", bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
}

func (suite *PostTestSuite) TestCreatePostInvalidContent() {
	user := factory.CreateUser()

	// Test empty content
	postReq := map[string]string{
		"content": "",
	}

	reqBody, err := json.Marshal(postReq)
	suite.NoError(err)

	req, err := http.NewRequest("POST", "/api/posts", bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", suite.getBearerToken(user))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func (suite *PostTestSuite) TestGetAllPosts() {
	// Create users and posts
	user1 := factory.CreateUser()
	user2 := factory.CreateUserWithPassword("password456")
	factory.CreatePostWithContent(user1.ID, "Post by user 1")
	factory.CreatePostWithContent(user2.ID, "Post by user 2")

	req, err := http.NewRequest("GET", "/api/posts", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var posts []domain.Post
	err = json.Unmarshal(rr.Body.Bytes(), &posts)
	suite.NoError(err)

	assert.Len(suite.T(), posts, 2)
}

func (suite *PostTestSuite) TestGetPostByID() {
	user := factory.CreateUser()
	post := factory.CreatePostWithContent(user.ID, "Test post content")

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/posts/%s", post.ID), nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var retrievedPost domain.Post
	err = json.Unmarshal(rr.Body.Bytes(), &retrievedPost)
	suite.NoError(err)

	assert.Equal(suite.T(), post.ID, retrievedPost.ID)
	assert.Equal(suite.T(), "Test post content", retrievedPost.Content)
}

func (suite *PostTestSuite) TestGetPostsByUserID() {
	user1 := factory.CreateUser()
	user2 := factory.CreateUserWithPassword("password456")
	factory.CreatePostWithContent(user1.ID, "Post 1 by user 1")
	factory.CreatePostWithContent(user1.ID, "Post 2 by user 1")
	factory.CreatePostWithContent(user2.ID, "Post by user 2")

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/posts/user/%s", user1.ID), nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var posts []domain.Post
	err = json.Unmarshal(rr.Body.Bytes(), &posts)
	suite.NoError(err)

	assert.Len(suite.T(), posts, 2)
	for _, post := range posts {
		assert.Equal(suite.T(), user1.ID, post.UserID)
	}
}

func (suite *PostTestSuite) TestUpdatePost() {
	user := factory.CreateUser()
	post := factory.CreatePostWithContent(user.ID, "Original content")

	updateReq := map[string]string{
		"content": "Updated content",
	}

	reqBody, err := json.Marshal(updateReq)
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/api/posts/%s", post.ID), bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", suite.getBearerToken(user))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var updatedPost domain.Post
	err = json.Unmarshal(rr.Body.Bytes(), &updatedPost)
	suite.NoError(err)

	assert.Equal(suite.T(), "Updated content", updatedPost.Content)
	assert.Equal(suite.T(), post.ID, updatedPost.ID)
}

func (suite *PostTestSuite) TestUpdatePostUnauthorized() {
	user1 := factory.CreateUser()
	user2 := factory.CreateUserWithPassword("password456")
	post := factory.CreatePostWithContent(user1.ID, "Original content")

	updateReq := map[string]string{
		"content": "Updated content",
	}

	reqBody, err := json.Marshal(updateReq)
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/api/posts/%s", post.ID), bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", suite.getBearerToken(user2))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusForbidden, rr.Code)
}

func (suite *PostTestSuite) TestDeletePost() {
	user := factory.CreateUser()
	post := factory.CreatePostWithContent(user.ID, "Content to delete")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/posts/%s", post.ID), nil)
	suite.NoError(err)
	req.Header.Set("Authorization", suite.getBearerToken(user))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusNoContent, rr.Code)

	// Verify post is deleted
	getReq, err := http.NewRequest("GET", fmt.Sprintf("/api/posts/%s", post.ID), nil)
	suite.NoError(err)

	getRr := httptest.NewRecorder()
	suite.router.ServeHTTP(getRr, getReq)

	assert.Equal(suite.T(), http.StatusNotFound, getRr.Code)
}

func (suite *PostTestSuite) TestDeletePostUnauthorized() {
	user1 := factory.CreateUser()
	user2 := factory.CreateUserWithPassword("password456")
	post := factory.CreatePostWithContent(user1.ID, "Content to delete")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/posts/%s", post.ID), nil)
	suite.NoError(err)
	req.Header.Set("Authorization", suite.getBearerToken(user2))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusForbidden, rr.Code)
}

// Helper method to get bearer token
func (suite *PostTestSuite) getBearerToken(user domain.User) string {
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

func TestPostTestSuite(t *testing.T) {
	suite.Run(t, new(PostTestSuite))
}
