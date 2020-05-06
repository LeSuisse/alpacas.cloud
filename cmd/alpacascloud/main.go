package main

import (
	"bytes"
	"errors"
	"github.com/LeSuisse/alpacas.cloud/pkg/images"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var im images.Images

func Index(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob:; connect-src 'self'; frame-ancestors 'none'; form-action 'none'; base-uri 'none';")

	http.ServeFile(c.Writer, c.Request, "./web/dist/index.html")
}

func OpenAPISpec(c *gin.Context) {
	c.Header("Content-Type", "application/json;charset=utf-8")
	http.ServeFile(c.Writer, c.Request, "./web/dist/openapi.json")
}

type GetAlpacaParameters struct {
	Width  int `form:"width" binding:"min=0"`
	Height int `form:"height" binding:"min=0"`
}

func Alpaca(c *gin.Context) {
	var requestParameters GetAlpacaParameters
	if err := c.BindQuery(&requestParameters); err != nil {
		c.String(http.StatusBadRequest, "Parameters are not valid")
		return
	}

	alpacaImg, imageErr := im.Get(images.ImageOpts{
		MaxWidth:  requestParameters.Width,
		MaxHeight: requestParameters.Height,
	})

	if imageErr != nil {
		log.Println(imageErr)
		var e *images.RequestedSizeTooBigError
		if errors.As(imageErr, &e) {
			c.String(http.StatusNotFound, "Cannot find an alpaca with the requested size")
			return
		}
		c.String(http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	readerImg := bytes.NewReader(alpacaImg)

	c.DataFromReader(http.StatusOK, readerImg.Size(), "image/jpeg", readerImg, nil)
}

type GetAlpacaPlaceHolderParameters struct {
	PlaceholderSize string `uri:"placeholder_size" binding:"required,min=3"`
}

func AlpacaPlaceholder(c *gin.Context) {
	var requestParameters GetAlpacaPlaceHolderParameters
	if err := c.BindUri(&requestParameters); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	rePlaceholderQueryString := regexp.MustCompile(`^([1-9]\d*)x([1-9]\d*)$`)
	params := rePlaceholderQueryString.FindStringSubmatch(requestParameters.PlaceholderSize)
	if len(params) != 3 {
		c.String(http.StatusBadRequest, "Parameters are not valid")
		return
	}

	width, errW := strconv.Atoi(params[1])
	height, errH := strconv.Atoi(params[2])

	if errW != nil || errH != nil {
		c.String(http.StatusBadRequest, "Parameters are not valid")
		return
	}

	alpacaImg, imageErr := im.GetPlaceHolder(images.ImageOpts{
		MaxWidth:  width,
		MaxHeight: height,
	})

	if imageErr != nil {
		log.Println(imageErr)
		var e *images.RequestedSizeTooBigError
		if errors.As(imageErr, &e) {
			c.String(http.StatusNotFound, "Cannot find an alpaca with the requested size")
			return
		}
		c.String(http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	readerImg := bytes.NewReader(alpacaImg)

	c.DataFromReader(http.StatusOK, readerImg.Size(), "image/jpeg", readerImg, nil)
}

func main() {
	var err error
	im, err = images.New(os.Getenv("IMAGES_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	metricsPassword := os.Getenv("METRICS_PASSWORD")

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(SecurityHeaders())
	internalAssets := router.Group("/")
	internalAssets.Use(InternalAssetsHeaders())
	internalAssets.GET("/", Index)
	internalAssets.Static("/assets", "./web/dist/")
	router.HEAD("/openapi.json", OpenAPISpec)
	router.GET("/openapi.json", OpenAPISpec)
	router.GET("/alpaca", Alpaca)
	router.GET("/placeholder/:placeholder_size", AlpacaPlaceholder)
	if metricsPassword != "" {
		metrics := router.Group("/metrics")
		metrics.Use(gin.BasicAuth(gin.Accounts{"metrics": metricsPassword}))
		metrics.Use(InternalAssetsHeaders())
		metrics.GET("", ginprom.PromHandler(promhttp.Handler()))
	}

	log.Fatal(router.Run(":8080"))
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; form-action 'none'; base-uri 'none';")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Cross-Origin-Resource-Policy", "cross-origin")
		c.Header("Cross-Origin-Embedder-Policy", "require-corp")
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
	}
}

func InternalAssetsHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "")
		c.Header("Cross-Origin-Resource-Policy", "same-origin")
	}
}
