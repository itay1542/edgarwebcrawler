package edgarwebcrawler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type ItemType uint8

const (
	File ItemType = iota
	Dir
)

type EdgarItem struct {
	//file or dir
	ItemType         string `json:"type"`
	Href, Size, Name string
	LastModified     string `json:"last-modified"`
}

type EdgarDirectory struct {
	Name      string      `json:"name"`
	ParentDir string      `json:"parent-dir"`
	Item      []EdgarItem `json:"item"`
}

type IdxDownloader interface {
	Download(baseEdgarDirectoryUrl string, recursive bool) error
}

//this implements IdxDownloader
type FormFourIdxDownloader struct {
	wg sync.WaitGroup
}

func NewIdxDownloader() *FormFourIdxDownloader {
	return &FormFourIdxDownloader{
		wg: sync.WaitGroup{},
	}
}

func (dl *FormFourIdxDownloader) Download(baseDir, baseEdgarDirectoryUrl string, recursive bool,
	startFromYear uint16) error {
	index, err := getIndex(baseEdgarDirectoryUrl)
	if err != nil {
		return err
	}
	filterYearDirectories(index, startFromYear)
	errc := make(chan error)
	wgDone := make(chan bool)
	if recursive {
		dl.wg.Add(1)
		go dl.crawlDirectory(*index, baseDir, baseEdgarDirectoryUrl, errc)
	}
	go func() {
		dl.wg.Wait()
		wgDone <- true
		close(wgDone)
	}()
	select {
	case <-wgDone:
		return nil
	case e := <-errc:
		close(errc)
		return e
	}
}

func getIndex(url string) (*EdgarDirectory, error) {
	req_url := fmt.Sprintf("%s/index.json", url)
	fmt.Printf("Requesting %s \n", req_url)
	request, err := http.NewRequest("GET", req_url, nil)

	request.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("got error when requesting resource %s \n", err)
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	type tempParsedType struct {
		Directory EdgarDirectory `json:"directory"`
	}
	var tempParsed tempParsedType
	json.Unmarshal(data, &tempParsed)
	return &tempParsed.Directory, nil
}

func (dl *FormFourIdxDownloader) crawlDirectory(directory EdgarDirectory,
	destdir, baseEdgarDirectoryUrl string, errc chan error) {
	defer dl.wg.Done()
	// sleep is needed to not exceed edgar request rate limit
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	for _, item := range directory.Item {
		fmt.Printf("item found: %s", item.Name)
		if item.ItemType == "dir" {
			//this is the full path from /full-index to the directory
			newContextDir := filepath.Join(destdir, directory.Name, item.Name)
			fmt.Printf("creating new directory: %s \n", newContextDir)
			err := os.MkdirAll(newContextDir, os.ModePerm)
			if err != nil {
				errc <- err
				return
			}
			//directory name starts with full-index so we remove it
			newUrl := baseEdgarDirectoryUrl + directory.Name[10:] + item.Name
			fmt.Printf("new Edgar url: %s \n", newUrl)
			newIndex, err := getIndex(newUrl)
			if err != nil {
				errc <- err
				return
			}
			dl.wg.Add(1)
			go dl.crawlDirectory(*newIndex, destdir, baseEdgarDirectoryUrl, errc)
		} else if item.ItemType == "file" && item.Name == "form.idx" {
			fmt.Printf("Found form.idx \n")
			fileUrl := baseEdgarDirectoryUrl + directory.Name[10:] + item.Href
			err := downloadFile(filepath.Join(destdir, directory.Name, item.Name), fileUrl)
			if err != nil {
				errc <- err
			}

		}
	}
}

func downloadFile(filepath string, url string) error {
	fmt.Printf("getting file from url: %s \n", url)
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if strings.Contains(resp.Header.Get("Content-Type"), "html") {
		fmt.Println("got html response type from file url, retrying")
		resp.Body.Close()
		sleepDuration := time.Duration(rand.Intn(5)) * time.Second
		fmt.Printf("sleeping for %d nanoseconds... \n", sleepDuration)
		time.Sleep(sleepDuration)
		downloadFile(filepath, url)
	} else {
		defer resp.Body.Close()
		// Create the file
		out, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		return err
	}
	return nil
}
