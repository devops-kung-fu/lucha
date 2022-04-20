package util

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintTabbed(t *testing.T) {
	output := CaptureOutput(func() {
		PrintTabbed("[TEST]")
	})

	assert.Contains(t, output, "\t", "Console output does not contain a tab character")
	assert.GreaterOrEqual(t, len(output), 0, "No information logged to STDOUT")
	assert.Equal(t, strings.Count(output, "\n"), 1, "Expected one line of log output")
}

func TestPrintSuccess(t *testing.T) {
	output := CaptureOutput(func() {
		PrintSuccess("[TEST]")
	})

	assert.GreaterOrEqual(t, len(output), 0, "No information logged to STDOUT")
	assert.Equal(t, strings.Count(output, "\n"), 1, "Expected one line of log output")
}

func TestPrintWarning(t *testing.T) {
	output := CaptureOutput(func() {
		PrintWarning("[TEST]")
	})

	assert.GreaterOrEqual(t, len(output), 0, "No information logged to STDOUT")
	assert.Equal(t, strings.Count(output, "\n"), 1, "Expected one line of log output")
}

func TestPrintInfo(t *testing.T) {
	output := CaptureOutput(func() {
		PrintInfo("[TEST]")
	})

	assert.GreaterOrEqual(t, len(output), 0, "No information logged to STDOUT")
	assert.Equal(t, strings.Count(output, "\n"), 1, "Expected one line of log output")
}

func TestPrintErr(t *testing.T) {
	output := CaptureOutput(func() {
		PrintErr("[TEST]", errors.New("Test Error"))
	})

	assert.GreaterOrEqual(t, len(output), 0, "No information logged to STDOUT")
	assert.Equal(t, strings.Count(output, "\n"), 1, "Expected two lines of log output")
}
