package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestSaveTokenAndLoadToken(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	tok := &oauth2.Token{
		AccessToken:  "test-access",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh",
		Expiry:       time.Now().Add(time.Hour),
	}
	if err := SaveToken(tok); err != nil {
		t.Fatalf("SaveToken() error: %v", err)
	}

	got, err := LoadToken()
	if err != nil {
		t.Fatalf("LoadToken() error: %v", err)
	}
	if got == nil {
		t.Fatal("LoadToken() returned nil")
	}
	if got.AccessToken != tok.AccessToken {
		t.Errorf("AccessToken = %q; want %q", got.AccessToken, tok.AccessToken)
	}
	if got.RefreshToken != tok.RefreshToken {
		t.Errorf("RefreshToken = %q; want %q", got.RefreshToken, tok.RefreshToken)
	}
	if got.TokenType != tok.TokenType {
		t.Errorf("TokenType = %q; want %q", got.TokenType, tok.TokenType)
	}
}

func TestLoadTokenFileNotFound(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	tok, err := LoadToken()
	if err != nil {
		t.Fatalf("LoadToken() error: %v", err)
	}
	if tok != nil {
		t.Fatal("LoadToken() should return nil when file does not exist")
	}
}

func TestTokenPath(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	path, err := TokenPath()
	if err != nil {
		t.Fatalf("TokenPath() error: %v", err)
	}
	want := filepath.Join(dir, ".time-broker", "tokens.json")
	if path != want {
		t.Errorf("TokenPath() = %q; want %q", path, want)
	}
}

func TestSaveTokenCreatesFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	tok := &oauth2.Token{AccessToken: "test"}
	if err := SaveToken(tok); err != nil {
		t.Fatalf("SaveToken() error: %v", err)
	}

	path, _ := TokenPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("token file %s was not created", path)
	}
}
