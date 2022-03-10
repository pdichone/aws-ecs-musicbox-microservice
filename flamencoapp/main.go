package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
)

const defaultPort = "9001"
const defaultArtists = "Paco de Lucia, Estrella Morente"
const defaultStage = "dev"

func getServerPort() string {
	port := os.Getenv("PORT")
	if port != "" {
		return port
	}

	return defaultPort
}

func getArtists() string {
	artists := os.Getenv("ARTISTS")
	if artists != "" {
		return artists
	}

	return defaultArtists
}

func getXRAYAppName() string {
	appName := os.Getenv("XRAY_APP_NAME")
	if appName != "" {
		return appName
	}

	return "flamenco"
}

func getStage() string {
	stage := os.Getenv("STAGE")
	if stage != "" {
		return stage
	}

	return defaultStage
}

type flamencoHandler struct{}

func (h *flamencoHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("flamenco artists requested, responding with", getArtists())
	fmt.Fprint(writer, getArtists())
}

type pingHandler struct{}

func (h *pingHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("ping to flamenco-svc requested, responding with HTTP 200")
	writer.WriteHeader(http.StatusOK)
}

func main() {
	log.Println("starting server, listening on port " + getServerPort())
	xraySegmentNamer := xray.NewFixedSegmentNamer(getXRAYAppName())
	http.Handle("/", xray.Handler(xraySegmentNamer, &flamencoHandler{}))
	http.Handle("/ping", xray.Handler(xraySegmentNamer, &pingHandler{}))
	http.ListenAndServe(":"+getServerPort(), nil)
}
