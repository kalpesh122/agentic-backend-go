package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAppliesDefaultsAndParsesLists(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://a:b@localhost:5432/x")
	t.Setenv("API_KEYS", "k1:alice, k2")
	t.Setenv("CORS_ORIGINS", "http://a,http://b")
	c, err := Load()
	require.NoError(t, err)
	assert.Equal(t, 8080, c.Port)
	assert.Equal(t, map[string]string{"k1": "alice", "k2": "k2"}, c.APIKeyMap())
	assert.Equal(t, []string{"http://a", "http://b"}, c.CORSOrigins)
}

func TestLoadReportsEveryProblem(t *testing.T) {
	t.Setenv("DATABASE_URL", "mysql://x")
	t.Setenv("ENV", "staging")
	t.Setenv("PORT", "70000")
	_, err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DATABASE_URL")
	assert.Contains(t, err.Error(), "ENV")
	assert.Contains(t, err.Error(), "PORT")
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	_, err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DATABASE_URL")
}
