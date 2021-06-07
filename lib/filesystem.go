package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type FileSystem struct {
	fs afero.Fs
}

func NewOsFs() FileSystem {
	var d FileSystem
	d.fs = afero.NewOsFs()
	return d
}

func (f FileSystem) Afero() (afs *afero.Afero) {
	afs = &afero.Afero{Fs: f.fs}
	return
}

//AppendIgnore appends the provided filename to the .luchaignore file
func (f FileSystem) AppendIgnore(filename string) (err error) {
	fn, _ := filepath.Abs(".luchaignore")
	contains, _ := f.Afero().FileContainsBytes(fn, []byte(fmt.Sprintf("%s\n", filename)))
	if !contains {
		file, err := f.fs.OpenFile(fn,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer func() {
			err = file.Close()
		}()

		if _, err := file.WriteString(fmt.Sprintf("%s\n", filename)); err != nil {
			return err
		}
	}
	return
}
func (f FileSystem) LoadIgnore() (err error) {

	filename, _ := filepath.Abs(".luchaignore")
	file, err := f.fs.Open(filename)

	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var ignores []string

	for scanner.Scan() {
		text := scanner.Text()
		if !strings.HasPrefix(text, "#") {
			ignores = append(ignores, text)
		}
	}

	IgnoreFiles = ignores
	return
}

func (f FileSystem) LoadRules(version string) (config Configuration, err error) {
	err = f.LoadIgnore()
	if err != nil {
		return
	}
	filename, _ := filepath.Abs("lucha.yaml")
	yamlFile, err := f.Afero().ReadFile(filename)

	if err != nil {
		return
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return
	}

	err = config.checkVersion(version)

	Rules = config.Lucha.Rules

	return
}

func (f FileSystem) RefreshRules(version string) (config Configuration, err error) {
	luchaDir, _ := LuchaDir()

	exists, _ := f.Afero().DirExists(luchaDir)
	if !exists {
		err = f.Afero().Mkdir(luchaDir, 0700)
		if err != nil {
			return
		}
	}
	_, err = f.DownloadURL("https://raw.githubusercontent.com/devops-kung-fu/lucha/main/lucha.yaml", luchaDir)
	if err != nil {
		return
	}
	return
}

func (f FileSystem) IsTextFile(file ScanFile) bool {
	buf, _ := f.Afero().ReadFile(fmt.Sprintf("%s/%s", file.Path, file.Info.Name()))
	size := 0
	for start := 0; start < len(buf); start += size {
		var r rune
		if r, size = utf8.DecodeRune(buf[start:]); r == utf8.RuneError {
			return false
		}
	}
	return true
}

func ScanFiles(path string) (files []ScanFile, err error) {
	tempFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range tempFiles {
		if !f.IsDir() && !Contains(IgnoreFiles, f.Name()) {
			files = append(files, ScanFile{
				Path: path,
				Info: f,
			})
		}
	}
	return
}

func ScanFilesRecursive(path string) (files []ScanFile, err error) {
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && !Contains(IgnoreFiles, info.Name()) {
				if StartsWith(IgnoreFiles, path) {
					return nil
				}
				files = append(files, ScanFile{
					Path: path,
					Info: info,
				})
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return
}

func LineCounter(r io.Reader) (int, error) {

	var count int
	const lineBreak = '\n'

	buf := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}

	return count, nil
}
