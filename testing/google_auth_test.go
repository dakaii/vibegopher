package testing

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dakaii/vibegopher/internal/api"
	"github.com/dakaii/vibegopher/internal/auth"
	"github.com/dakaii/vibegopher/internal/controller"
	"github.com/dakaii/vibegopher/internal/database"
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type fakeGoogleValidator struct {
	profile *domain.GoogleProfile
	err     error
}

func (f fakeGoogleValidator) Validate(ctx context.Context, rawToken, audience string) (*domain.GoogleProfile, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.profile, nil
}

type GoogleAuthTestSuite struct {
	suite.Suite
	router http.Handler
}

func (s *GoogleAuthTestSuite) SetupTest() {
	db := database.GetDatabase(true)
	repos := repository.InitRepositories(db)
	controllers := controller.InitControllers(repos)
	controllers.UserController.WithGoogleValidator(fakeGoogleValidator{
		profile: &domain.GoogleProfile{
			Sub:           "google-sub-123",
			Email:         "alice@example.com",
			EmailVerified: true,
			GivenName:     "Alice",
		},
	})
	s.router = api.SetupRouter(controllers)
}

func (s *GoogleAuthTestSuite) TearDownTest() {
	TruncateAllTables()
}

func (s *GoogleAuthTestSuite) TestGoogleAuthCreatesUserAndReturnsJWT() {
	body, _ := json.Marshal(map[string]string{"id_token": "fake-google-token"})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/google", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)

	require.Equal(s.T(), http.StatusOK, rr.Code)

	var token domain.AuthToken
	require.NoError(s.T(), json.Unmarshal(rr.Body.Bytes(), &token))
	require.Equal(s.T(), "Bearer", token.TokenType)
	require.NotEmpty(s.T(), token.Token)

	user, err := auth.VerifyJWT(token.Token)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "alice", user.Username)
	require.NotContains(s.T(), token.Token, "Password")
}

func TestGoogleAuthTestSuite(t *testing.T) {
	suite.Run(t, new(GoogleAuthTestSuite))
}
