package main

import (
	"fmt"
	"go_day06/pkg/adapters/db"
	"go_day06/pkg/adapters/http"
	"go_day06/pkg/config"
	"go_day06/pkg/entities/admin"
	"go_day06/pkg/usecases/auth"
	"go_day06/pkg/usecases/postHandler"
	"go_day06/pkg/zip"
	"log"
)

var credentialsFile = "admin_credentials.txt"

func main() {
	err := zip.Unzip("static", "static.zip")
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.New(credentialsFile)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	adm := admin.NewAdmin(cfg)
	storage := db.New(cfg)
	defer db.Close(storage)
	a := auth.New(adm)
	posthandler := postHandler.New(storage)
	server := http.New(a, posthandler)
	_ = http.StartServer(8888, server)
}
