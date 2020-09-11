package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

//SSHConnect is provide connect to remote server
func SSHConnect() {
	command := "uptime"
	var usInfo userInfo
	var kP keyPath
	key, err := ioutil.ReadFile(kP.privetKey)
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal(err)
	}
	hostKeyCallBack, err := knownhosts.New(kP.knowHost)
	if err != nil {
		log.Fatal(err)
	}
	config := &ssh.ClientConfig{
		User: usInfo.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallBack,
	}
	client, err := ssh.Dial("tcp", usInfo.servIP+":"+usInfo.port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	ss, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer ss.Close()
	var stdoutBuf bytes.Buffer
	ss.Stdout = &stdoutBuf
	ss.Run(command)
	fmt.Println(stdoutBuf.String())
}
