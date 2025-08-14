package auth

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type AuthConfig struct {
	Token  string `json:"token"`
	APIUrl string `json:"api_url"`
}

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".snapshell-config.json"
	}
	return filepath.Join(homeDir, ".snapshell-config.json")
}

func LoadConfig() (*AuthConfig, error) {
	configPath := GetConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config AuthConfig
	err = json.Unmarshal(data, &config)
	return &config, err
}

func SaveConfig(config *AuthConfig) error {
	configPath := GetConfigPath()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0600)
}

func PerformTokenLogin(apiURL string) error {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("could not start local server: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	callbackURL := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	tokenURL := fmt.Sprintf("%s/api/token?callback=%s", apiURL, callbackURL)
	fmt.Println("If the browser doesn't open automatically, please visit the URL above.")

	if err := openBrowser(tokenURL); err != nil {
		fmt.Printf("Could not open browser automatically: %v\n", err)
	}

	tokenCh := make(chan string)
	server := &http.Server{}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" && r.Method == "POST" {
			var body struct {
				Token string `json:"token"`
			}
			err := json.NewDecoder(r.Body).Decode(&body)
			if err == nil {
				token = body.Token
			}
		}
		if token == "" {
			http.Error(w, "Missing token", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Login successful! You may close this window.")
		tokenCh <- token
	})
	go func() {
		server.Serve(listener)
	}()
	token := <-tokenCh
	server.Close()

	fmt.Printf("Token received: %s\n", token)

	// Save token to config
	config := &AuthConfig{
		Token:  token,
		APIUrl: apiURL,
	}
	if err := SaveConfig(config); err != nil {
		return fmt.Errorf("error saving config: %v", err)
	}
	fmt.Printf("✅ CLI login successful! Token saved to %s\n", GetConfigPath())
	return nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	return exec.Command(cmd, args...).Start()
}
