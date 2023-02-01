// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"

	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
)

const (
	ROOT = "/home/affeeal/Desktop/programming/III-term/CN_autonomus/lab-06/files"
	NAME = "affeeal"
	PASS = "123456"
	PORT = 2121
	HOST = "localhost"
)

func main() {
	var (
		root = flag.String("root", ROOT, "Root directory to serve")
		user = flag.String("user", NAME, "Username for login")
		pass = flag.String("pass", PASS, "Password for login")
		port = flag.Int("port", PORT, "Port")
		host = flag.String("host", HOST, "Host")
	)
	flag.Parse()

	factory := &filedriver.FileDriverFactory{
		RootPath: *root,
		Perm:     server.NewSimplePerm("user", "group"),
	}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     *port,
		Hostname: *host,
		Auth:     &server.SimpleAuth{Name: *user, Password: *pass},
	}

	log.Printf("Starting ftp server on %v:%v", opts.Hostname, opts.Port)
	log.Printf("Username %v, Password %v", *user, *pass)
	server := server.NewServer(opts)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
