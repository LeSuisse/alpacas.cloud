package main

import (
	"log"
	"os"

	"github.com/LeSuisse/alpacas.cloud/pkg/signal"
	"github.com/signal-golang/textsecure"
)

func getConfig() (*textsecure.Config, error) {
	telNumber := os.Getenv("TEL_NUMBER")
	if telNumber == "" {
		log.Fatal("TEL_NUMBER environment variable must be set")
	}
	storageDirectory := os.Getenv("STORAGE_DIRECTORY")
	if storageDirectory == "" {
		log.Fatal("STORAGE_DIRECTORY environment variable must be set")
	}

	return &textsecure.Config{
		Tel:                telNumber,
		StorageDir:         storageDirectory,
		UnencryptedStorage: true,
	}, nil
}

func getLocalContacts() ([]textsecure.Contact, error) {
	return []textsecure.Contact{}, nil
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

	err = textsecure.StartListening()
	if err != nil {
		log.Fatal(err)
	}
}
