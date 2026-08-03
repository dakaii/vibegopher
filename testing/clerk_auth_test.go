package testing

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dakaii/vibegopher/internal/api"
	"github.com/dakaii/vibegopher/internal/controller"
	"github.com/dakaii/vibegopher/internal/database"
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type fakeClerkValidator struct {
	profile *domain.ClerkProfile
	err     error
}

func (f fakeClerkValidator) Validate(ctx context.Context, rawToken string) (*domain.ClerkProfile, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.profile, nil
}

func (f fakeClerkValidator) Enrich(ctx context.Context, clerkUserID string) (*domain.ClerkProfile, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.profile, nil
}

type ClerkAuthTestSuite struct {
	suite.Suite
	router http.Handler
}

func (s *ClerkAuthTestSuite) SetupTest() {
	os.Setenv("CLERK_SECRET_KEY", "sk_test_fake_for_suite")
	db := database.GetDatabase(true)
	repos := repository.InitRepositories(db)
	controllers := controller.InitControllers(repos)
	controllers.UserController.WithClerkValidator(fakeClerkValidator{
		profile: &domain.ClerkProfile{
			Sub:       "user_clerk_123",
			Email:     "alice@example.com",
			GivenName: "Alice",
		},
	})
	s.router = api.SetupRouter(controllers)
}

func (s *ClerkAuthTestSuite) TearDownTest() {
	os.Unsetenv("CLERK_SECRET_KEY")
	TruncateAllTables()
}

func (s *ClerkAuthTestSuite) TestClerkBearerCreatesUserOnMe() {
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.Header.Set("Authorization", "Bearer fake-clerk-session")
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)

	require.Equal(s.T(), http.StatusOK, rr.Code)

	var me map[string]any
	require.NoError(s.T(), json.Unmarshal(rr.Body.Bytes(), &me))
	require.Equal(s.T(), "alice", me["username"])
	require.Equal(s.T(), "alice@example.com", me["email"])
	require.NotEmpty(s.T(), me["id"])
}

func TestClerkAuthTestSuite(t *testing.T) {
	suite.Run(t, new(ClerkAuthTestSuite))
}
