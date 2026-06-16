package baremetal

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *ssh.Client
	host   string
	port   string
	user   string
}

func Connect(host, port, user, keyPath, password string) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	if keyPath != "" {
		keyBytes, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("ssh: failed to read key file %s: %w", keyPath, err)
		}
		signer, err := ssh.ParsePrivateKey(keyBytes)
		if err != nil {
			return nil, fmt.Errorf("ssh: failed to parse private key: %w", err)
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else if password != "" {
		config.Auth = []ssh.AuthMethod{ssh.Password(password)}
	} else {
		return nil, fmt.Errorf("ssh: no authentication method provided (set keyPath or password)")
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("ssh: failed to dial %s: %w", addr, err)
	}

	log.Printf("[ssh] Connected to %s@%s", user, addr)
	return &SSHClient{client: client, host: host, port: port, user: user}, nil
}

func (c *SSHClient) Run(cmd string) (string, error) {
	sess, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("ssh: failed to create session: %w", err)
	}
	defer sess.Close()

	output, err := sess.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("ssh: command failed: %s\noutput: %s", err, string(output))
	}
	return string(output), nil
}

func (c *SSHClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
