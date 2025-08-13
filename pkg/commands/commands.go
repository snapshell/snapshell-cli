package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/snapshell/snapshell-cli/pkg/auth"
	"github.com/snapshell/snapshell-cli/pkg/snapshot"

	"github.com/spf13/cobra"
)

var (
	apiURL   string
	label    string
	typeFlag string
	private  bool
	expires  int
	file     string
)

var rootCmd = &cobra.Command{
	Use:   "snapshell",
	Short: "Convert CLI output to shareable snapshots",
	Long: `SnapShell CLI converts raw CLI output (like terraform plan, npm audit) into clean, styled, and shareable web snapshots.

Examples:
  terraform plan | snapshell --label="My Plan"
  npm audit | snapshell --label="Security Audit"
  snapshell --file=plan.txt --label="My Plan"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var content []byte
		var err error
		if file != "" {
			content, err = os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", file, err)
			}
			// If label is empty, use the file's base name
			if label == "" {
				label = filepath.Base(file)
			}
		} else {
			content, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("error reading stdin: %v", err)
			}
			// If label is empty, use inferred type plus datetime
			if label == "" {
				snapshotType := typeFlag
				if snapshotType == "" {
					snapshotType = snapshot.DetectSnapshotType(string(content))
				}
				label = fmt.Sprintf("%s-%s", snapshotType, time.Now().Format("2006-01-02_15-04-05"))
			}
		}
		if len(content) == 0 {
			return fmt.Errorf("no input provided. Use --file or pipe content to stdin")
		}
		return snapshot.CreateSnapshot(content, apiURL, label, typeFlag, private, expires)
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login via browser to get CLI token",
	Long:  `Open browser to authenticate and get a CLI token. This is the recommended authentication method.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return auth.PerformTokenLogin(apiURL)
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear authentication and logout",
	Long:  `Remove stored authentication tokens and logout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := auth.GetConfigPath()
		if err := os.Remove(configPath); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("error removing config: %v", err)
			}
		}
		fmt.Println("✅ Logged out successfully")
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiURL, "api", "https://snapshell.dev", "API base URL")
	rootCmd.Flags().StringVar(&label, "label", "", "Snapshot label")
	rootCmd.Flags().StringVar(&typeFlag, "type", "", "Snapshot type (auto-detected if not specified)")
	rootCmd.Flags().BoolVar(&private, "private", true, "Make snapshot private")
	rootCmd.Flags().IntVar(&expires, "expires", 30, "Snapshot expiration in days")
	rootCmd.Flags().StringVar(&file, "file", "", "Read from file instead of stdin")
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
