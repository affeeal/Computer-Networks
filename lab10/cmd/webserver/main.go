package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/affeeal/dashboard/internal/app/webserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/webserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := webserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := webserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
