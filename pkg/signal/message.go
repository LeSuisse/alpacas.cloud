package signal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/signal-golang/textsecure"
	"io"
	"log"
	"net/http"
	"strings"
)

const FetchImageURL = "https://alpacas.cloud/alpaca?width=800"

var (
	receiveMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signal_received_message_total",
		Help: "Total number of received messages",
	})
	responseMessages = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "signal_response_message_total",
			Help: "Total number of responded messages",
		},
		[]string{"status"},
	)
	processingMessageDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "signal_processing_message_duration_seconds",
			Help:    "A message processing duration histogram",
			Buckets: []float64{.1, .25, .5, .75, 1, 2.5, 5},
		},
	)
)

func MessageHandler(incomingMessage *textsecure.Message) {
	timer := prometheus.NewTimer(processingMessageDuration)
	defer timer.ObserveDuration()
	receiveMessages.Inc()
	if !doesMessageTriggersAResponse(incomingMessage.Message()) {
		return
	}

	response, err := http.Get(FetchImageURL)
	if err != nil {
		responseMessages.WithLabelValues("api_call_failure").Inc()
		log.Print(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		responseMessages.WithLabelValues("api_service_failure").Inc()
		log.Printf("API responded with an HTTP status code %d", response.StatusCode)
		return
	}

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
		responseMessages.WithLabelValues("send_failure").Inc()
		log.Println("Could not send message: ", err)
	}
	responseMessages.WithLabelValues("success").Inc()
}
