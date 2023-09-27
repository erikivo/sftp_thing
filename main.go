package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	sftpHost     = "HOST"
	sftpPort     = 0
	sftpUsername = "sftp_user"
	sftpPassword = "PASSWORD"
	remotePath   = "REMOTE PATH"
)

func main() {
	// Create a CSV file
	csvData := [][]string{
		{"Name", "Email"},
		{"Javier", "javier@example.com"},
		{"Alex", "alex@example.com"},
	}

	fileName := "data.csv"
	if err := createCSVFile(fileName, csvData); err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}

	// Connect to the SFTP server
	client, err := connectToSFTP()
	if err != nil {
		fmt.Printf("Error connecting to SFTP server: %v\n", err)
		return
	}
	defer client.Close()

	fmt.Println("connected to sftp successfully.")
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

func connectToSFTP() (*sftp.Client, error) {
	config := &ssh.ClientConfig{
		User: sftpUsername,
		Auth: []ssh.AuthMethod{
			ssh.Password(sftpPassword),
		},
		Timeout: 5 * time.Second,
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
