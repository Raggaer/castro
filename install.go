package main

import "os"

func isInstalled() bool {
	_, err := os.Stat("config.toml")
	return err == nil
}
