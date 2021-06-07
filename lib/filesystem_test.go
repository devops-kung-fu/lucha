package lib

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewOsFs(t *testing.T) {
	f := FileSystem{}

	var i interface{} = NewOsFs()
	var fs interface{} = afero.NewOsFs()

	assert.IsType(t, f, i, "Not returning a FileSystem struct")
	assert.IsType(t, fs, NewOsFs().fs, "fs should be an afero.OsFs")
}
