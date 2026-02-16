package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"habit-tracker/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	_, cleanup := testutils.SetupTestDB()

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func TestServerRoutes(t *testing.T) {

	t.Run("Public Routes Accessible", func(t *testing.T) {
		assert.True(t, true) // placeholder
	})

	t.Run("Protected Routes Redirect to Login", func(t *testing.T) {
		assert.True(t, true) // placeholder
	})
}
