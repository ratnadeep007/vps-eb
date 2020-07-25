package utils

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/digitalocean/godo"
)

// CreateVM -> Create a vm in vps
func CreateVM() string {
	fingerprint := getFingerPrint()
	machineName := GetCurrentDirName()
	fingerPrintCreate := godo.DropletCreateSSHKey{
		Fingerprint: fingerprint,
	}
	var fingerprints []godo.DropletCreateSSHKey
	fingerprints = append(fingerprints, fingerPrintCreate)
	createRequest := &godo.DropletCreateRequest{
		Name:    machineName,
		Region:  "blr1",
		Size:    "s-1vcpu-1gb",
		SSHKeys: fingerprints,
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-18-04-x64",
		},
	}
	ctx := context.TODO()
	client := godo.NewFromToken(os.Getenv("DIGITALOCEAN_SECRET_KEY"))
	newDroplet, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Machine created, waiting to get machine ready")
	for {
		getDroplet, _, err := client.Droplets.Get(ctx, newDroplet.ID)
		if err != nil {
			log.Fatal(err.Error())
		}
		if getDroplet.Status == "active" {
			break
		}
		time.Sleep(10 * time.Second)
	}
	for {
		getDroplet, _, err := client.Droplets.Get(ctx, newDroplet.ID)
		if err != nil {
			log.Fatal(err.Error())
		}
		ipaddr, err := getDroplet.PublicIPv4()
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("Checking machine ip")
		isUp := CheckPing(ipaddr)
		if isUp {
			break
		}
		time.Sleep(10 * time.Second)
	}
	getDroplet, _, err := client.Droplets.Get(ctx, newDroplet.ID)
	if err != nil {
		log.Fatal(err.Error())
	}
	ipaddr, err := getDroplet.PublicIPv4()
	fmt.Println("Machine ip: " + ipaddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Machine ready")
	return ipaddr
}

func getFingerPrint() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		println("ssh key file not found in home directory")
	}
	key, err := ioutil.ReadFile(homeDir + "/.ssh/id_rsa.pub")
	if err != nil {
		println("ssh key file not found in home directory")
	}
	parts := strings.Fields(string(key))
	if len(parts) < 2 {
		log.Fatal("bad key")
	}
	k, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Fatal(err)
	}
	fp := md5.Sum([]byte(k))
	var fpString []string
	for _, val := range fp {
		// println(val)
		fpString = append(fpString, fmt.Sprintf("%02x", val))
	}
	// fmt.Print("MD5:")
	// for i, b := range fp {
	// 	fmt.Printf("%02x", b)
	// 	if i < len(fp)-1 {
	// 		fmt.Print(":")
	// 	}
	// }
	return strings.Join(fpString[:], ":")
}
