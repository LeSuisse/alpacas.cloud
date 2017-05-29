package main

import (
	"fmt"
	"github.com/LeSuisse/alpacas.cloud/pkg/images"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var im images.Images

func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Alpacas everywhere\n")
}

func Alpaca(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, im.Random())
}

func main() {
	var err error
	im, err = images.LoadImages("/tmp/alpacas")
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/alpaca", Alpaca)

	log.Fatal(http.ListenAndServe(":8080", router))
}
