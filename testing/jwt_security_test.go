package testing

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/dakaii/vibegopher/internal/auth"
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
)

func TestJWTSecurity(t *testing.T) {
	fmt.Println("Testing JWT Security Fix...")

	os.Setenv("AUTH_SECRET", "test-secret-key")

	user := domain.User{
		ID:       uuid.New(),
		Username: "testuser",
		Password: "super-secret-hash",
	}

	token, err := auth.GenerateJWT(user)
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}
	fmt.Printf("Generated token: %s\n", token.Token)

	verifiedUser, err := auth.VerifyJWT(token.Token)
	if err != nil {
		t.Fatalf("Valid token verification failed: %v", err)
	}
	fmt.Printf("✅ PASS: Valid token verified successfully for user: %s\n", verifiedUser.Username)

	payload := jwtPayload(t, token.Token)
	if strings.Contains(payload, "Password") || strings.Contains(payload, "super-secret-hash") || strings.Contains(payload, "password") {
		t.Fatalf("JWT payload embeds password material: %s", payload)
	}
	fmt.Println("✅ PASS: JWT does not embed password claims")

	noneToken := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJ1c2VybmFtZSI6ImhhY2tlciIsImlkIjoiMTIzNDUifQ."
	if _, err = auth.VerifyJWT(noneToken); err == nil {
		t.Fatal("'none' algorithm attack was not rejected")
	}
	fmt.Printf("✅ PASS: 'none' algorithm attack was properly rejected: %v\n", err)

	invalidToken := token.Token[:len(token.Token)-10] + "invalidsig"
	if _, err = auth.VerifyJWT(invalidToken); err == nil {
		t.Fatal("Invalid signature was not rejected")
	}
	fmt.Printf("✅ PASS: Invalid signature was properly rejected: %v\n", err)

	fmt.Println("\n🎉 All JWT security tests passed!")
}

func jwtPayload(t *testing.T, token string) string {
	t.Helper()
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		t.Fatalf("malformed token")
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	return string(raw)
}
