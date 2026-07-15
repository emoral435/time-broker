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
	"sort"
	"strings"
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

var redirectPort = ":8085"
var authTimeout = 2 * time.Minute
var testListener net.Listener

const (
	ProviderName = "google"
)

type Provider struct {
	config   *oauth2.Config
	token    *oauth2.Token
	service  *calendar.Service
	endpoint string
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
		_ = p.createService()
	}

	return p
}

func (g *Provider) Name() string {
	return ProviderName
}

func (g *Provider) EnsureAuthenticated() error {
	if g.token == nil {
		return fmt.Errorf("authentication required; run 'time-broker auth' to authenticate with your calendar provider")
	}

	if g.service == nil {
		if err := g.createService(); err != nil {
			return fmt.Errorf("failed to initialize calendar service: %w", err)
		}
	}

	if err := g.validateToken(); err != nil {
		if g.token.RefreshToken != "" {
			if refreshErr := g.refreshToken(); refreshErr != nil {
				return fmt.Errorf("token validation failed and refresh failed (%v); run 'time-broker auth' to re-authenticate", err)
			}
			if err := g.createService(); err != nil {
				return fmt.Errorf("failed to re-initialize calendar service: %w", err)
			}
			if err := g.validateToken(); err != nil {
				return fmt.Errorf("token validation failed after refresh: %w; run 'time-broker auth' to re-authenticate", err)
			}
			return nil
		}
		return fmt.Errorf("token validation failed: %w; run 'time-broker auth' to re-authenticate", err)
	}

	return nil
}

func (g *Provider) refreshToken() error {
	ctx := context.Background()
	ts := g.config.TokenSource(ctx, g.token)
	tok, err := ts.Token()
	if err != nil {
		return fmt.Errorf("refresh token: %w", err)
	}
	g.token = tok
	return config.SaveToken(tok)
}

func (g *Provider) Auth() error {
	if g.token != nil && g.service != nil {
		if err := g.validateToken(); err == nil {
			return nil
		}
	}

	if g.config.ClientID == "" || g.config.ClientSecret == "" {
		return fmt.Errorf("google OAuth credentials not configured: set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables, or inject them at build time")
	}

	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no code in callback response")
			return
		}
		fmt.Fprint(w, `<html><body><h1>Authentication to time-broker successful!</h1><p>You can close this tab.</p></body></html>`)
		codeChan <- code
	})

	var listener net.Listener
	if testListener != nil {
		listener = testListener
		testListener = nil
	} else {
		var err error
		listener, err = net.Listen("tcp", redirectPort)
		if err != nil {
			return fmt.Errorf("start callback server on port %s: %w", redirectPort, err)
		}
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
			return fmt.Errorf("save token error: %w", err)
		}
		if err := g.createService(); err != nil {
			return fmt.Errorf("authentication succeeded but failed to initialize calendar service: %w", err)
		}
		if err := g.validateToken(); err != nil {
			return fmt.Errorf("authentication succeeded but token validation failed: %w; run 'time-broker auth' again", err)
		}
		return nil
	case err := <-errChan:
		return err
	case <-time.After(authTimeout):
		return fmt.Errorf("authentication timed out after %v", authTimeout)
	}
}

