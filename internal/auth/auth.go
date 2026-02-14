package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	driveapi "google.golang.org/api/drive/v3"
)

// GetClient returns an authenticated HTTP client for Google Drive.
// It uses OAuth2 with a local redirect to capture the authorization code.
// Tokens are cached in ~/.config/storage/token.json.
func GetClient() (*http.Client, error) {
	credPath, err := credentialsPath()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(credPath)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read credentials file at %s: %w\n\n"+
				"To set up credentials:\n"+
				"1. Go to https://console.cloud.google.com/apis/credentials\n"+
				"2. Create an OAuth 2.0 Client ID (Desktop app)\n"+
				"3. Download the JSON and save it to %s",
			credPath, err, credPath,
		)
	}

	config, err := google.ConfigFromJSON(b, driveapi.DriveReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %w", err)
	}

	tok, err := loadToken()
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		if err := saveToken(tok); err != nil {
			return nil, err
		}
	}

	return config.Client(context.Background(), tok), nil
}

func credentialsPath() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "credentials.json"), nil
}

func tokenPath() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "token.json"), nil
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to find home directory: %w", err)
	}
	dir := filepath.Join(home, ".config", "storage")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("unable to create config directory: %w", err)
	}
	return dir, nil
}

func loadToken() (*oauth2.Token, error) {
	path, err := tokenPath()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	return tok, json.NewDecoder(f).Decode(tok)
}

func saveToken(tok *oauth2.Token) error {
	path, err := tokenPath()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("unable to save token: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(tok)
}

// getTokenFromWeb starts a local HTTP server on a random port to receive
// the OAuth2 callback, opens the browser for consent, and returns the token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	// Start local server to capture the redirect
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no code in callback")
			fmt.Fprintln(w, "Error: no authorization code received.")
			return
		}
		codeCh <- code
		fmt.Fprintln(w, "Authorization successful! You can close this tab.")
	})

	server := &http.Server{
		Addr:    "127.0.0.1:9874",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	config.RedirectURL = "http://127.0.0.1:9874/callback"

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Open this URL in your browser to authorize:\n\n%s\n\n", authURL)

	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		return nil, fmt.Errorf("authorization failed: %w", err)
	}

	_ = server.Close()

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("unable to exchange code for token: %w", err)
	}

	return tok, nil
}
