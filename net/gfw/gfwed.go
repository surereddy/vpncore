package gfw

import (
	"bufio"
	"encoding/base64"
	"io"
	"os"
	"regexp"
	"time"
	"net/http"
)

const gfwlistURL = "https://yii.li/gfwlist"

func downloadRemoteContent(remoteLink string) (io.ReadCloser, error) {
	response, err := http.Get(remoteLink)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}


func CreateGFWList(url string, file_path string) (gfwlist *ItemSet, err error) {
	if url == "" {
		url = gfwlistURL
	}

	gfwlist = NewItemSet(file_path, 4000)

	var content io.ReadCloser
	for err = os.ErrNotExist; err != nil; time.Sleep(5 * time.Second) {
		content, err = downloadRemoteContent(url)
	}
	defer content.Close()

	decoder := base64.NewDecoder(base64.StdEncoding, content)

	commentPattern, _ := regexp.Compile(`^\!|\[|^@@|^\d+\.\d+\.\d+\.\d+`)
	domainPattern, _ := regexp.Compile(`([\w\-\_]+\.[\w\.\-\_]+)[\/\*]*`)
	scanner := bufio.NewScanner(decoder)
	scanner.Split(bufio.ScanLines)
	gfwlist.Lock()
	for scanner.Scan() {
		t := scanner.Text()
		if commentPattern.MatchString(t) {
			continue
		}
		ss := domainPattern.FindStringSubmatch(t)
		if len(ss) > 1 {
			gfwlist.AddItem(ss[1])
		}
	}
	gfwlist.Unlock()
	gfwlist.Save()
	return
}
