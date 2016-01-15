package utils

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenKey(t *testing.T) {
	filter := make(map[string]bool)
	for i := 0; i < 10000; i++ {
		key := GenerateKey()
		assert.True(t, !filter[key], "should not gen old key")
		filter[key] = true
	}
}

func TestGetNowSecond(t *testing.T) {
	now := GetNowSecond()
	time.Sleep(1 * time.Second)
	assert.True(t, GetNowSecond()-now == 1)
}
