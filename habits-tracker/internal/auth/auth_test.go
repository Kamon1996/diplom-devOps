package auth

import (
	"context"
	"testing"

	"habit-tracker/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	authService := NewService(db)
	ctx := context.Background()

	t.Run("Hash and Check Password", func(t *testing.T) {
		password := "testpassword123"

		hash, err := authService.HashPassword(password)
		require.NoError(t, err)
		require.NotEmpty(t, hash)

		err = authService.CheckPassword(hash, password)
		assert.NoError(t, err)

		err = authService.CheckPassword(hash, "wrongpassword")
		assert.Error(t, err)
	})

	t.Run("Register User", func(t *testing.T) {
		err := authService.Register(ctx, "testuser", "password123")
		assert.NoError(t, err)

		err = authService.Register(ctx, "testuser", "password123")
		assert.Error(t, err)
	})

	t.Run("Authenticate User", func(t *testing.T) {
		username := "authuser"
		password := "authpassword"

		err := authService.Register(ctx, username, password)
		require.NoError(t, err)

		userID, err := authService.Authenticate(ctx, username, password)
		assert.NoError(t, err)
		assert.NotEmpty(t, userID)

		_, err = authService.Authenticate(ctx, username, "wrongpassword")
		assert.Error(t, err)

		_, err = authService.Authenticate(ctx, "nonexistent", password)
		assert.Error(t, err)
	})
}
