package main

import (
	"fmt"
	handler2 "github.com/mainak90/SimpleProx/handlers"
	log "github.com/sirupsen/logrus"
	"github.com/mainak90/SimpleProx/config"
	"net/http"
)

func main() {
	server()
}

func server(){
	conf := config.New()
	log.WithFields(log.Fields{"Controller": "server"}).Info(fmt.Printf("%+v\n", conf))
	handler := handler2.NewHandler(conf)
	log.WithFields(log.Fields{"Controller": "server"}).Info("Starting server at port: ", conf.ListenHost)
	err := http.ListenAndServe(conf.ListenHost, handler)
	if err != nil {
		fmt.Println("startup failed:", err)
	}
}