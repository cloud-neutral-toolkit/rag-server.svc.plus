package integration

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// Try loading .env, but don't fail if missing (CI envs)
	_ = godotenv.Load("../../.env")
	os.Exit(m.Run())
}

func TestServerHealth(t *testing.T) {
	baseURL := os.Getenv("SERVER_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	client := &http.Client{Timeout: 5 * time.Second}

	// Retry loop to wait for server start in real integration scenarios
	var resp *http.Response
	var err error
	for i := 0; i < 5; i++ {
		resp, err = client.Get(baseURL + "/health")
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		t.Skip("Skipping integration test: server not reachable")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
