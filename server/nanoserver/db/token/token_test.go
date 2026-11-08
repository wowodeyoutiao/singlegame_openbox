package token

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUserToken_GenerateToken(t *testing.T) {
	// Setup
	viper.Set("token.jwtTokenSignKey", "test-secret-key")
	userTokenFactory := CreateUserFactory()

	tests := []struct {
		name        string
		userid      int64
		accountName string
		platform    string
		serverId    int32
		expireAt    int64
		wantErr     bool
	}{
		{
			name:        "Valid token generation",
			userid:      12345,
			accountName: "testuser",
			platform:    "web",
			serverId:    1,
			expireAt:    3600, // 1 hour
			wantErr:     false,
		},
		{
			name:        "Token with zero expiration",
			userid:      67890,
			accountName: "user2",
			platform:    "mobile",
			serverId:    2,
			expireAt:    0,
			wantErr:     false,
		},
		{
			name:        "Token with negative expiration",
			userid:      11111,
			accountName: "user3",
			platform:    "desktop",
			serverId:    3,
			expireAt:    -100,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := userTokenFactory.GenerateToken(tt.userid, tt.accountName, tt.platform, tt.serverId, tt.expireAt)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Parse the token to verify its contents
				claims, parseErr := userTokenFactory.ParseToken(token)
				assert.NoError(t, parseErr)
				assert.Equal(t, tt.userid, claims.UserId)
				assert.Equal(t, tt.accountName, claims.AccountName)
				assert.Equal(t, tt.platform, claims.Platform)
				assert.Equal(t, tt.serverId, claims.ServerId)

				// Verify expiration time is set correctly
				expectedExpiry := time.Now().Unix() + tt.expireAt
				assert.InDelta(t, expectedExpiry, claims.ExpiresAt, 5) // Allow 5 seconds tolerance
			}
		})
	}
}

func TestUserToken_ConcurrentGeneration(t *testing.T) {
	// Setup
	viper.Set("token.jwtTokenSignKey", "test-secret-key")
	userTokenFactory := CreateUserFactory()

	// Channel to collect results
	type result struct {
		userID int64
		token  string
		err    error
	}
	results := make(chan result, 100)

	// Launch 100 goroutines
	for i := 1; i <= 100; i++ {
		go func(userID int64) {
			token, err := userTokenFactory.GenerateToken(
				userID,
				"testuser"+string(rune(userID)), // accountName
				"web",                           // platform
				1,                               // serverId
				3600,                            // expireAt (1 hour)
			)
			results <- result{userID: userID, token: token, err: err}
		}(int64(i))
	}

	// Collect all results
	tokens := make(map[int64]string)
	userIDs := make([]int64, 0, 100)
	for i := 0; i < 100; i++ {
		res := <-results
		assert.NoError(t, res.err)
		assert.NotEmpty(t, res.token)
		tokens[res.userID] = res.token
		userIDs = append(userIDs, res.userID)
	}

	// Verify all tokens are unique and valid
	assert.Equal(t, 100, len(tokens))

	// Verify all UserIDs from 1-100 are present
	expectedUserIDs := make(map[int64]bool)
	for i := int64(1); i <= 100; i++ {
		expectedUserIDs[i] = true
	}

	// Check that all expected UserIDs are present and no duplicates
	for userID := range tokens {
		assert.True(t, expectedUserIDs[userID], "Unexpected UserID: %d", userID)
		delete(expectedUserIDs, userID)
	}
	assert.Empty(t, expectedUserIDs, "Missing UserIDs: %v", expectedUserIDs)

	// Verify each token contains correct UserID
	for userID, token := range tokens {
		claims, err := userTokenFactory.ParseToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserId)
	}
}
