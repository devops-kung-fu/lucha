package lib

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_LoadIgnore(t *testing.T) {

	fs := FileSystem{
		fs:         afero.NewMemMapFs(),
		SearchPath: ".",
	}

	createTestFileSystem(fs)
	_, err := LoadIgnore(fs)
	assert.Error(t, err)

	fs.Afero().WriteFile(".luchaignore", []byte("test.txt"), 0666)

	ignores, err := LoadIgnore(fs)
	assert.NoError(t, err)
	assert.Len(t, ignores, 1)

}
