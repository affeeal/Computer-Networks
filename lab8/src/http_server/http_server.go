package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/julienschmidt/httprouter"
)

type ServerConnWrapper struct {
	conn *ftp.ServerConn
}

func (conn ServerConnWrapper) fileHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := ps.ByName("path")
	ftpResp, err := conn.conn.Retr(path)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	defer ftpResp.Close()

	switch filepath.Ext(path) {
	case ".cpp":
		_, codeFileName := filepath.Split(path)
		objectFileName := codeFileName[:len(codeFileName)-4] + "exe"

		codeFile, err := os.Create(codeFileName)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Fatal(err)
		}
		defer os.Remove(codeFileName)

		if _, err = io.Copy(codeFile, ftpResp); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Fatal(err)
		}

		compileCmd := exec.Command("g++", "-o", objectFileName, codeFileName)
		if err = compileCmd.Run(); err != nil {
			log.Fatal(err)
		}

		args := make([]string, 0)
		if r.Method == http.MethodPost {
			r.ParseForm()
			for _, value := range r.Form {
				for _, arg := range value {
					args = append(args, arg)
				}
			}
		} else if r.Method == http.MethodGet {
			params := r.URL.Query()
			for _, value := range params {
				for _, arg := range value {
					args = append(args, arg)
				}
			}
		}

		runCmd := exec.Command("./"+objectFileName, args...)
		runCmd.Stdout = w
		if err = runCmd.Run(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Fatal(err)
		}

		err = os.Remove(objectFileName)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Fatal(err)
		}
	default:
		_, err = io.Copy(w, ftpResp)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Fatal(err)
		}
	}
}

func main() {
	log.Println("Dialing FTP-server...")
	conn, err := ftp.Dial("localhost:2121", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	connWrapper := ServerConnWrapper{conn: conn}
	log.Println("Login FTP-server...")
	err = conn.Login("admin", "123456")
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.GET("/files/*path", connWrapper.fileHandler)
	router.POST("/files/*path", connWrapper.fileHandler)
	log.Println("HTTP-server is listening on 8080...")
	err = http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
