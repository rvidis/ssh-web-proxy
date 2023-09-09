package sshservice

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/pkg/sftp"
	gossh "golang.org/x/crypto/ssh"
)

func TestGetPageContent(t *testing.T) {
	address := "0.0.0.0"
	port := 2022
	idrsaPrivate := filepath.Join("test_resources", "id_rsa_insecure")
	idrsaPublic := filepath.Join("test_resources", "id_rsa_insecure.pub")
	user := "app-admin"

	// Create a mock sftp server to test the retrieval of file content

	mockSftpServer, err := NewMockSftpServer(
		fmt.Sprintf("%s:%d", address, port),
		idrsaPublic,
		user,
	)
	if err != nil {
		t.Fatalf("unable to create mock server: %v\n", err)
	}

	go func() {
		mockSftpServer.ListenAndServe()
	}()
	t.Cleanup(func() {
		mockSftpServer.Close()
	})
	time.Sleep(1 * time.Second)

	// Create system under test

	service, err := NewService(
		address,
		port,
		user,
		idrsaPrivate,
		"test_site",
	)
	assertNoError(t, err)

	// Cetrieve the page content

	pageContent, err := service.GetPageContent("index.html")
	assertNoError(t, err)

	// Assert

	expectedContent := "<html>Test remote file</html>"
	content := strings.TrimSpace(pageContent)
	if content != expectedContent {
		t.Errorf("Expected: \"%v\" but got \"%v\"", expectedContent, content)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("TestGetPageContent failed: %v", err)
	}
}

type MockSftpServer struct {
	server *ssh.Server
}

func NewMockSftpServer(addr string, idrsaPublic string, user string) (*MockSftpServer, error) {
	fileIdrsa, err := ioutil.ReadFile(idrsaPublic)
	if err != nil {
		return nil, fmt.Errorf("could not fing public rsa %s: %v", idrsaPublic, err)
	}

	return &MockSftpServer{
		server: &ssh.Server{
			Addr: addr,
			SubsystemHandlers: map[string]ssh.SubsystemHandler{
				"sftp": SftpHandler,
			},
			PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
				if ctx.User() != user {
					fmt.Printf("user does not match, expected \"%s\" but got \"%s\"\n", user, ctx.User())
					return false
				}
				marshalledKey := gossh.MarshalAuthorizedKey(key)
				if bytes.Equal(marshalledKey, fileIdrsa) {
					return true
				}
				fmt.Println("id_rsa_insecure test keys don't match")
				return false
			},
		},
	}, nil
}

func SftpHandler(sess ssh.Session) {
	debugStream := ioutil.Discard
	str, _ := os.Getwd()
	rootDirectory := filepath.Join(str, "test_resources")
	serverOptions := []sftp.ServerOption{
		sftp.WithDebug(debugStream),
		sftp.WithServerWorkingDirectory(rootDirectory),
	}
	server, err := sftp.NewServer(
		sess,
		serverOptions...,
	)
	if err != nil {
		log.Printf("sftp server init error: %s\n", err)
		return
	}
	if err := server.Serve(); err == io.EOF {
		server.Close()
		fmt.Println("sftp client exited session.")
	} else if err != nil {
		fmt.Println("sftp server completed with error:", err)
	}
}

func (h *MockSftpServer) ListenAndServe() error {
	return h.server.ListenAndServe()
}

func (h *MockSftpServer) Close() error {
	return h.server.Close()
}
