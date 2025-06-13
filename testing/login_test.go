package testing

import (
	"bytes"
	"encoding/json"
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

type LoginTestSuite struct {
	suite.Suite
	router http.Handler
}

func (suite *LoginTestSuite) SetupTest() {
	db := database.GetDatabase(true)
	repos := repository.InitRepositories(db)
	controllers := controller.InitControllers(repos)
	suite.router = api.SetupRouter(controllers)
}

func (suite *LoginTestSuite) TearDownTest() {
	TruncateAllTables()
}

func (suite *LoginTestSuite) TestLoginUser() {
	// Create a user first
	user := factory.CreateUser()

	// Prepare login request
	loginReq := map[string]string{
		"username": user.Username,
		"password": "password123", // Default password from factory
	}

	reqBody, err := json.Marshal(loginReq)
	suite.NoError(err)

	// Make login request
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var authToken domain.AuthToken
	err = json.Unmarshal(rr.Body.Bytes(), &authToken)
	suite.NoError(err)

	assert.Equal(suite.T(), "Bearer", authToken.TokenType)
	assert.NotEmpty(suite.T(), authToken.Token)
	assert.Greater(suite.T(), authToken.ExpiresIn, int64(0))
}

func (suite *LoginTestSuite) TestLoginWithInvalidCredentials() {
	// Create a user first
	user := factory.CreateUser()

	// Prepare login request with wrong password
	loginReq := map[string]string{
		"username": user.Username,
		"password": "wrongpassword",
	}

	reqBody, err := json.Marshal(loginReq)
	suite.NoError(err)

	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
}

func (suite *LoginTestSuite) TestLoginWithNonExistentUser() {
	// Prepare login request for non-existent user
	loginReq := map[string]string{
		"username": "nonexistent",
		"password": "password123",
	}

	reqBody, err := json.Marshal(loginReq)
	suite.NoError(err)

	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}
