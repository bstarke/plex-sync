package cmd

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type SftpClient struct {
	host, user string
	port       int
	*sftp.Client
}

func NewConn(host, user string, port int) (client *SftpClient, err error) {
	switch {
	case `` == strings.TrimSpace(host),
		`` == strings.TrimSpace(user),
		0 >= port || port > 65535:
		return nil, errors.New("invalid parameters")
	}
	client = &SftpClient{
		host: host,
		user: user,
		port: port,
	}
	err = client.connect()
	if err != nil {
		return nil, err
	}
	return client, err
}

func (sc *SftpClient) connect() (err error) {
	privateKeyFile, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "id_ed25519"))
	if err != nil {
		log.Panicf("error opening private key: %v", err)
	}
	pemBytes, err := io.ReadAll(privateKeyFile)
	if err != nil {
		log.Panicf("error reading private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Panicf("error parsing private key: %v", err)
	}
	sshConfig := &ssh.ClientConfig{
		User:              sc.user,
		Auth:              []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyAlgorithms: []string{ssh.KeyAlgoED25519},
	}
	sshConfig.HostKeyCallback = checkKnownHosts()
	addr := fmt.Sprintf("%s:%d", sc.host, sc.port)
	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		log.Panicf("error : %v", err)
	}
	sc.Client, err = sftp.NewClient(sshClient)
	if err != nil {
		log.Panicf("error : %v", err)
	}
	return
}

func checkKnownHosts() ssh.HostKeyCallback {
	kh, e := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if e != nil {
		log.Fatalln(e)
	}
	return kh
}

// Put Upload file to sftp server
func (sc *SftpClient) Put(localFile, remoteFile string) (err error) {
	srcFile, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	err = sc.MkdirAll(parent)
	if err != nil {
		log.Printf("Error creating path : %v\tError : %v\n", parent, err)
	}
	err = sc.Chown(parent, 998, 998) //998 = plex on my server, not sure how to make this dynamic
	if err != nil {
		log.Printf("Error setting owner : %v\n", err)
	}
	err = sc.Chmod(parent, 0775)
	if err != nil {
		log.Printf("Error setting mode : %v\n", err)
	}

	dstFile, err := sc.Create(remoteFile)
	if err != nil {
		log.Printf("Error creating file : %v\tError : %v\n", remoteFile, err)
		return
	}
	defer dstFile.Close()
	err = dstFile.Chmod(os.FileMode(0664))
	if err != nil {
		log.Printf("Error setting mode : %v\n", err)
		return
	}
	err = dstFile.Chown(998, 998) //998 = plex on my server, not sure how to make this dynamic
	if err != nil {
		log.Printf("Error setting owner : %v\n", err)
	}
	_, err = io.Copy(dstFile, srcFile)
	return
} // Put Upload file to sftp server

func (sc *SftpClient) Get(localFile, remoteFile string) (err error) {
	srcFile, err := sc.Open(remoteFile)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// Make remote directories recursion
	parent := filepath.Dir(localFile)
	err = sc.MkdirAll(parent)
	if err != nil {
		log.Printf("Error creating path : %v\tError : %v\n", parent, err)
	}
	err = sc.Chown(parent, 998, 998) //998 = plex on my server, not sure how to make this dynamic
	if err != nil {
		log.Printf("Error setting owner : %v\n", err)
	}
	err = sc.Chmod(parent, 0775)
	if err != nil {
		log.Printf("Error setting mode : %v\n", err)
	}

	dstFile, err := os.Create(localFile)
	if err != nil {
		log.Printf("Error creating file : %v\tError : %v\n", remoteFile, err)
		return
	}
	defer dstFile.Close()
	err = dstFile.Chmod(os.FileMode(0664))
	if err != nil {
		log.Printf("Error setting mode : %v\n", err)
		return
	}
	err = dstFile.Chown(998, 998) //998 = plex on my server, not sure how to make this dynamic
	if err != nil {
		log.Printf("Error setting owner : %v\n", err)
	}
	_, err = io.Copy(dstFile, srcFile)
	return
}

// PutReader Upload file to sftp server
func (sc *SftpClient) PutReader(reader io.ReadCloser, remoteFile string) (err error) {
	defer reader.Close()
	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		_ = sc.Mkdir(path)
	}

	dstFile, err := sc.Create(remoteFile)
	if err != nil {
		return
	}
	defer dstFile.Close()
	err = dstFile.Chmod(os.FileMode(0666))
	if err != nil {
		return
	}
	_, err = io.Copy(dstFile, reader)
	return
}
