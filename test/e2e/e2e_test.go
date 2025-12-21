//go:build e2e

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL = "http://localhost:8000"
	timeout = 30 * time.Second
)

func waitForAPI(t *testing.T) {
	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := client.Get(baseURL + "/health")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				t.Log("API is ready")
				return
			}
		}
		time.Sleep(1 * time.Second)
	}

	t.Fatal("API did not become ready in time")
}

func TestE2E_GetNews_Success(t *testing.T) {
	waitForAPI(t)

	resp, err := http.Get(baseURL + "/news?id=1")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// ReadAll returns bytes
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	bodyStr := string(body)
	assert.Contains(t, bodyStr, "Lorem ipsum")
}

func TestE2E_GetNews_NotFound(t *testing.T) {
	waitForAPI(t)

	resp, err := http.Get(baseURL + "/news?id=666")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "article not found", result["error"])
}

func TestE2E_GetNews_MissingID(t *testing.T) {
	waitForAPI(t)

	resp, err := http.Get(baseURL + "/news")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "id parameter is required", result["error"])
}
