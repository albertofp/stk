package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/gomail.v2"

	"github.com/albertofp/stk/utils"
)

type Config struct {
	sender   string
	receiver string
	fileName string
	password string
	host     string
}

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
	password := os.Getenv("STK_PWD")
	if password == "" {
		log.Fatal("STK_PWD not set")
	}

	cfg := Config{
		sender:   "albertopluecker@gmail.com",
		receiver: "kFrjSBVr5BNtbd@kindle.com",
		fileName: fileName,
		password: password,
		host:     "smtp.gmail.com",
	}
	err = utils.WithSpinner(func() error {
		return sendMail(cfg)
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Book sent!")
}

func sendMail(cfg Config) error {
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.sender)
	m.SetHeader("To", cfg.receiver)
	m.SetHeader("Subject", "Book sent to Kindle: "+os.Args[1])
	m.Attach(cfg.fileName)

	d := gomail.NewDialer(cfg.host, 587, cfg.sender, cfg.password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
