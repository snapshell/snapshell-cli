package main

import (
	"fmt"
	"os"

	"github.com/snapshell/snapshell-cli/pkg/commands"
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

// Global flags
var (
	apiURL   string
	label    string
	typeFlag string
	private  bool
	expires  int
	file     string
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
