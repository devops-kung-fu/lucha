package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/spf13/afero"
)

var (
	//Rules contains the loaded rules from lucha.yaml
	Rules []Rule
	//Ignores contains the names of files that shouldn't be processed from the .luchaignore file

)

//FileSystem encapsulates the Afero fs Filesystem
type FileSystem struct {
	fs         afero.Fs
	SearchPath string
	Recursive  bool
	IncludeGit bool
}

//AbsoluteSearchPath returns the the absolute path for the (possibly) relative search path
func (fs FileSystem) AbsoluteSearchPath() string {
	path, err := filepath.Abs(fs.SearchPath)
	if err != nil {
		fmt.Println(err)
	}
	return path
}

//NewOsFs returns a new local os file system
func NewOsFs() FileSystem {
	return FileSystem{
		fs: afero.NewOsFs(),
	}
}

//Afero returns the Afero system
func (fs FileSystem) Afero() *afero.Afero {
	return &afero.Afero{Fs: fs.fs}
}

//IsTextFile examines a file and returns true if the file is UTF-8
func isUTF8(fs FileSystem, file afero.File) bool {
	buf, _ := fs.Afero().ReadFile(file.Name()) //fmt.Sprintf("%s/%s", file.Path, file.Info.Name()))
	size := 0
	for start := 0; start < len(buf); start += size {
		var r rune
		if r, size = utf8.DecodeRune(buf[start:]); r == utf8.RuneError {
			return false
		}
	}
	return true
}

func shouldIgnore(file string, ignoreList []string) (ignore bool) {
	var absIgnore []string

	for _, i := range ignoreList {
		path, _ := filepath.Abs(i)
		absIgnore = append(absIgnore, path)
	}
	return !matchIgnore(absIgnore, file)
}

func matchIgnore(s []string, str string) (matches bool) {

	for _, v := range s {
		if v == str {
			matches = true
			return
		}
		matches = strings.HasPrefix(str, v)
	}

	return
}

func shouldIgnoreDir(fs FileSystem, f os.FileInfo, path string) bool {
	if f.IsDir() && f.Name() == ".git" {
		return !fs.IncludeGit
	}
	if (f.IsDir() && !fs.Recursive) && fs.AbsoluteSearchPath() != path {
		return true
	}
	return false
}

//BuildFileList gathers all of the files from the searchpath down the folder tree
func BuildFileList(fs FileSystem) (fileList []string, err error) {

	ignores, _ := LoadIgnore(fs)
	err = fs.Afero().Walk(fs.AbsoluteSearchPath(), func(path string, f os.FileInfo, err error) error {
		if shouldIgnoreDir(fs, f, path) {
			return filepath.SkipDir
		}
		if shouldIgnore(path, ignores) {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		return
	}

	return
}
