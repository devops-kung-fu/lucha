package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var (
	f FileSystem = FileSystem{
		fs: afero.NewMemMapFs(),
	}
	version string = "1.0.0"
	config  Configuration
)

func createTestFileSystem() {
	file, _ := f.fs.OpenFile("test.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString(fmt.Sprintf("%s\n", "test"))

	f.AppendIgnore("test")
	config := Configuration{
		Version: version,
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

	filename, _ := filepath.Abs("lucha.yaml")
	f.Afero().WriteFile(filename, lf, 0644)

}

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

func TestFileSystem_LoadIgnore(t *testing.T) {
	createTestFileSystem()
	err := f.LoadIgnore()
	assert.NoError(t, err, "Should be no error loading the .luchaignore file")
}

func TestFileSystem_LuchaRulesFile(t *testing.T) {
	createTestFileSystem()
	pwd, _ := os.Getwd()
	filename := fmt.Sprintf("%s/lucha.yaml", pwd)

	file, err := f.LuchaRulesFile("lucha.yaml")

	assert.NoError(t, err, "There should be no error")
	assert.Equal(t, file, filename, "Paths should be equal")
}

func TestFileSystem_LoadRules(t *testing.T) {
	createTestFileSystem()
	pwd, _ := os.Getwd()
	filename := fmt.Sprintf("%s/lucha.yaml", pwd)
	config, err := f.LoadRules(version, filename)
	versionErr := config.checkVersion(version)

	assert.NoError(t, err, "There should be no error loading lucha.yaml")
	assert.NoError(t, versionErr, "Version should have matched 1.0.0")
	assert.Equal(t, 1, len(config.Lucha.Rules), "There should have only been one rule")
}

func TestFileSystem_IsTextFile(t *testing.T) {
	createTestFileSystem()

	info, _ := f.Afero().Stat("test.txt")
	scanFile := ScanFile{
		Path:   ".",
		Info:   info,
		Issues: []Issue{},
	}
	test := f.IsTextFile(scanFile)

	assert.True(t, test, "test.txt should have been set up as a text file")
}

func TestScanFiles(t *testing.T) {

}

func TestScanFilesRecursive(t *testing.T) {

}
