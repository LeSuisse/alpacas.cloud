package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Alpacas everywhere\n")
}

func Alpaca(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "A random alpaca\n")
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/alpaca", Alpaca)

	log.Fatal(http.ListenAndServe(":8080", router))
}
