package system

import (
	"log"
	"os"
	"path/filepath"
)

func GetSharedLocalDir() string {
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("Error acquiring user home directory: ", err)
		}
		return filepath.Join(home, ".local", "share")
	}
	return dataDir
}

func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Fatal("Error closing File: ", err)
	}
}
