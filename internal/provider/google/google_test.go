package google

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/emoral435/time-broker/internal/config"
)

func TestName(t *testing.T) {
	p := New()
	if got := p.Name(); got != ProviderName {
		t.Errorf("Name() = %q; want %q", got, ProviderName)
	}
}

func TestNewWithEnvVars(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "test-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "test-secret")

	p := New()
	if p == nil {
		t.Fatal("New() returned nil")
	}
	if p.config == nil {
		t.Fatal("New() didn't initialize config")
	}
	if p.config.ClientID != "test-id" {
		t.Errorf("ClientID = %q; want %q", p.config.ClientID, "test-id")
	}
	if p.config.ClientSecret != "test-secret" {
		t.Errorf("ClientSecret = %q; want %q", p.config.ClientSecret, "test-secret")
	}
}

func TestNewWithBuildFlags(t *testing.T) {
	ClientID = "build-id"
	ClientSecret = "build-secret"
	defer func() {
		ClientID = ""
		ClientSecret = ""
	}()

	os.Unsetenv("GOOGLE_CLIENT_ID")
	os.Unsetenv("GOOGLE_CLIENT_SECRET")

	p := New()
	if p.config.ClientID != "build-id" {
		t.Errorf("ClientID = %q; want %q", p.config.ClientID, "build-id")
	}
	if p.config.ClientSecret != "build-secret" {
		t.Errorf("ClientSecret = %q; want %q", p.config.ClientSecret, "build-secret")
	}
}

func TestNewBuildFlagsOverrideEnvVars(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "env-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "env-secret")

	ClientID = "build-id"
	ClientSecret = "build-secret"
	defer func() {
		ClientID = ""
		ClientSecret = ""
	}()

	p := New()
	if p.config.ClientID != "build-id" {
		t.Errorf("ClientID = %q; want 'build-id' (build flags take priority)", p.config.ClientID)
	}
	if p.config.ClientSecret != "build-secret" {
		t.Errorf("ClientSecret = %q; want 'build-secret' (build flags take priority)", p.config.ClientSecret)
	}
}

func TestNewTokenLoadedFromFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("GOOGLE_CLIENT_ID", "test-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "test-secret")

	tok := &oauth2.Token{
		AccessToken: "cached-access-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Hour),
	}
	if err := config.SaveToken(tok); err != nil {
		t.Fatalf("SaveToken() error: %v", err)
	}

	p := New()
	if p.token == nil {
		t.Fatal("token should be loaded from file")
	}
	if p.token.AccessToken != "cached-access-token" {
		t.Errorf("AccessToken = %q; want %q", p.token.AccessToken, "cached-access-token")
	}
	if p.service != nil {
		t.Log("service was created from cached token (requires network)")
	}
}

func TestFreeSlotsNotAuthenticated(t *testing.T) {
	p := &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     "test-id",
			ClientSecret: "test-secret",
		},
	}
	_, err := p.FreeSlots(time.Now(), time.Hour)
	if err == nil {
		t.Fatal("FreeSlots() expected error when not authenticated")
	}
	if !strings.Contains(err.Error(), "not authenticated") {
		t.Errorf("FreeSlots() error = %q; want 'not authenticated'", err)
	}
}

func TestBookNotAuthenticated(t *testing.T) {
	p := &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     "test-id",
			ClientSecret: "test-secret",
		},
	}
	err := p.Book("test", time.Now(), time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("Book() expected error when not authenticated")
	}
	if !strings.Contains(err.Error(), "not authenticated") {
		t.Errorf("Book() error = %q; want 'not authenticated'", err)
	}
}

func TestFreeSlotsNotYetImplemented(t *testing.T) {
	p := &GoogleProvider{
		service: &calendar.Service{},
	}
	_, err := p.FreeSlots(time.Now(), time.Hour)
	if err == nil {
		t.Fatal("FreeSlots() expected 'not yet implemented'")
	}
	if !strings.Contains(err.Error(), "not yet implemented") {
		t.Errorf("FreeSlots() error = %q; want 'not yet implemented'", err)
	}
}

func TestBookNotYetImplemented(t *testing.T) {
	p := &GoogleProvider{
		service: &calendar.Service{},
	}
	err := p.Book("test", time.Now(), time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("Book() expected 'not yet implemented'")
	}
	if !strings.Contains(err.Error(), "not yet implemented") {
		t.Errorf("Book() error = %q; want 'not yet implemented'", err)
	}
}

func TestRandomState(t *testing.T) {
	s1 := randomState()
	s2 := randomState()

	if len(s1) != 32 {
		t.Errorf("randomState() length = %d; want 32", len(s1))
	}
	if s1 == s2 {
		t.Error("randomState() returned duplicate values")
	}
}

func TestRandomStateHex(t *testing.T) {
	s := randomState()
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("randomState() contains non-hex character: %c", c)
		}
	}
}

func TestCreateServiceSadPath_NilToken(t *testing.T) {
	p := &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     "test-id",
			ClientSecret: "test-secret",
			Endpoint:     google.Endpoint,
		},
		token: nil,
	}
	p.createService()
	if p.service == nil {
		t.Fatal("service should be created; calendar.NewService is lazy and doesn't authenticate eagerly")
	}
}

func TestCreateServiceSadPath_ExpiredToken(t *testing.T) {
	p := &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     "test-id",
			ClientSecret: "test-secret",
			Endpoint: oauth2.Endpoint{
				TokenURL: "https://invalid-token.example.com/token",
			},
		},
		token: &oauth2.Token{
			AccessToken:  "expired-token",
			TokenType:    "Bearer",
			RefreshToken: "expired-refresh",
			Expiry:       time.Now().Add(-1 * time.Hour),
		},
	}
	p.createService()
	if p.service == nil {
		t.Fatal("service should be created; calendar.NewService is lazy and doesn't eagerly authenticate")
	}
}

