package sync

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SyncClient struct {
	SftpClient *sftp.Client
	Conn       *ssh.Client
	Config     SyncConfig
}

func NewClient(cfg SyncConfig) (SyncClient, error) {
	var c SyncClient

	key, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		log.Println("Error reading ssh key")
		return c, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Println("Error parsing private key")
		return c, err
	}

	sshConfig := &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // WARN: replace this in prod
	}

	// Connect to the SSH server
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), sshConfig)
	if err != nil {
		fmt.Println("Failed to connect to SSH server:", err)
		return c, err
	}

	// Open SFTP session
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Println("Failed to open SFTP session:", err)
		return c, err
	}

	return SyncClient{
		SftpClient: sftpClient,
		Conn:       conn,
		Config:     cfg,
	}, nil
}

func (c SyncClient) Close() error {
	if err := c.SftpClient.Close(); err != nil {
		log.Println("Failed to close sftp connection: ", err)
	}

	if err := c.Conn.Close(); err != nil {
		log.Println("Failed to close ssh connection: ", err)
		return err
	}

	return nil
}
