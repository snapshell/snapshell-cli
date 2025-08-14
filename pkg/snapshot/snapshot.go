package snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/snapshell/snapshell-cli/pkg/auth"
)

type SnapshotRequest struct {
	Label         string `json:"label"`
	Type          string `json:"type"`
	Content       string `json:"content"`
	IsPrivate     bool   `json:"isPrivate"`
	ExpiresInDays int    `json:"expiresInDays"`
}

type APIResponse struct {
	Snapshot struct {
		ID string `json:"id"`
	} `json:"snapshot"`
}

func DetectSnapshotType(content string) string {

	// Trivy JSON detection (must be first)
	trimmed := strings.TrimSpace(content)
	if strings.HasPrefix(trimmed, "{") {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(trimmed), &m); err == nil {
			keys := []string{"SchemaVersion", "ArtifactType", "Metadata", "Results"}
			found := true
			for _, k := range keys {
				if _, ok := m[k]; !ok {
					found = false
					break
				}
			}
			if found {
				return "trivy"
			}
		}
	}
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "# npm audit report" {
			return "npm-audit"
		}
	}
	if strings.Contains(content, "added") || strings.Contains(content, "changed") ||
		strings.Contains(content, "audited") && strings.Contains(content, "packages") {
		for _, line := range lines {
			if strings.Contains(line, "added") && strings.Contains(line, "packages") {
				return "npm"
			}
		}
	}
	if strings.Contains(content, "Terraform will perform") ||
		strings.Contains(content, "Plan:") && (strings.Contains(content, "to add") || strings.Contains(content, "to change") || strings.Contains(content, "to destroy")) {
		return "terraform"
	}
	return "terraform"
}

func CreateSnapshot(content []byte, apiURL, label, typeFlag string, private bool, expires int) error {
	config, _ := auth.LoadConfig()
	if config != nil && config.APIUrl != "" {
		apiURL = config.APIUrl
	}
	snapshotType := typeFlag
	if snapshotType == "" {
		snapshotType = DetectSnapshotType(string(content))
		fmt.Fprintf(os.Stderr, "Auto-detected snapshot type: %s\n", snapshotType)
	}
	req := SnapshotRequest{
		Label:         label,
		Type:          snapshotType,
		Content:       string(content),
		IsPrivate:     private,
		ExpiresInDays: expires,
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}
	apiEndpoint := apiURL + "/api/snapshots"
	httpReq, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if config != nil {
		if config.Token != "" {
			httpReq.Header.Set("Authorization", "Bearer "+config.Token)
		}
	} else {
		expires = 1 // Default to 1 day if not authenticated
	}
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error sending request to %s: %v\nMake sure the web server is running on %s", apiEndpoint, err, apiURL)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized && config != nil {
		return fmt.Errorf("authentication expired or invalid. Please re-authenticate:\n  snapshell login")
	}
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("error parsing response: %v", err)
	}
	fullURL := apiURL + "/snapshots/" + apiResp.Snapshot.ID
	fmt.Println(fullURL)
	return nil
}
