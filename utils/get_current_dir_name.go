package utils

import (
	"os"
	"strings"
)

// GetCurrentDirName -> Get Current directory Name
func GetCurrentDirName() string {
	currentDirPath, _ := os.Getwd()
	currentDirArray := strings.Split(currentDirPath, "/")
	currentDir := currentDirArray[len(currentDirArray)-1]
	return currentDir
}
