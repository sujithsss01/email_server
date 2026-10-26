package main

import (
	"fmt"
	"net/smtp"
)

func main() {
	client, err := smtp.Dial("127.0.0.1:1025")
	if err != nil {
		fmt.Println("Connection failed")
		return
	}
	defer client.Close()
	fmt.Println("Connection successful")

	if err := client.Mail("sender@example.com"); err != nil {
		fmt.Println("Failed to set sender")
		return
	}
	if err := client.Rcpt("recipient@example.com"); err != nil {
		fmt.Println("Failed to set reciever")
		return
	}
	wc, err := client.Data()
	if err != nil {
		fmt.Println("Failed to open body")
		return
	}
	fmt.Fprintf(wc, "From: sender@example.com\r\n")
	fmt.Fprintf(wc, "To: recipient@example.com\r\n")
	fmt.Fprintf(wc, "Subject: Test Email\r\n")
	fmt.Fprintf(wc, "\r\n")
	fmt.Fprintf(wc, "Fuck you\r\n")
	err = wc.Close()
	if err != nil {
		fmt.Println("Failed to close body")
	}

}
