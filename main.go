package main

import (
	"theskylab.in/vps-eb/commands"
	"theskylab.in/vps-eb/utils"
)

func main() {
	utils.CreateArchive()
	ipaddr := utils.CreateVM()
	utils.CopyToRemote(ipaddr)
	commands.ExtractAndCheck(ipaddr)
}
