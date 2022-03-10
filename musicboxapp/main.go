package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/pkg/errors"
)

const defaultPort = "9000"
const defaultStage = "dev"

func getServerPort() string {
	port := os.Getenv("PORT")
	if port != "" {
		return port
	}

	return defaultPort
}

func getStage() string {
	stage := os.Getenv("STAGE")
	if stage != "" {
		return stage
	}

	return defaultStage
}

func getXRAYAppName() string {
	appName := os.Getenv("XRAY_APP_NAME")
	if appName != "" {
		return appName
	}

	return "musicbox-front"
}

func getFlamencoEndpoint() (string, error) {
	flamecoEndpoint := os.Getenv("FLAMENCO_HOST")
	if flamecoEndpoint == "" {
		return "", errors.New("FLAMENCO_HOST is not set")
	}
	return flamecoEndpoint, nil
}

func getOperaEndpoint() (string, error) {
	operaEndpoint := os.Getenv("OPERA_HOST")
	if operaEndpoint == "" {
		return "", errors.New("OPERA_HOST is not set")
	}
	return operaEndpoint, nil
}

type flamencoHandler struct{}

func (h *flamencoHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	artists, err := getFlamencoArtists(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - Unexpected Error"))
		return
	}

	fmt.Fprintf(writer, `{"flamenco artists":"%s"}`, artists)
}

func getFlamencoArtists(request *http.Request) (string, error) {
	flamencoEndpoint, err := getFlamencoEndpoint()
	if err != nil {
		return "-n/a-", err
	}

	client := xray.Client(&http.Client{})
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s", flamencoEndpoint), nil)
	if err != nil {
		return "-n/a-", err
	}

	resp, err := client.Do(req.WithContext(request.Context()))
	if err != nil {
		return "-n/a-", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "-n/a-", err
	}

	flamencoArtists := strings.TrimSpace(string(body))
	if len(flamencoArtists) < 1 {
		return "-n/a-", errors.New("Empty response from flamencoArtists")
	}

	return flamencoArtists, nil
}

type operaHandler struct{}

func (h *operaHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	artists, err := getOperaArtists(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - Unexpected Error"))
		return
	}

	fmt.Fprintf(writer, `{"Opera artists":"%s"}`, artists)
}
func getOperaArtists(request *http.Request) (string, error) {
	operaEndpoint, err := getOperaEndpoint()
	if err != nil {
		return "-n/a-", err
	}

	client := xray.Client(&http.Client{})
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s", operaEndpoint), nil)
	if err != nil {
		return "-n/a-", err
	}

	resp, err := client.Do(req.WithContext(request.Context()))
	if err != nil {
		return "-n/a-", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "-n/a-", err
	}

	operaArtists := strings.TrimSpace(string(body))
	if len(operaArtists) < 1 {
		return "-n/a-", errors.New("Empty response from operaArtists")
	}

	return operaArtists, nil
}

type pingHandler struct{}

func (h *pingHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("ping requested, responding with HTTP 200")
	writer.WriteHeader(http.StatusOK)
}

func main() {
	log.Println("Starting server, listening on port " + getServerPort())

	xray.SetLogger(xraylog.NewDefaultLogger(os.Stderr, xraylog.LogLevelInfo))

	flamecoEndpoint, err := getFlamencoEndpoint()
	if err != nil {
		log.Fatalln(err)
	}
	operaEndpoint, err := getOperaEndpoint()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Using -flamenco- service at " + flamecoEndpoint)
	log.Println("Using -opera- service at " + operaEndpoint)

	xraySegmentNamer := xray.NewFixedSegmentNamer(getXRAYAppName())

	http.Handle("/flamenco", xray.Handler(xraySegmentNamer, &flamencoHandler{}))
	http.Handle("/opera", xray.Handler(xraySegmentNamer, &operaHandler{}))
	http.Handle("/ping", xray.Handler(xraySegmentNamer, &pingHandler{}))
	log.Fatal(http.ListenAndServe(":"+getServerPort(), nil))
}
