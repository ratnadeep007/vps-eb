package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"theskylab.in/vps-eb/utils"
)

// ExtractAndCheck ...
func ExtractAndCheck(ipaddr string) {
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			publicKey(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", ipaddr+":22", config)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()
	// extractCmd := "mkdir " + utils.GetCurrentDirName() + " && tar -xvf " + utils.GetCurrentDirName() + ".tar -C ./" + utils.GetCurrentDirName()
	folderName := utils.GetCurrentDirName()
	extractCmd := fmt.Sprintf(`if [ -d %s ]; then
rm -rf %s;
fi
mkdir %s
tar -xvf %s.tar -C ./%s
`, folderName, folderName, folderName, folderName, folderName)
	// println(extractCmd)
	// extractCmd := "ls && pwd"
	RunCommand(extractCmd, conn)
	lang := CheckLang()
	if lang == "go" {
		InstallGo(conn)
		BuildGo(conn)
		RunExectuable(conn)
	}

}

// CheckLang ...
func CheckLang() string {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, file := range files {
		if strings.Contains(file.Name(), "go.mod") {
			return "go"
		}
		if strings.Contains(file.Name(), ".go") {
			return "go"
		}
		if strings.Contains(file.Name(), "package.json") {
			return "node"
		}
	}
	return ""
}

// RunCommand ...
func RunCommand(cmd string, conn *ssh.Client) {
	sess, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()
	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stderr, sessStderr)
	err = sess.Run(cmd) // eg., /usr/bin/whoami
	if err != nil {
		panic(err)
	}
}

func publicKey() ssh.AuthMethod {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err.Error())
	}
	key, err := ioutil.ReadFile(homedir + "/.ssh/id_rsa")
	if err != nil {
		log.Fatal(err.Error())
	}
	singer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal(err.Error())
	}
	return ssh.PublicKeys(singer)
}
