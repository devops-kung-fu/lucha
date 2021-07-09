package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLuchaDir(t *testing.T) {
	path, err := LuchaDir()
	assert.NoError(t, err)
	assert.Contains(t, path, ".lucha")
}
