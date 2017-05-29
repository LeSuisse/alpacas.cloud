package main

import (
	"fmt"
	"github.com/LeSuisse/alpacas.cloud/pkg/images"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"time"
)

var im images.Images

func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Alpacas everywhere\n")
}

func Alpaca(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	f, err := os.Open(im.Random())
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	http.ServeContent(w, r, "", time.Time{}, f)
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