func (g *Provider) FreeSlots(start, end time.Time, minDuration time.Duration) ([]provider.Slot, error) {
	if err := g.EnsureAuthenticated(); err != nil {
		return nil, err
	}

	events, err := g.service.Events.List("primary").
		TimeMin(start.Format(time.RFC3339)).
		TimeMax(end.Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	type busySlot struct {
		start time.Time
		end   time.Time
	}

	var busy []busySlot
	for _, ev := range events.Items {
		s, _ := parseEventDateTime(ev.Start)
		e, _ := parseEventDateTime(ev.End)
		if s.IsZero() || e.IsZero() {
			continue
		}
		if s.Before(start) {
			s = start
		}
		if e.After(end) {
			e = end
		}
		busy = append(busy, busySlot{start: s, end: e})
	}

	sort.Slice(busy, func(i, j int) bool {
		return busy[i].start.Before(busy[j].start)
	})

	var free []provider.Slot
	cursor := start
	for _, b := range busy {
		if b.start.After(cursor) && b.start.Sub(cursor) >= minDuration {
			free = append(free, provider.Slot{Start: cursor, End: b.start})
		}
		if b.end.After(cursor) {
			cursor = b.end
		}
	}
	if cursor.Before(end) && end.Sub(cursor) >= minDuration {
		free = append(free, provider.Slot{Start: cursor, End: end})
	}

	return free, nil
}

func (g *Provider) Book(title, description string, start, end time.Time, allDay bool) error {
	if err := g.EnsureAuthenticated(); err != nil {
		return err
	}

	event := &calendar.Event{
		Summary:     title,
		Description: description,
	}

	if allDay {
		event.Start = &calendar.EventDateTime{
			Date: start.Format("2006-01-02"),
		}
		event.End = &calendar.EventDateTime{
			Date: end.Format("2006-01-02"),
		}
	} else {
		event.Start = &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: resolveTimezone(start),
		}
		event.End = &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: resolveTimezone(end),
		}
	}

	if g.service.Events == nil {
		return fmt.Errorf("calendar service not fully initialized")
	}

	_, err := g.service.Events.Insert("primary", event).Do()
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (g *Provider) EventsForDay(day time.Time) ([]provider.Event, error) {
	if g.service == nil {
		return nil, fmt.Errorf("not authenticated; run 'time-broker auth'")
	}

	day = day.In(time.Local)
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)

	events, err := g.service.Events.List("primary").
		TimeMin(start.Format(time.RFC3339)).
		TimeMax(end.Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	var result []provider.Event
	for _, ev := range events.Items {
		s, allDay := parseEventDateTime(ev.Start)
		e, _ := parseEventDateTime(ev.End)

		result = append(result, provider.Event{
			Title:       ev.Summary,
			Description: ev.Description,
			Location:    ev.Location,
			Start:       s,
			End:         e,
			AllDay:      allDay,
		})
	}

	return result, nil
}

func (g *Provider) Timezone() *time.Location {
	if g.service == nil {
		return time.Local
	}

	cal, err := g.service.CalendarList.Get("primary").Do()
	if err != nil {
		return time.Local
	}

	loc, err := time.LoadLocation(cal.TimeZone)
	if err != nil {
		return time.Local
	}

	return loc
}

func parseEventDateTime(edt *calendar.EventDateTime) (time.Time, bool) {
	if edt == nil {
		return time.Time{}, false
	}

	if edt.Date != "" {
		t, err := time.Parse("2006-01-02", edt.Date)
		if err != nil {
			return time.Time{}, true
		}
		return t, true
	}

	if edt.DateTime != "" {
		t, err := time.Parse(time.RFC3339, edt.DateTime)
		if err != nil {
			return time.Time{}, false
		}
		return t, false
	}

	return time.Time{}, false
}

func (g *Provider) createService(opts ...option.ClientOption) error {
	ctx := context.Background()
	ts := g.config.TokenSource(ctx, g.token)
	saving := &saveTokenSource{
		src: oauth2.ReuseTokenSource(g.token, ts),
	}
	client := oauth2.NewClient(ctx, saving)
	allOpts := append([]option.ClientOption{option.WithHTTPClient(client)}, opts...)
	if g.endpoint != "" {
		allOpts = append(allOpts, option.WithEndpoint(g.endpoint))
	}
	srv, err := calendar.NewService(ctx, allOpts...)
	if err != nil {
		return fmt.Errorf("create calendar service: %w", err)
	}
	g.service = srv
	return nil
}

func (g *Provider) validateToken() error {
	if g.service == nil {
		return fmt.Errorf("calendar service not initialized")
	}
	if g.service.CalendarList == nil {
		return fmt.Errorf("calendar service not fully configured")
	}
	_, err := g.service.CalendarList.List().MaxResults(1).Do()
	if err != nil {
		return fmt.Errorf("calendar API: %w", err)
	}
	return nil
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

var openURL = func(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		return fmt.Errorf("unsupported OS platform: %s", runtime.GOOS)
	}
}

func resolveTimezone(t time.Time) string {
	loc := t.Location()
	if loc != time.Local {
		return loc.String()
	}

	if tz := os.Getenv("TZ"); tz != "" {
		if _, err := time.LoadLocation(tz); err == nil {
			return tz
		}
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if name, ok := timezoneFromEtcLocaltime(); ok {
			return name
		}
	}

	_, offset := t.Zone()
	sign := "+"
	if offset < 0 {
		sign = "-"
		offset = -offset
	}
	hours := offset / 3600
	minutes := (offset % 3600) / 60
	return fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)
}

var timezoneFromEtcLocaltime = func() (string, bool) {
	link, err := os.Readlink("/etc/localtime")
	if err != nil {
		return "", false
	}
	const prefix = "zoneinfo/"
	idx := strings.Index(link, prefix)
	if idx < 0 {
		return "", false
	}
	return link[idx+len(prefix):], true
}
