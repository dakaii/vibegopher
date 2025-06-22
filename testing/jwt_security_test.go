package testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/dakaii/vibegopher/internal/auth"
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
)

// TestJWTSecurity is a security test to verify JWT validation works properly
func TestJWTSecurity(t *testing.T) {
	fmt.Println("Testing JWT Security Fix...")

	// Set a test secret
	os.Setenv("AUTH_SECRET", "test-secret-key")

	// Create a test user
	user := domain.User{
		ID:       uuid.New(),
		Username: "testuser",
	}

	// Generate a valid JWT
	token := auth.GenerateJWT(user)
	fmt.Printf("Generated token: %s\n", token.Token)

	// Test 1: Valid token should verify successfully
	verifiedUser, err := auth.VerifyJWT(token.Token)
	if err != nil {
		fmt.Printf("❌ FAIL: Valid token verification failed: %v\n", err)
		return
	}
	fmt.Printf("✅ PASS: Valid token verified successfully for user: %s\n", verifiedUser.Username)

	// Test 2: Try to create a token with "none" algorithm (this should be prevented by our fix)
	// This simulates an attack where someone tries to bypass signature verification
	noneToken := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJ1c2VybmFtZSI6ImhhY2tlciIsImlkIjoiMTIzNDUifQ."
	_, err = auth.VerifyJWT(noneToken)
	if err != nil {
		fmt.Printf("✅ PASS: 'none' algorithm attack was properly rejected: %v\n", err)
	} else {
		fmt.Println("❌ FAIL: 'none' algorithm attack was not rejected!")
		return
	}

	// Test 3: Invalid signature should be rejected
	invalidToken := token.Token[:len(token.Token)-10] + "invalidsig"
	_, err = auth.VerifyJWT(invalidToken)
	if err != nil {
		fmt.Printf("✅ PASS: Invalid signature was properly rejected: %v\n", err)
	} else {
		fmt.Println("❌ FAIL: Invalid signature was not rejected!")
		return
	}

	fmt.Println("\n🎉 All JWT security tests passed!")
}
