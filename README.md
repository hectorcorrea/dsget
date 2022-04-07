# dsget
A small program to download files (bitstreams) from DataSpace collections

## Sample of usage
Download the [executable](https://github.com/hectorcorrea/dsget/releases) and then issue a command as follows to get a *list of the files* (bitstreams) for the items in a given collection and their sizes. For example:

```
./dsget -collection=261 -get=list
```

The `collection` parameter is required and it must match the ID of an existing collection in DataSpace.

By default `dsget` connects to DataSpace at https://dataspace-dev.princeton.edu but you can pass a `-url` parameter if you want to point to a different DataSpace.

You can also request the *actual files* (bitstreams) to be downloaded with a command as follows:

```
./dsget -collection=261 -get=files
```

by default files are downloaded to the `./data` folder but you can pass the `-data` parameter to request that files are stored at a different location.

When fetching files via the `-get=files` parameter `dsfetch` will download the files and make sure their MD5 hash matches the one reported by DataSpace.

If you stop the program while downloading files and run it again it at a later time `dsget` will recognize the files already on disk and only download files that still need to be downloaded.


## Source code
The source is written in Go. If you want to play with it you can follow these steps (assumming you have Go installed):

```
git clone https://github.com/hectorcorrea/dsget.git
cd dsget
go build
./dsget
```

`main.go` is the entry point and parses the command line options. `downloader.go` has the functionality to connect to DataSpace and download the files.



