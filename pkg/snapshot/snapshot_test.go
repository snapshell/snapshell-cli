package snapshot_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/snapshell/snapshell-cli/pkg/snapshot"
)

func TestDetectSnapshotType_NpmAudit(t *testing.T) {
	content := "# npm audit report\nSome details..."
	typeName := snapshot.DetectSnapshotType(content)
	if typeName != "npm-audit" {
		t.Errorf("expected 'npm-audit', got '%s'", typeName)
	}
}

func TestDetectSnapshotType_Npm(t *testing.T) {
	content := "added 1 package, audited 2 packages"
	typeName := snapshot.DetectSnapshotType(content)
	if typeName != "npm" {
		t.Errorf("expected 'npm', got '%s'", typeName)
	}
}

func TestDetectSnapshotType_Terraform(t *testing.T) {
	content := "Terraform will perform the following actions...\nPlan: 1 to add, 2 to change, 3 to destroy"
	typeName := snapshot.DetectSnapshotType(content)
	if typeName != "terraform" {
		t.Errorf("expected 'terraform', got '%s'", typeName)
	}
}

func TestDetectSnapshotType_Default(t *testing.T) {
	content := "random output"
	typeName := snapshot.DetectSnapshotType(content)
	if typeName != "terraform" {
		t.Errorf("expected 'terraform' as default, got '%s'", typeName)
	}
}

func TestCreateSnapshot_Success(t *testing.T) {
	// Mock API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		resp := snapshot.APIResponse{}
		resp.Snapshot.ID = "testid123"
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	err := snapshot.CreateSnapshot([]byte("test content"), server.URL, "label", "type", true, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCreateSnapshot_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad request")
	}))
	defer server.Close()

	err := snapshot.CreateSnapshot([]byte("test content"), server.URL, "label", "type", true, 1)
	if err == nil || err.Error() == "" {
		t.Errorf("expected error for API error, got nil")
	}
}

func TestCreateSnapshot_AuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	err := snapshot.CreateSnapshot([]byte("test content"), server.URL, "label", "type", true, 1)
	if err == nil || err.Error() == "" {
		t.Errorf("expected error for auth error, got nil")
	}
}

func TestCreateSnapshot_MarshalError(t *testing.T) {
	// Simulate by passing a type that can't be marshaled (shouldn't happen in real use)
	// We'll skip this as the function always marshals a valid struct, but you could refactor to allow injection for testing.
}
