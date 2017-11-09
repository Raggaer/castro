---
name: Download
---

# Download

You can either compile Castro or download an already compiled version.

## Already compiled

You can grab the latest stable release from the [github release page](https://github.com/Raggaer/castro/releases).

You can grab the latest compiled version from [AppVeyor artifacts](https://ci.appveyor.com/project/Raggaer/castro/build/artifacts). You will need to download an executable for your system and the **buildOutput/release.zip** folder.

> **Please note not all AppVeyor builds are tested**

Extract **buildOutput/release.zip** on a folder and add your executable. You are now ready to [install Castro](https://castroaac.org/docs/home/installation)

## Compiling

Building castro from source is a very easy process. You will need to have **Go** installed on your system.

### Installing Go

Make sure you have [Go](https://golang.org/) installed. More information about this topic can be found [here](https://golang.org/doc/install/source). You will need to set-up `GOPATH` and `GOROOT` variables on your system.

### Getting castro

First get `dep` using `go get -u github.com/golang/dep/cmd/dep`; the version control tool Castro uses. 

After installing `dep` download the latest stable source from the releases page. Run `dep ensure` to populate your `vendor` directory.

### Building

You need to run `go build github.com/raggaer/castro`. This will compile and create a castro artifact. Just drop it on your castro directory and you are ready to go.