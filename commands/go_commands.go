package commands

import (
	"fmt"

	"golang.org/x/crypto/ssh"
	"theskylab.in/vps-eb/utils"
)

// InstallGo ...
func InstallGo(conn *ssh.Client) {
	installCmd := fmt.Sprintf(`wget -q https://storage.googleapis.com/golang/getgo/installer_linux;
chmod +x installer_linux;
./installer_linux;
source /root/.bash_profile;
rm installer_linux; 
`)
	RunCommand(installCmd, conn)
}

// BuildGo ...
func BuildGo(conn *ssh.Client) {
	folderName := utils.GetCurrentDirName()
	buildCmd := fmt.Sprintf(`source /root/.bash_profile;
cd %s && ls && pwd \
&& go mod download \
&& go build -o ../%s-exec *.go`, folderName, folderName)
	RunCommand(buildCmd, conn)
}

// RunExectuable ...
func RunExectuable(conn *ssh.Client) {
	folderName := utils.GetCurrentDirName()
	execCmd := fmt.Sprintf(`chmod +x %s-exec;
echo "[Unit]
Description=%s Service by EB

[Service]
ExecStart=/root/%s-exec

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/%s.service;
systemctl start %s
`, folderName, folderName, folderName, folderName, folderName)
	RunCommand(execCmd, conn)
}
