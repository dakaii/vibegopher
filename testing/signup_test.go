package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dakaii/graphyy/internal/api"
	"github.com/dakaii/graphyy/internal/controller"
	"github.com/dakaii/graphyy/internal/database"
	"github.com/dakaii/graphyy/internal/domain"
	"github.com/dakaii/graphyy/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SignUpTestSuite struct {
	suite.Suite
	router http.Handler
}

func (suite *SignUpTestSuite) SetupTest() {
	db := database.GetDatabase(true)
	repos := repository.InitRepositories(db)
	controllers := controller.InitControllers(repos)
	suite.router = api.SetupRouter(controllers)
}

func (suite *SignUpTestSuite) TearDownTest() {
	TruncateAllTables()
}

func (suite *SignUpTestSuite) TestCreateUser() {
	// Prepare request
	signupReq := map[string]string{
		"username": "testuser",
		"password": "password123",
	}

	reqBody, err := json.Marshal(signupReq)
	suite.NoError(err)

	// Make request
	req, err := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusCreated, rr.Code)

	var authToken domain.AuthToken
	err = json.Unmarshal(rr.Body.Bytes(), &authToken)
	suite.NoError(err)

	assert.Equal(suite.T(), "Bearer", authToken.TokenType)
	assert.NotEmpty(suite.T(), authToken.Token)
	assert.Greater(suite.T(), authToken.ExpiresIn, int64(0))
}

func (suite *SignUpTestSuite) TestCreateUserWithInvalidUsername() {
	// Test with short username (should fail validation)
	signupReq := map[string]string{
		"username": "short",
		"password": "password123",
	}

	reqBody, err := json.Marshal(signupReq)
	suite.NoError(err)

	req, err := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func TestSignUpTestSuite(t *testing.T) {
	suite.Run(t, new(SignUpTestSuite))
}
