---
name: Installation
---

# Installation

You can either compile Castro or download an already compiled version.

## Already compiled

You can grab the latest stable release from the [github release page](https://github.com/Raggaer/castro/releases). 

However you can get the latest commit compiled at AppVeyor. Please note all AppVeyor releases are not tested.

## Compiling

Building castro from source is a very easy process. You will need to have **Go** installed on your system.

### Installing Go

Make sure you have [Go](https://golang.org/) installed. More information about this topic can be found [here](https://golang.org/doc/install/source). You will need to set-up `GOPATH` and `GOROOT` variables on your system.

### Getting castro

You can run `go get github.com/raggaer/castro` to get the latest commit. However it is recommended that you download the latest version source and extract it on `GOPATH/src/github.com/raggaer/castro`

### Building

You need to run `go build github.com/raggaer/castro`. This will compile and create a castro artifact. Just drop it on your castro directory and you are ready to go.