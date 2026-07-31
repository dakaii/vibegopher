package critic

import (
	"net"
	"net/url"
	"testing"
)

func TestValidatePublicURLBlocksLocal(t *testing.T) {
	cases := []string{
		"http://localhost/x",
		"http://127.0.0.1/x",
		"http://metadata.google.internal/",
		"file:///etc/passwd",
	}
	for _, raw := range cases {
		u, err := url.Parse(raw)
		if err != nil && raw != "file:///etc/passwd" {
			t.Fatalf("parse %s: %v", raw, err)
		}
		if raw == "file:///etc/passwd" {
			if err := validatePublicURL(u); err == nil {
				t.Fatalf("expected block for %s", raw)
			}
			continue
		}
		if err := validatePublicURL(u); err == nil {
			t.Fatalf("expected block for %s", raw)
		}
	}
}

func TestIsPublicIP(t *testing.T) {
	if isPublicIP(net.ParseIP("10.0.0.1")) {
		t.Fatal("private ip should be blocked")
	}
	if isPublicIP(net.ParseIP("169.254.169.254")) {
		t.Fatal("link local / metadata should be blocked")
	}
	if !isPublicIP(net.ParseIP("1.1.1.1")) {
		t.Fatal("public ip should be allowed")
	}
}
