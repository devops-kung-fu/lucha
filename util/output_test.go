package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureOutput(t *testing.T) {
	output := CaptureOutput(func() {
		fmt.Print("TEST")
	})
	assert.NotNil(t, output)
	assert.Len(t, output, 4)
}
