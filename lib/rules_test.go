package lib

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_LuchaDir(t *testing.T) {
	path, err := LuchaDir()
	assert.NoError(t, err)
	assert.Contains(t, path, ".lucha")
}

func Test_LoadRules(t *testing.T) {
	fs := FileSystem{
		fs: afero.NewMemMapFs(),
	}

	createTestFileSystem(fs)

	config, err := LoadRules(fs, "1.0.0", "lucha.yaml")
	versionErr := config.checkVersion("1.0.0")

	assert.NoError(t, err)
	assert.NoError(t, versionErr)
	assert.Equal(t, 1, len(config.Lucha.Rules))
}
