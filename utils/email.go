package utils

import (
	"crypto/tls"
	"ecommerce-golang/config"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	fmt.Println("Parsing templates", paths)

	if err != nil {
		return nil, err
	}
	return template.ParseFiles(paths...)
}

func SendEmail() error {
	config, _ := config.LoadConfig(".")

	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := "bochuangjie@gmail.com"
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	m := gomail.NewMessage()
	message := []byte("This is a test")

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetBody("text/html", string(message))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	fmt.Println("message sent")
	return nil
}
