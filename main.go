package main

import (
	"flag"
	"fmt"
)

var url, get, data string
var collectionId int
var downloader Downloader

func init() {
	flag.IntVar(&collectionId, "collection", 0, "Collection ID (required)")
	flag.StringVar(&url, "url", "https://dataspace-dev.princeton.edu", "DataSpace URL")
	flag.StringVar(&get, "get", "list", "list or files")
	flag.StringVar(&data, "data", "./data", "Folder where data will be downloaded to")

	flag.Usage = func() {
		fmt.Println("Downloads files (bitstreams) from a DataSpace collection")
		flag.PrintDefaults()
		fmt.Println()
	}

	flag.Parse()
}

func main() {
	downloader = NewDownloader(url, data)

	if collectionId == 0 {
		flag.Usage()
		fmt.Println("Error: No collection ID indicated")
		return
	}

	if get == "list" {
		listItems()
	} else if get == "files" {
		downloadItems()
	} else {
		flag.Usage()
		fmt.Printf("Error: Unknown get value (%s)\r\n", get)
	}
}

// List the items in the collection
func listItems() {
	items := downloader.GetItems(collectionId)
	totalSizeBytes := 0
	fmt.Printf("ID, Handle, FileName, FileSize, FileLink\r\n")
	for _, item := range items {
		for _, bitstream := range item.Bitstreams {
			fmt.Printf("%d, %s, %s, %d, %s, %s, %s\r\n",
				item.Id, item.Handle,
				bitstream.Name, bitstream.SizeBytes, bitstream.RetrieveLink,
				bitstream.CheckSum.Algorithm, bitstream.CheckSum.Value)
			totalSizeBytes += bitstream.SizeBytes
		}
	}

	var totalSizeGB float32 = float32(totalSizeBytes) / (1024 * 1024 * 1024)
	fmt.Printf("Total size: %d bytes (%0.2f GB)\r\n", totalSizeBytes, totalSizeGB)
}

// Downloads the bitstream for the items in the collection
func downloadItems() {
	fmt.Printf("Downloading from: %s\r\n", downloader.DataSpaceUrl)
	items := downloader.GetItems(collectionId)
	for _, item := range items {
		fmt.Printf("Processing item %d, %s, %s\r\n", item.Id, item.Handle, item.Name)
		downloader.DownloadItem(item)
	}
}
