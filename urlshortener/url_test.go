package urlshortener

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNotExpiring(t *testing.T) {
	assert.False(t, MustNewURL("http://blabla.net", nil).Expiring())
}

func TestExpiring(t *testing.T) {
	now := time.Now()
	assert.True(t, MustNewURL("http://blabla.net", &now).Expiring())
}

func TestNotExpired(t *testing.T) {
	exp := time.Now().Add(-1 * time.Hour)
	assert.True(t, MustNewURL("http://blabla.net", &exp).ExpiredAt(time.Now()))
}

func TestExpired(t *testing.T) {
	exp := time.Now().Add(1 * time.Hour)
	assert.False(t, MustNewURL("http://blabla.net", &exp).ExpiredAt(time.Now()))
}
