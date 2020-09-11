package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type userInfo struct {
	user   string
	servIP string
	port   string
}

//IsValidIP parse input ipaddress
func IsValidIP(ip string) bool {
	res := net.ParseIP(ip)
	if res == nil {
		return false
	}
	return true
}

func (usInf *userInfo) getInformToConnect() {
	reader := bufio.NewReader(os.Stdin)
	//User for SSH - server
	usInf.user = ""
	fmt.Println("Enter ur user for connect to remote server")
	usInf.user, _ = reader.ReadString('\n')
	usInf.user = strings.TrimSuffix(usInf.user, "\n")
	fmt.Println("SSH user is " + usInf.user)
	fmt.Println()
	//Server IP
	usInf.servIP = ""
	for IsValidIP(usInf.servIP) != true {
		fmt.Println("Enter ur server ipaddress")
		//fmt.Scanln(usInf.servIP)
		usInf.servIP, _ = reader.ReadString('\n')
		usInf.servIP = strings.TrimSuffix(usInf.servIP, "\n")
		if IsValidIP(usInf.servIP) == false {
			fmt.Println("Incorrect ipaddress")
		}

	}
	fmt.Println("SSH IP address is " + usInf.servIP)
	fmt.Println()
	//Port for SSH-server
	fmt.Println("Enter port to connect")
	usInf.port, _ = reader.ReadString('\n')
	usInf.port = strings.TrimSuffix(usInf.port, "\n")
	fmt.Println("SSH port is " + usInf.port)
	fmt.Println()

}

type keyPath struct {
	privetKey string
	publicKey string
	knowHost  string
}

func (path *keyPath) KeyFolder() {
	reader := bufio.NewReader(os.Stdin)
	/*
		fmt.Println("Enter your path, where programm will genegrate Public Key.\nExample /home/user/folder/key.pub")
		path.publicKey, _ = reader.ReadString('\n')
		path.publicKey = strings.TrimSuffix(path.publicKey, "\n")
		fmt.Println("Path for Public Key is '" + path.publicKey + "'")
		fmt.Println()
	*/
	fmt.Println("Enter your path, where programm will genegrate Keys.\nExample /home/user/folder")
	keyPath := ""
	keyPath, _ = reader.ReadString('\n')
	keyPath = strings.TrimSuffix(keyPath, "\n")
	fmt.Println("Path for Private Key is '" + keyPath + "'")
	fmt.Println()
	fmt.Println("Enter your name for Keys")
	keyName, _ := reader.ReadString('\n')
	keyName = strings.TrimSuffix(keyName, "\n")
	fmt.Println()
	path.privetKey = keyPath + "/" + keyName
	path.publicKey = path.privetKey + ".pub"
	path.knowHost = keyPath + "/known_hosts"

}

func main() {
	var userInfo userInfo
	userInfo.getInformToConnect()
	//fmt.Println(userInfo)
	var keyPath keyPath
	keyPath.KeyFolder()
	fmt.Println(keyPath)
	bitSize := 4096
	privateKey, err := GeneratePrivateKey(bitSize)
	if err != nil {
		log.Fatal(err)
	}
	publicKeyBytes, err := GeneratedPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal(err)
	}
	privateKeysBytes := EncodePrivteKeyToPEM(privateKey)
	err = WriteKeyToFiles(privateKeysBytes, keyPath.privetKey)
	if err != nil {
		log.Fatal(err)
	}
	err = WriteKeyToFiles(publicKeyBytes, keyPath.publicKey)
	if err != nil {
		log.Fatal(err)
	}
	SSHConnect()

}
