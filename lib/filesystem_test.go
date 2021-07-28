package lib

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func createTestFileSystem(fs FileSystem) {
	file, _ := fs.fs.OpenFile("test.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString(fmt.Sprintf("%s\n", "test"))
	fs.fs.MkdirAll("foo/bar", 0644)

	file, _ = fs.fs.OpenFile("foo/foo.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString(fmt.Sprintf("%s\n", "test"))

	file, _ = fs.fs.OpenFile("foo/bar/bar.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString(fmt.Sprintf("%s\n", "test"))

	config := Configuration{
		Version: "1.0.0",
		Lucha: Lucha{
			Rules: []Rule{
				{
					Code:        "DFKM001",
					Name:        "test",
					Description: "Matches if the word `test` (case insensitive) exists in a test file",
					Message:     "Line may contains the word test",
					Attribution: "DKFM",
					Regex:       "(?i)\\btest\\b",
					Severity:    0,
				},
			},
		},
	}

	lf, _ := yaml.Marshal(config)

	fs.Afero().WriteFile("lucha.yaml", lf, 0644)
}

func Test_NewOsFs(t *testing.T) {
	f := FileSystem{}

	var i interface{} = NewOsFs()
	var fs interface{} = afero.NewOsFs()

	assert.IsType(t, f, i)
	assert.IsType(t, fs, NewOsFs().fs)
}

func TestFileSystem_BuildFileList(t *testing.T) {
	fs := FileSystem{
		fs:         afero.NewMemMapFs(),
		SearchPath: ".",
	}

	createTestFileSystem(fs)
	fs.Afero().WriteFile(".luchaignore", []byte("test.txt"), 0666)

	//test if no ignore as well

	b, err := fs.Afero().Exists("lucha.yaml")
	assert.True(t, b)
	assert.NoError(t, err)

	// scanFiles, err := BuildFileList(fs)

	// assert.NoError(t, err, "There should be no error")
	// assert.Equal(t, 3, len(scanFiles), "Expecting .luchaignore, lucha.yaml, and test.txt")
}

// 	scanFiles, err = f.BuildFileList(".", true)

// 	assert.NoError(t, err, "There should be no error")
// 	assert.Equal(t, 5, len(scanFiles), "Expecting 5 files")

// 	_, err = f.BuildFileList("...", false)
// 	assert.Error(t, err, "There should be an error because the folder ... shouldn't exist")

// 	_, err = f.BuildFileList("...", true)
// 	assert.Error(t, err, "There should be an error because the folder ... shouldn't exist")
// }

func TestFileSystem_AbsoluteSearchPath(t *testing.T) {
	fs := FileSystem{
		fs:         afero.NewMemMapFs(),
		SearchPath: ".",
	}
	assert.Contains(t, fs.AbsoluteSearchPath(), "/lucha/lib")
}

func Test_shouldIgnoreDir(t *testing.T) {
	fs := FileSystem{
		fs:         afero.NewMemMapFs(),
		SearchPath: ".",
	}
	fs.Afero().Mkdir(".git", 0644)
	fi, _ := fs.Afero().ReadDir(fs.SearchPath)
	assert.Len(t, fi, 1)

	dir := fi[0]

	shouldIgnore := shouldIgnoreDir(fs, dir, "")
	assert.True(t, shouldIgnore)

	fs.IncludeGit = true
	shouldIgnore = shouldIgnoreDir(fs, dir, "")
	assert.False(t, shouldIgnore)
}

func Test_isUTF8(t *testing.T) {
	fs := FileSystem{
		fs:         afero.NewMemMapFs(),
		SearchPath: ".",
	}

	fs.Afero().WriteFile("test.txt", []byte("test string"), 0666)
	fi, _ := fs.fs.OpenFile("test.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	b := isUTF8(fs, fi)
	assert.True(t, b)
}
