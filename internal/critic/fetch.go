package critic

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	maxRedirects   = 3
	maxSnippetBody = 64 << 10 // 64 KiB
	maxSnippetText = 500
)

// FetchURLSnippet downloads a public http(s) URL and returns a short text snippet.
// It blocks private/link-local/metadata addresses (including at dial time) to reduce SSRF risk.
func FetchURLSnippet(ctx context.Context, rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid url")
	}
	if err := validatePublicURL(parsed); err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 8 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}
				ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
				if err != nil {
					return nil, fmt.Errorf("dns lookup failed")
				}
				var lastErr error
				for _, ipa := range ips {
					if !isPublicIP(ipa.IP) {
						lastErr = fmt.Errorf("blocked address")
						continue
					}
					d := net.Dialer{Timeout: 5 * time.Second}
					conn, err := d.DialContext(ctx, network, net.JoinHostPort(ipa.IP.String(), port))
					if err == nil {
						return conn, nil
					}
					lastErr = err
				}
				if lastErr == nil {
					lastErr = fmt.Errorf("no addresses")
				}
				return nil, lastErr
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("too many redirects")
			}
			return validatePublicURL(req.URL)
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "VibeGopherCritic/1.0 (+https://github.com/dakaii/vibegopher)")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	ctype := strings.ToLower(resp.Header.Get("Content-Type"))
	if ctype != "" && !strings.Contains(ctype, "text/") && !strings.Contains(ctype, "json") && !strings.Contains(ctype, "xml") && !strings.Contains(ctype, "html") {
		return "", fmt.Errorf("unsupported content type")
	}

	limited := io.LimitReader(resp.Body, maxSnippetBody)
	buf, err := io.ReadAll(limited)
	if err != nil {
		return "", err
	}

	text := stripTags(string(buf))
	text = strings.Join(strings.Fields(text), " ")
	if len(text) > maxSnippetText {
		text = text[:maxSnippetText] + "..."
	}
	if text == "" {
		return "(empty body)", nil
	}
	return text, nil
}

func validatePublicURL(u *url.URL) error {
	if u == nil {
		return fmt.Errorf("invalid url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported scheme")
	}
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("missing host")
	}
	lower := strings.ToLower(host)
	if lower == "localhost" || strings.HasSuffix(lower, ".localhost") || lower == "metadata.google.internal" {
		return fmt.Errorf("blocked host")
	}

	// Fast-path literal IPs (no DNS). Hostnames are re-checked at dial time.
	if ip := net.ParseIP(host); ip != nil {
		if !isPublicIP(ip) {
			return fmt.Errorf("blocked address")
		}
		return nil
	}
	return nil
}

func isPublicIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() || ip.IsUnspecified() {
		return false
	}
	// Cloud metadata ranges commonly abused for SSRF.
	if ip4 := ip.To4(); ip4 != nil {
		if ip4[0] == 169 && ip4[1] == 254 {
			return false
		}
	}
	return true
}

func stripTags(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return b.String()
}
