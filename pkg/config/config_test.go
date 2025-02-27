package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnv(t *testing.T) {
	key := "asdfhjasfetklq2453"
	value := "1234"
	def := "defaultValue"
	err := os.Setenv(key, value)
	require.NoError(t, err, "setting env %s", key)

	// assigned key
	got := getEnv(key, def)
	assert.Equal(t, value, got, "env variable should match %s", value)

	// non-existent key
	got = getEnv("adsddsfaadsfadsfadsfadsf", def)
	assert.Equal(t, def, got, "env variable should match %s", def)

	err = os.Unsetenv(key)
	assert.NoError(t, err, "setting env %s", key)
}

func TestNewTestingConfig(t *testing.T) {
	cfg := NewTestingConfig()
	assert.Equal(t, cfg.ServerPort, ":0")
}

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()

	envFile := filepath.Join(tmpDir, ".env")

	envContent := `SERVER_PORT=9090
SECRET_KEY=mysecret
DISK_UPLOADS_LOCATION=/tmp/uploads`

	if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to write temp env file: %v", err)
	}

	var cfg Config

	cfg.LoadConfig(envFile)

	assert.Equal(t, "9090", cfg.ServerPort, "Expected SERVER_PORT to be 9090")
	assert.Equal(t, "mysecret", cfg.JWTsecretKey, "Expected SECRET_KEY to be mysecret")
	assert.Equal(t, "/tmp/uploads", cfg.UploadsLocation, "Expected DISK_UPLOADS_LOCATION to be /tmp/uploads")
}

func TestLoadConfigDefaults(t *testing.T) {
	os.Clearenv()

	var cfg Config

	cfg.LoadConfig("")

	assert.Equal(t, ":8080", cfg.ServerPort, "Expected default SERVER_PORT to be :8080")
	assert.Equal(t, "als;dgasdfkasbf2ql4q", cfg.JWTsecretKey, "Expected default SECRET_KEY")
	assert.Equal(t, "uploads", cfg.UploadsLocation, "Expected default DISK_UPLOADS_LOCATION to be uploads")
}
