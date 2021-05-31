package main

import (
	"github.com/signal-golang/textsecure/config"
	"github.com/signal-golang/textsecure/contacts"
	"log"
	"net/http"
	"os"

	"github.com/LeSuisse/alpacas.cloud/pkg/signal"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/signal-golang/textsecure"
)

func getConfig() (*config.Config, error) {
	telNumber := os.Getenv("TEL_NUMBER")
	if telNumber == "" {
		log.Fatal("TEL_NUMBER environment variable must be set")
	}
	storageDirectory := os.Getenv("STORAGE_DIRECTORY")
	if storageDirectory == "" {
		log.Fatal("STORAGE_DIRECTORY environment variable must be set")
	}

	return &config.Config{
		Tel:                telNumber,
		StorageDir:         storageDirectory,
		UnencryptedStorage: true,
	}, nil
}

func getLocalContacts() ([]contacts.Contact, error) {
	return []contacts.Contact{}, nil
}

func main() {
	client := &textsecure.Client{
		GetConfig:        getConfig,
		GetLocalContacts: getLocalContacts,
		GetVerificationCode: func() string {
			log.Fatal("Phone number is expected to be already verified, please register it first if needed")
			return ""
		},
		MessageHandler: signal.MessageHandler,
		RegistrationDone: func() {
		},
	}
	err := textsecure.Setup(client)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan bool)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
		stop <- true
	}()
	go func() {
		err = textsecure.StartListening()
		if err != nil {
			log.Fatal(err)
		}
		stop<-true
	}()

	<-stop
}
