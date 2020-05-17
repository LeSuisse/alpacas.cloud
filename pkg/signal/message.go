package signal

import (
	"github.com/signal-golang/textsecure"
	"io"
	"log"
	"net/http"
	"strings"
)

const FetchImageURL = "https://alpacas.cloud/alpaca?width=800"

func MessageHandler(incomingMessage *textsecure.Message) {
	if !doesMessageTriggersAResponse(incomingMessage.Message()) {
		return
	}

	response, err := http.Get(FetchImageURL)
	if err != nil {
		log.Print(err)
		return
	}
	defer response.Body.Close()

	sendImageResponse(incomingMessage, response.Body)
}

func doesMessageTriggersAResponse(messageContent string) bool {
	var triggeringKeywords = [...]string{"alpaca", "alpaga"}

	lowercaseMessageContent := strings.ToLower(messageContent)
	for _, keyword := range triggeringKeywords {
		if strings.Contains(lowercaseMessageContent, keyword) {
			return true
		}
	}

	return false
}

func sendImageResponse(incomingMessage *textsecure.Message, r io.Reader) {
	incomingMessageGroup := incomingMessage.Group()

	var err error
	if incomingMessageGroup == nil {
		_, err = textsecure.SendAttachment(incomingMessage.Source(), "", r, 0)
	} else {
		_, err = textsecure.SendGroupAttachment(incomingMessageGroup.Hexid, "", r, 0)
	}
	if err != nil {
		log.Println("Could not send message: ", err)
	}
}
