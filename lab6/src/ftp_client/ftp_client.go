package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

const (
	ADDRESS = "students.yss.su:21"
	PASSWORD = "3Ru7yOTA"
	USER = "ftpiu8"

	PATH = "../../media/"
)

const COMMANDS = `STOR <file-name> - загрузить файл из каталога examples на FTP-сервер;
RETR <file-name> - сохранить файл в каталог downloads с FTP-сервера;
MKD <dir-name>   - создать каталог на FTP-сервере;
DELE <file-name> - удалить файл с FTP-сервера;
NLST             - показать содержимое текущего каталога FTP-сервера;
QUIT             - закрыть соединение с FTP-сервером.`

func stor(conn *ftp.ServerConn, name string) {
	file, err := os.Open(PATH + name)
	if err != nil {
		log.Println("Failed to open file: ", err)
		return
	}
	defer file.Close()

	if err = conn.Stor(name, file); err != nil {
		log.Fatal("Failed to issue STOR FTP command: ", err)
	}
}

func retr(conn *ftp.ServerConn, name string) {
	resp, err := conn.Retr(name)
	if err != nil {
		log.Println("Failed to issue RETR FTP command: ", err)
		return
	}
	defer resp.Close()

	file, err := os.Create(PATH + name)
	if err != nil {
		log.Println("Failed to create file: ", err)
		return
	}
	defer file.Close()

	if _, err = io.Copy(file, resp); err != nil {
		log.Println("Failed to copy respond: ", err)
	}
}

func mkd(conn *ftp.ServerConn, name string) {
	if err := conn.MakeDir(name); err != nil {
		log.Println("Failed to issue MKD FTP command: ", err)
	}
}

func dele(conn *ftp.ServerConn, name string) {
	if err := conn.Delete(name); err != nil {
		log.Println("Failed to issue DELE FTP command: ", err)
	}
}

func nlst(conn *ftp.ServerConn) {
	path, err := conn.CurrentDir()
	if err != nil {
		log.Println("Failed to issue PWD FTP command: ", err)
		return
	}

	entries, err := conn.NameList(path)
	if err != nil {
		log.Println("Failed to issue NLST FTP command: ", err)
		return
	}

	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func quit(conn *ftp.ServerConn) {
	if err := conn.Quit(); err != nil {
	    log.Fatal("Failed to quit: ", err)
	}
}

func main() {
	log.Printf("Connecting to %s FTP-server...\n", ADDRESS)
	conn, err := ftp.Dial(ADDRESS, ftp.DialWithTimeout(5 * time.Second))
	if err != nil {
	    log.Fatal("Failed to dial: ", err)
	}

	err = conn.Login(USER, PASSWORD)
	if err != nil {
	    log.Fatal("Failed to login: ", err)
	}

	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		split := strings.Split(cmd, " ")
		switch split[0] {
		case "STOR": // загрузить файл на сервер
			if len(split) != 2 {
				log.Println("Not enough arguments! Format: STOR <file-name>")
			} else {
				stor(conn, split[1])
			}
		case "RETR": // скачать файл с сервера
			if len(split) != 2 {
				log.Println("Not enough arguments! Format: RETR <file-name>")
			} else {
				retr(conn, split[1])
			}
		case "MKD": // создать новый каталог
			if len(split) != 2 {
				log.Println("Not enough arguments! Format: MKD <dir-name>")
			} else {
				mkd(conn, split[1])
			}
		case "DELE": // удалить файл с сервера
			if len(split) != 2 {
				log.Println("Not enough arguments! Format: DELE <file-name>")
			} else {
				dele(conn, split[1])
			}
		case "NLST": // вывести содержимое каталога
			nlst(conn)
		case "QUIT": // закрыть соединение
			quit(conn)
			return
		default: // вывести доступные команды
			fmt.Println(COMMANDS)
		}
		fmt.Print("> ")
	}
}
