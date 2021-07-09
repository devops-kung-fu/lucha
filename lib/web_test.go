//Package lib Functionality for the Hookz CLI
package lib

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_DownloadFile(t *testing.T) {
	fs := FileSystem{
		fs: afero.NewMemMapFs(),
	}
	_, err := DownloadFile(fs, "x", "x")
	assert.Error(t, err, "URL should be a valid URI")
}

func TestWriteCounter_Write(t *testing.T) {
	wc := WriteCounter{}
	count, err := wc.Write([]byte("test"))
	assert.NoError(t, err, "There should be no error")
	assert.Equal(t, 4, count, "4 bytes should have been written")
}

func TestDownloadFile(t *testing.T) {
	newFs := FileSystem{
		fs: afero.NewMemMapFs(),
	}
	URL := "https://raw.githubusercontent.com/devops-kung-fu/lucha/main/lucha.yaml"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", URL,
		httpmock.NewBytesResponder(200, []byte("test")))

	filename, err := DownloadFile(newFs, ".", URL)
	assert.NoError(t, err)
	assert.Equal(t, "lucha.yaml", filename)

	httpmock.GetTotalCallCount()
}
