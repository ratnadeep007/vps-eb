package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

// CopyToRemote ...
func CopyToRemote(ipaddr string) {
	fmt.Println("Copying source code")
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err.Error())
	}
	// fmt.Print(a ...interface{})
	clientConfig, err := auth.PrivateKey("root", homedir+"/.ssh/id_rsa", ssh.InsecureIgnoreHostKey())
	if err != nil {
		log.Fatal(err.Error())
	}
	client := scp.NewClient(ipaddr+":22", &clientConfig)
	err = client.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	f, _ := os.Open(GetCurrentDirName() + ".tar")
	defer client.Close()
	defer f.Close()

	err = client.CopyFile(f, "/root/"+f.Name(), "0655")
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println("Copy done")
	os.Remove(GetCurrentDirName() + ".tar")
}
