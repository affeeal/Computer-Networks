package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

var (
	hostname = flag.String("hostname", "smtp.gmail.com", "authentication hostname")
	hostport = flag.String("hostport", "587", "authentication host port")
	username = flag.String("username", "ilya.afanasyev.26@gmail.com", "authentication username")
	password = flag.String("password", "", "authentication password") // пароль указываю вручную
)

func main() {
	flag.Parse()
	auth := smtp.PlainAuth("", *username, *password, *hostname)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		switch scanner.Text() {
		case "SEND":
			var (
				to            []string
				subject, body string
			)
			fmt.Print("To: ")
			scanner.Scan()
			to = append(to, scanner.Text())

			fmt.Print("Subject: ")
			scanner.Scan()
			subject = scanner.Text()

			fmt.Print("Body: ")
			scanner.Scan()
			body = scanner.Text()

			msg := []byte("To: " + to[0] + "\r\n" +
				"Subject: " + subject + "\r\n" +
				"\r\n" +
				body + "\r\n")

			err := smtp.SendMail(*hostname+":"+*hostport, auth, *username, to, msg)
			if err != nil {
				log.Println(err)
			}
			log.Println("The email has been sent successfully!")
		case "QUIT":
			os.Exit(0)
		default:
			fmt.Println("SEND - отправить письмо;\n" +
				"QUIT - выйти из приложения.")
		}
		fmt.Print("> ")
	}
}
