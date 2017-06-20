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
	fmt.Fprintln(w, "Alpacas are everywhere\nThe alpaca you are looking for is at GET /alpaca")
}

func Alpaca(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	f, err := os.Open(im.Get())
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	http.ServeContent(w, r, "", time.Time{}, f)
}

func main() {
	var err error
	im, err = images.New(os.Getenv("IMAGES_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/alpaca", Alpaca)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", &Server{router}))
}

type Server struct {
	router *httprouter.Router
}

func (s *Server) ServeHTTP (w http.ResponseWriter, req *http.Request) {
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none';")
	s.router.ServeHTTP(w, req)
}