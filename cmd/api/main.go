package main

import (
	"flag"
	"github.com/simonhylander/gorsk/pkg/api"
	"github.com/simonhylander/gorsk/pkg/utl/config"
	"log"
	"os"
)

func main() {
	cfgPath := flag.String("p", "./cmd/api/conf.local.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)

	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	err = api.Start(cfg)

	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}