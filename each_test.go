package xlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIgnoredPath(t *testing.T) {
	assert.True(t, IsIgnoredPath(".git/config"))
	assert.True(t, IsIgnoredPath(".versions/config"))
	assert.False(t, IsIgnoredPath("index.md"))
	assert.False(t, IsIgnoredPath("something/something"))
}

func TestIsNil(t *testing.T) {
	assert.True(t, isNil[Page](nil))
	assert.True(t, isNil[*Page](nil))
}
