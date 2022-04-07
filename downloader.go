package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Checksum struct {
	Value     string `json:"value"`
	Algorithm string `json:"checkSumAlgorithm"`
}
type Bitstream struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	BundleName   string `json:"bundleName"`
	CheckSum     Checksum
	Format       string `json:"format"` // MPEG
	RetrieveLink string `json:"retrieveLink"`
	SizeBytes    int    `json:"sizeBytes"`
	SequenceId   int    `json:"sequenceId"`
}

type Item struct {
	Handle     string `json:"handle"`
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Bitstreams []Bitstream
}

// Downloader implements the functionality to download files from DataSpace.
type Downloader struct {
	DataSpaceUrl string
	LocalDir     string
}

// Create a new Downloader instance.
func NewDownloader(dataSpaceUrl string, dataDir string) Downloader {
	downloader := Downloader{
		DataSpaceUrl: dataSpaceUrl,
		LocalDir:     dataDir,
	}
	return downloader
}

// Get the item information for a given collection Id.
func (d Downloader) GetItems(collectionId int) []Item {
	// TODO: handle more than 100 items in a collection
	// expand=all
	queryString := "limit=100&offset=0&expand=bitstreams"
	url := fmt.Sprintf("%s/rest/collections/%d/items?%s", d.DataSpaceUrl, collectionId, queryString)
	bytes, err := httpGet(url)
	if err != nil {
		log.Fatal(err)
	}

	var items []Item
	json.Unmarshal(bytes, &items)
	return items
}

// Downloads the files (bitstreams) for a given item
func (d Downloader) DownloadItem(item Item) {

	// Make sure we have a directory for the item
	itemPath := filepath.Join(d.LocalDir, item.Handle)
	if !dirExist(itemPath) {
		err := createDir(itemPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Download the individual files (bitstreams)
	for _, b := range item.Bitstreams {
		filename := filepath.Join(d.LocalDir, item.Handle, b.Name)

		// file already on disk
		if fileExist(filename) && fileMD5(filename) == b.CheckSum.Value {
			fmt.Printf("\tfile: %s - already on disk\r\n", b.Name)
			continue
		}

		// fetch the file from DataSpace
		fmt.Printf("\tdownloading: %s (%d bytes)\r\n", b.Name, b.SizeBytes)
		url := fmt.Sprintf("%s/rest%s", d.DataSpaceUrl, b.RetrieveLink)
		bytes, err := httpGet(url)
		if err != nil {
			log.Fatal(err)
		}

		// save it
		writeFile(filename, bytes)
		if err != nil {
			log.Fatal(err)
		}

		// make sure local MD5 is OK
		if fileMD5(filename) != b.CheckSum.Value {
			log.Fatal(errors.New("MD5 mismatch"))
		}
	}
}

func httpGet(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return []byte{}, err
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	// log.Printf("%s\r\n", bytes)
	return bytes, err
}

func fileExist(name string) bool {
	file, err := os.Open(name)
	if os.IsNotExist(err) {
		return false
	}
	defer file.Close()
	return true
}

func dirExist(name string) bool {
	return fileExist(name)
}

func createDir(fullpath string) error {
	return os.MkdirAll(fullpath, os.ModePerm)
}

func writeFile(filename string, bytes []byte) error {
	return ioutil.WriteFile(filename, bytes, 0644)
}

// Calculates the MD5 of a file
// Source: https://pkg.go.dev/crypto/md5@master#example-Sum
func fileMD5(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
