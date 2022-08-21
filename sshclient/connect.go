package sshclient

import (
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func SshToGit(filePath string) {
	var hostKey ssh.PublicKey

	// path to private key file
	dirname, err := os.UserHomeDir()
	key, err := ioutil.ReadFile(dirname + filePath)

	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: "user",
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", "git@github.com:22", config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}

	defer client.Close()
}
