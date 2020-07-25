package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// CreateArchive -> Create archive from git
func CreateArchive() {
	gitExec, err := exec.LookPath("git")
	if err != nil {
		fmt.Println("git is not found in path")
	}

	getDirectory, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err.Error())
	}
	isGitDirectoryFound := false
	for _, file := range getDirectory {
		if file.Name() == ".git" && file.IsDir() {
			isGitDirectoryFound = true
			break
		}
	}
	if !isGitDirectoryFound {
		println("Current directory is not a git repo")
		os.Exit(0)
	}
	var gitShaBuff bytes.Buffer
	var gitSha string
	gitShaCmd := &exec.Cmd{
		Path:   gitExec,
		Args:   []string{gitExec, "rev-parse", "HEAD"},
		Stdout: &gitShaBuff,
	}
	if err := gitShaCmd.Run(); err != nil {
		println("No commit found")
		os.Exit(0)
	}
	gitSha = gitShaBuff.String()
	println("Using commit: " + gitSha)
	currentDir := GetCurrentDirName()
	archiveCmd := &exec.Cmd{
		Path: gitExec,
		Args: []string{gitExec, "archive", "--format=tar", "--output=" + currentDir + ".tar", "HEAD"},
	}
	if err := archiveCmd.Run(); err != nil {
		println(err.Error())
	}
}
