package lib

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var (
	f FileSystem = FileSystem{
		fs: afero.NewMemMapFs(),
	}
	version string = "1.0.0"
	config  Configuration
)

func TestNewOsFs(t *testing.T) {
	f := FileSystem{}

	var i interface{} = NewOsFs()
	var fs interface{} = afero.NewOsFs()

	assert.IsType(t, f, i, "Not returning a FileSystem struct")
	assert.IsType(t, fs, NewOsFs().fs, "fs should be an afero.OsFs")
}

func TestFileSystem_AppendIgnore(t *testing.T) {
	err := f.AppendIgnore("TEST")
	assert.NoError(t, err, "No error should come out of this method")

	path, _ := os.Getwd()
	fullFileName := fmt.Sprintf("%s/.luchaignore", path)
	contains, _ := f.Afero().FileContainsBytes(fullFileName, []byte("TEST"))
	assert.True(t, contains, ".luchaignore file should have the phrase `TEST` in it")
}
