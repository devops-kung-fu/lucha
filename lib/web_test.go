package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DownloadURL(t *testing.T) {
	_, err := DownloadURL("x", "/tmp")
	assert.Error(t, err, "URL should be a valid URI")
}