func TestCreateServiceHappyPath(t *testing.T) {
	hit := false
	var mockServer *httptest.Server
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit = true
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"kind":    "discovery#restDescription",
			"name":    "calendar",
			"version": "v3",
			"baseUrl": mockServer.URL + "/calendar/v3/",
		})
	}))
	defer mockServer.Close()

	p := &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     "test-id",
			ClientSecret: "test-secret",
			Endpoint: oauth2.Endpoint{
				TokenURL: mockServer.URL + "/token",
			},
		},
		token: &oauth2.Token{
			AccessToken: "valid-token",
			TokenType:   "Bearer",
			Expiry:      time.Now().Add(time.Hour),
		},
	}

	p.createService(option.WithEndpoint(mockServer.URL))
	if p.service == nil {
		t.Fatal("service should be created successfully with valid token and mock endpoint")
	}
	if !hit {
		t.Log("note: calendar.NewService did not contact the mock endpoint (lazy initialization)")
	}
}

func TestAuth_AlreadyAuthenticated(t *testing.T) {
	p := &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     "test-id",
			ClientSecret: "test-secret",
		},
		token: &oauth2.Token{
			AccessToken: "valid-token",
			TokenType:   "Bearer",
			Expiry:      time.Now().Add(time.Hour),
		},
	}
	if err := p.Auth(); err != nil {
		t.Errorf("Auth() should return nil when already authenticated, got: %v", err)
	}
}

func TestAuth_MissingCredentials(t *testing.T) {
	p := &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     "",
			ClientSecret: "",
		},
	}
	err := p.Auth()
	if err == nil {
		t.Fatal("Auth() should return error when credentials are missing")
	}
	if !strings.Contains(err.Error(), "not configured") {
		t.Errorf("Auth() error = %q; want 'not configured'", err)
	}
}

func TestAuth_FullFlowSuccess(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	freePort := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	oldPort := redirectPort
	redirectPort = fmt.Sprintf(":%d", freePort)
	defer func() { redirectPort = oldPort }()

	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "mock-access-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	}))
	defer tokenServer.Close()

	authCalled := make(chan struct{})
	oldOpen := openURL
	openURL = func(rawurl string) error {
		close(authCalled)

		go func() {
			time.Sleep(50 * time.Millisecond)
			u, err := url.Parse(rawurl)
			if err != nil {
				t.Errorf("failed to parse auth URL: %v", err)
				return
			}
			state := u.Query().Get("state")
			callbackURL := fmt.Sprintf("http://localhost:%d/callback?code=test-code&state=%s", freePort, state)
			resp, err := http.Get(callbackURL)
			if err != nil {
				t.Errorf("callback request failed: %v", err)
				return
			}
			resp.Body.Close()
		}()
		return nil
	}
	defer func() { openURL = oldOpen }()

	p := &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			Endpoint: oauth2.Endpoint{
				AuthURL:  tokenServer.URL + "/auth",
				TokenURL: tokenServer.URL + "/token",
			},
			RedirectURL: fmt.Sprintf("http://localhost:%d/callback", freePort),
			Scopes:      []string{calendar.CalendarScope},
		},
	}

	if err := p.Auth(); err != nil {
		t.Fatalf("Auth() failed: %v", err)
	}
	if p.token == nil {
		t.Fatal("token should be set after successful auth")
	}
	if p.token.AccessToken != "mock-access-token" {
		t.Errorf("AccessToken = %q; want %q", p.token.AccessToken, "mock-access-token")
	}
}

func TestAuth_FullFlowCallbackError(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	freePort := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	oldPort := redirectPort
	redirectPort = fmt.Sprintf(":%d", freePort)
	defer func() { redirectPort = oldPort }()

	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "mock-access-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	}))
	defer tokenServer.Close()

	oldOpen := openURL
	openURL = func(rawurl string) error {
		go func() {
			time.Sleep(50 * time.Millisecond)
			callbackURL := fmt.Sprintf("http://localhost:%d/callback", freePort)
			resp, err := http.Get(callbackURL)
			if err != nil {
				t.Errorf("callback request failed: %v", err)
				return
			}
			resp.Body.Close()
		}()
		return nil
	}
	defer func() { openURL = oldOpen }()

	p := &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			Endpoint: oauth2.Endpoint{
				AuthURL:  tokenServer.URL + "/auth",
				TokenURL: tokenServer.URL + "/token",
			},
			RedirectURL: fmt.Sprintf("http://localhost:%d/callback", freePort),
			Scopes:      []string{calendar.CalendarScope},
		},
	}

	if err := p.Auth(); err == nil {
		t.Fatal("Auth() should return error when callback has no code")
	} else if !strings.Contains(err.Error(), "no code in callback") {
		t.Errorf("Auth() error = %q; want 'no code in callback'", err)
	}
}

func TestSaveTokenSource(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	ts := &saveTokenSource{
		src: oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: "saved-token",
			TokenType:   "Bearer",
			Expiry:      time.Now().Add(time.Hour),
		}),
	}

	tok, err := ts.Token()
	if err != nil {
		t.Fatalf("Token() error: %v", err)
	}
	if tok.AccessToken != "saved-token" {
		t.Errorf("AccessToken = %q; want %q", tok.AccessToken, "saved-token")
	}
}
