package testing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dakaii/vibegopher/internal/api"
	"github.com/dakaii/vibegopher/internal/auth"
	"github.com/dakaii/vibegopher/internal/controller"
	"github.com/dakaii/vibegopher/internal/database"
	"github.com/dakaii/vibegopher/internal/repository"
	"github.com/dakaii/vibegopher/testing/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MeTestSuite struct {
	suite.Suite
	router http.Handler
}

type MeResponse struct {
	Username string `json:"username"`
}

func (suite *MeTestSuite) SetupTest() {
	db := database.GetDatabase(true)
	repos := repository.InitRepositories(db)
	controllers := controller.InitControllers(repos)
	suite.router = api.SetupRouter(controllers)
}

func (suite *MeTestSuite) TearDownTest() {
	TruncateAllTables()
}

func (suite *MeTestSuite) TestGetMe() {
	// Create a user and get auth token
	user := factory.CreateUser()
	token, err := auth.GenerateJWT(user)
	suite.NoError(err)

	// Make authenticated request to /me endpoint
	req, err := http.NewRequest("GET", "/api/me", nil)
	suite.NoError(err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var meResponse MeResponse
	err = json.Unmarshal(rr.Body.Bytes(), &meResponse)
	suite.NoError(err)

	assert.Equal(suite.T(), user.Username, meResponse.Username)
}

func (suite *MeTestSuite) TestGetMeWithoutAuth() {
	// Make request without authorization header
	req, err := http.NewRequest("GET", "/api/me", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	// Should return unauthorized
	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
}

func (suite *MeTestSuite) TestGetMeWithInvalidToken() {
	// Make request with invalid token
	req, err := http.NewRequest("GET", "/api/me", nil)
	suite.NoError(err)
	req.Header.Set("Authorization", "Bearer invalid-token")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	// Should return unauthorized
	assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
}

func TestMeTestSuite(t *testing.T) {
	suite.Run(t, new(MeTestSuite))
}
