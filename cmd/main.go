package main

import (
	"L0"
	"L0/nats"
	"L0/pkg/handler"
	"L0/pkg/repository"
	"L0/pkg/service"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "fivvif",
		Password: "6187",
		DBName:   "vibe-db",
		SSLMode:  "disable",
	})
	repos := repository.NewRepository(db)
	service, err := service.NewService(repos)
	if err != nil {
		logrus.Errorf("Error while creating service : %s", err.Error())
	}
	handlers := handler.NewHandler(service)
	srv := new(L0.Server)
	go func() {
		if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
			log.Fatalf("error while running server : %s ", err.Error())
		}
	}()
	logrus.Println("Server started")

	clusterId := "test-cluster"
	clientId := "vibe"
	channelName := "vibeChannel"
	nats, err := nats.NewSubscribeToChannel(clusterId, clientId, channelName, repos, service)
	if err != nil {
		logrus.Fatalln(err.Error())
	}

	// Прокид json  в nats
	go func() {
		for {
			var filename string
			fmt.Scanln(&filename)
			file := fmt.Sprintf("json/%s", filename)
			jsonStr, err := ioutil.ReadFile(file)
			if err != nil {
				logrus.Errorf("Error while reading json : %s", err.Error())
			}

			if err := nats.Publish(channelName, jsonStr); err != nil {
				logrus.Errorf("Error while publishing message : %s", err.Error())
			}
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Println("Server shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Error("error while shutting down server", err.Error())
	}

	if err := nats.Close(); err != nil {
		logrus.Error("error while closing channel", err.Error())
	}
}
