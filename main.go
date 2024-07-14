package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/wneessen/go-mail"

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
	fileName := filepath.Join(cwd, os.Args[1])
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
		log.Fatal("Error: ", err)
	}
	fmt.Println("Book sent!")
}

func sendMail(cfg Config) error {
	m := mail.NewMsg()
	_ = m.From(cfg.sender)
	_ = m.To(cfg.receiver)
	m.Subject(fmt.Sprintf("Book sent to Kindle: %s", filepath.Base(cfg.fileName)))
	if err := checkFileSize(cfg.fileName); err != nil {
		return err
	}
	m.AttachFile(cfg.fileName)

	d, err := mail.NewClient(cfg.host, mail.WithPort(25), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(cfg.sender), mail.WithPassword(cfg.password))
	if err != nil {
		return err
	}

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func checkFileSize(fileName string) error {
	file, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if file.Size() > 25*1024*1024 {
		return fmt.Errorf("file size has to be smaller than 25Mb")
	}
	return nil
}
