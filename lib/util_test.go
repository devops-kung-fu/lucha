package lib

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintIf(t *testing.T) {
	result, err := captureStdout(func() { fmt.Println("Test") })

	assert.Equal(t, "Test\n", result, "Should match the string Test")
	assert.NoError(t, err, "No error should have been generated")
}

func captureStdout(f func()) (captured string, err error) {
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}
	origStdout := os.Stdout
	os.Stdout = w

	f()

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout = origStdout
	captured = string(buf[:n])
	return
}

func TestLuchaDir(t *testing.T) {
	path, err := LuchaDir()
	assert.NoError(t, err, "There should be no error")
	assert.NotEmpty(t, path, "Path should not be empty")
}

func TestStartsWith(t *testing.T) {
	testArray := []string{
		"foo",
		"test",
	}
	startsWith := StartsWith(testArray, "foobar")
	assert.True(t, startsWith, "The test array should contain an element that starts with `te`")

	startsWith = StartsWith(testArray, "ge")
	assert.False(t, startsWith, "The test array should not contain an element that starts with `ge`")

}

func TestContains(t *testing.T) {
	testArray := []string{
		"testing",
		"foo",
	}
	contains := Contains(testArray, "foo")
	assert.True(t, contains, "The test array should contain an element `foo`")

	contains = Contains(testArray, "ge")
	assert.False(t, contains, "The test array should not contain an element  `bar`")

}
