//Package lib Functionality for the Hookz CLI
package lib

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/dustin/go-humanize"
)

//WriteCounter encapsulates the total number of bytes captured and rendered
type WriteCounter struct {
	Total    uint64
	FileName string
}

//Write increments the total number of bytes and prints progress to STDOUT
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

//PrintProgress prints the current download progress to STDOUT
func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading %s... %s complete", wc.FileName, humanize.Bytes(wc.Total))
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(fs FileSystem, filepath string, URL string) (filename string, err error) {
	_, err = url.ParseRequestURI(URL)
	if err != nil {
		return
	}
	filename = path.Base(URL)
	fullFileName := fmt.Sprintf("%s/%s", filepath, filename)
	out, err := fs.Afero().Create(fmt.Sprintf("%s.tmp", fullFileName))
	if err != nil {
		return
	}
	defer func() {
		err = out.Close()
	}()

	resp, err := http.Get(URL)
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()

	counter := &WriteCounter{FileName: filename}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return
	}

	fmt.Print("\n")

	err = os.Rename(fmt.Sprintf("%s.tmp", fullFileName), fullFileName)
	if err != nil {
		return
	}

	err = os.Chmod(fullFileName, 0777)
	filename = fullFileName

	return
}
