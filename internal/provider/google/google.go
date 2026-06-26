package google

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/emoral435/time-broker/internal/config"
	"github.com/emoral435/time-broker/internal/provider"
)

var (
	ClientID     string
	ClientSecret string
)

const redirectPort = ":8085"

type Provider struct {
	config  *oauth2.Config
	token   *oauth2.Token
	service *calendar.Service
}

func New() *Provider {
	id := ClientID
	secret := ClientSecret
	if id == "" {
		id = os.Getenv("GOOGLE_CLIENT_ID")
	}
	if secret == "" {
		secret = os.Getenv("GOOGLE_CLIENT_SECRET")
	}

	cfg := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		RedirectURL:  "http://localhost" + redirectPort + "/callback",
		Scopes:       []string{calendar.CalendarScope},
		Endpoint:     google.Endpoint,
	}

	p := &Provider{config: cfg}

	tok, err := config.LoadToken()
	if err == nil && tok != nil {
		p.token = tok
		p.createService()
	}

	return p
}

func (g *Provider) Name() string {
	return "google"
}

func (g *Provider) Auth() error {
	if g.token != nil && g.token.Valid() {
		return nil
	}

	if g.config.ClientID == "" || g.config.ClientSecret == "" {
		return fmt.Errorf("google OAuth credentials not configured: set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables, or inject them at build time")
	}

	codeChan := make(chan string)
	errChan := make(chan error)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no code in callback response")
			return
		}
		fmt.Fprint(w, `<html><body><h1>Authentication successful!</h1><p>You can close this tab.</p></body></html>`)
		codeChan <- code
	})

	listener, err := net.Listen("tcp", redirectPort)
	if err != nil {
		return fmt.Errorf("start callback server on port %s: %w", redirectPort, err)
	}

	server := &http.Server{Handler: mux}
	go func() {
		err := server.Serve(listener)
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()
	defer server.Close()

	state := randomState()

	authURL := g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Println("Opening browser for Google authentication...")
	if err := openURL(authURL); err != nil {
		fmt.Printf("Open this URL in your browser: %s\n", authURL)
	}

	select {
	case code := <-codeChan:
		tok, err := g.config.Exchange(context.Background(), code)
		if err != nil {
			return fmt.Errorf("exchange auth code: %w", err)
		}
		g.token = tok
		if err := config.SaveToken(tok); err != nil {
			return fmt.Errorf("save token: %w", err)
		}
		g.createService()
		return nil
	case err := <-errChan:
		return err
	case <-time.After(2 * time.Minute):
		return fmt.Errorf("authentication timed out after 2 minutes")
	}
}

func (g *Provider) FreeSlots(day time.Time, minDuration time.Duration) ([]provider.Slot, error) {
	if g.service == nil {
		return nil, fmt.Errorf("not authenticated; run 'time-broker auth'")
	}
	return nil, fmt.Errorf("free: not yet implemented")
}

func (g *Provider) Book(title string, start, end time.Time) error {
	if g.service == nil {
		return fmt.Errorf("not authenticated; run 'time-broker auth'")
	}
	return fmt.Errorf("book: not yet implemented")
}

func (g *Provider) createService() {
	ctx := context.Background()
	ts := g.config.TokenSource(ctx, g.token)
	saving := &saveTokenSource{
		src: oauth2.ReuseTokenSource(g.token, ts),
	}
	client := oauth2.NewClient(ctx, saving)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return
	}
	g.service = srv
}

type saveTokenSource struct {
	src oauth2.TokenSource
}

func (s *saveTokenSource) Token() (*oauth2.Token, error) {
	tok, err := s.src.Token()
	if err != nil {
		return nil, err
	}
	if err := config.SaveToken(tok); err != nil {
		return nil, fmt.Errorf("save refreshed token: %w", err)
	}
	return tok, nil
}

func randomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func openURL(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
