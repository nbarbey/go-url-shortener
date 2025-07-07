package urlshortener

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStoreWithExpiration(t *testing.T) {
	s := NewInMemorySqlite()
	now := time.Now()
	err := s.Save("http://long.net", "http://short.uk", &now)
	require.NoError(t, err)

	u, err := s.Get("http://short.uk")
	require.NoError(t, err)

	actual := *u.expiration
	assert.Equal(t, now.Format(time.RFC3339), actual.Format(time.RFC3339))
}
