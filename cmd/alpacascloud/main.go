package main

import (
	"bytes"
	"errors"
	"github.com/LeSuisse/alpacas.cloud/pkg/images"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var im images.Images

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob:; connect-src 'self'; frame-ancestors 'none'; form-action 'none'; base-uri 'none';")
	http.ServeFile(w, r, "./web/dist/index.html")
}

func OpenAPISpec(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	http.ServeFile(w, r, "./web/dist/openapi.json")
}

func Alpaca(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var opts images.ImageOpts

	widthQuery := r.URL.Query().Get("width")
	if widthQuery != "" {
		width, err := strconv.Atoi(widthQuery)
		if err != nil || width < 1 {
			http.Error(w, "400 Bad Request - Invalid width parameter", http.StatusBadRequest)
			return
		}
		opts.MaxWidth = width
	}
	heightQuery := r.URL.Query().Get("height")
	if heightQuery != "" {
		height, err := strconv.Atoi(heightQuery)
		if err != nil || height < 1 {
			http.Error(w, "400 Bad Request - Invalid height parameter", http.StatusBadRequest)
			return
		}
		opts.MaxHeight = height
	}

	alpacaImg, imageErr := im.Get(opts)

	if imageErr != nil {
		log.Println(imageErr)
		var e *images.RequestedSizeTooBigError
		if errors.As(imageErr, &e) {
			http.Error(w, "Cannot find an alpaca with the requested size", http.StatusNotFound)
			return
		}
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	readerImg := bytes.NewReader(alpacaImg)

	http.ServeContent(w, r, "", time.Time{}, readerImg)
}

func main() {
	var err error
	im, err = images.New(os.Getenv("IMAGES_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/openapi.json", OpenAPISpec)
	router.GET("/alpaca", Alpaca)
	router.ServeFiles("/assets/*filepath", http.Dir("./web/dist/"))

	log.Fatal(http.ListenAndServe(":8080", &Server{router}))
}

type Server struct {
	router *httprouter.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; form-action 'none'; base-uri 'none';")
	s.router.ServeHTTP(w, req)
}
