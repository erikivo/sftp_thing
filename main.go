package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	sftpHost     = ""
	sftpPort     = 22
	sftpUsername = ""
	sftpPassword = ""
	remotePath   = ""
)

func main() {
	// Create a text file
	fileContent := "Name\tEmail\nJavier\tjavier@example.com\nAlex\talex@example.com"

	fileName := "test.txt"
	if err := createTextFile(fileName, fileContent); err != nil {
		fmt.Printf("Error creating text file: %v\n", err)
		return
	}

	client, err := connectToSFTP()
	if err != nil {
		fmt.Printf("Error connecting to SFTP server: %v\n", err)
		return
	}
	defer client.Close()

	err = uploadToSFTP(client, fileName)
	if err != nil {
		fmt.Printf("Error uploading CSV file: %v\n", err)
		return
	}

	fmt.Println("CSV file uploaded successfully.")
}

func connectToSFTP() (*sftp.Client, error) {
	config := &ssh.ClientConfig{
		User: sftpUsername,
		Auth: []ssh.AuthMethod{
			ssh.Password(sftpPassword),
		},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sftpHost, sftpPort), config)
	if err != nil {
		return nil, err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func uploadToSFTP(client *sftp.Client, localPath string) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer localFile.Close()

	remoteFilePath := filepath.Join(remotePath, filepath.Base(localPath))
	fmt.Println("Remote file path:", remoteFilePath)

	remoteFile, err := client.OpenFile(remoteFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return fmt.Errorf("client open file: %w", err)
	}
	defer remoteFile.Close()

	if _, err := localFile.Seek(0, 0); err != nil {
		return fmt.Errorf("localfile: %w", err)
	}

	if _, err := remoteFile.ReadFrom(localFile); err != nil {
		return fmt.Errorf("remoteFile: %w", err)
	}

	return nil
}

func createCSVFile(fileName string, data [][]string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func createTextFile(fileName, content string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
