package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/gomail.v2"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fileName := fmt.Sprintf("%s/%s", cwd, os.Args[1])
	if fileName == "" {
		fmt.Println("Please provide a file")
		os.Exit(1)
	}
	supportedFormats := []string{".mobi", ".pdf", ".epub", ".azw3", ".txt", ".html"}
	if !slices.Contains(supportedFormats, filepath.Ext(fileName)) {
		fmt.Println("Unsupported file format.\n Supported formats: ", supportedFormats)
		os.Exit(1)
	}
	sender := "albertopluecker@gmail.com"
	receiver := "kFrjSBVr5BNtbd@kindle.com"
	password := os.Getenv("STK_PWD")
	if password == "" {
		log.Fatal("STK_PWD not set")
	}
	host := "smtp.gmail.com"

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "Book sent to Kindle: "+os.Args[1])
	m.Attach(fileName)

	d := gomail.NewDialer(host, 587, sender, password)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
