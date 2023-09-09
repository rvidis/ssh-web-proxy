package sshservice

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SSHService struct {
	address     string
	hostBaseDir string
	config      *ssh.ClientConfig
}

func NewService(
	hostAddress string,
	hostPort int,
	hostUser string,
	identityFile string,
	hostBaseDir string,
) (*SSHService, error) {
	autMethod, err := getPublicKey(identityFile)
	if err != nil {
		return nil, err
	}

	service := &SSHService{
		config: &ssh.ClientConfig{
			User: hostUser,
			Auth: []ssh.AuthMethod{
				autMethod,
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
		address:     fmt.Sprintf("%s:%d", hostAddress, hostPort),
		hostBaseDir: hostBaseDir,
	}

	return service, nil
}

func getPublicKey(file string) (ssh.AuthMethod, error) {
	var key ssh.Signer

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read key file: %w", err)
	}

	key, err = ssh.ParsePrivateKey(b)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return ssh.PublicKeys(key), nil
}

func (service *SSHService) GetPageContent(filePath string) (string, error) {
	conn, err := ssh.Dial("tcp", service.address, service.config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connecto to [%s]: %v\n", service.address, err)
		os.Exit(1)
	}

	defer conn.Close()

	// Create new SFTP client
	sc, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start SFTP subsystem: %v\n", err)
		os.Exit(1)
	}
	defer sc.Close()

	content, err := getFileContent(sc, path.Join(service.hostBaseDir, filePath))

	if err != nil {
		return "", err
	}

	return content, nil
}

func getFileContent(sc *sftp.Client, remoteFile string) (string, error) {
	fmt.Printf("getting remote file [%s] ...\n", remoteFile)
	file, err := sc.OpenFile(remoteFile, os.O_RDONLY)

	if err != nil {
		return "", fmt.Errorf("unable to open remote file: %v\n", err)
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("unable to get remote file stat: %v\n", err)
	}

	fileBytes := make([]byte, stat.Size())
	_, err = file.Read(fileBytes)
	if err != nil {
		return "", fmt.Errorf("unable to read remote file: %v\n", err)
	}

	return string(fileBytes), nil
}
