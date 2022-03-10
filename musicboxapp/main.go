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

	return "jukebox-front"
}

func getMetalEndpoint() (string, error) {
	metalEndpoint := os.Getenv("METAL_HOST")
	if metalEndpoint == "" {
		return "", errors.New("METAL_HOST is not set")
	}
	return metalEndpoint, nil
}

func getPopEndpoint() (string, error) {
	popEndpoint := os.Getenv("POP_HOST")
	if popEndpoint == "" {
		return "", errors.New("POP_HOST is not set")
	}
	return popEndpoint, nil
}

type metalHandler struct{}

func (h *metalHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	artists, err := getMetalArtists(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - Unexpected Error"))
		return
	}

	fmt.Fprintf(writer, `{"metal artists":"%s"}`, artists)
}

func getMetalArtists(request *http.Request) (string, error) {
	metalEndpoint, err := getMetalEndpoint()
	if err != nil {
		return "-n/a-", err
	}

	client := xray.Client(&http.Client{})
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s", metalEndpoint), nil)
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

	metalArtists := strings.TrimSpace(string(body))
	if len(metalArtists) < 1 {
		return "-n/a-", errors.New("Empty response from metalArtists")
	}

	return metalArtists, nil
}

type popHandler struct{}

func (h *popHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	artists, err := getPopArtists(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - Unexpected Error"))
		return
	}

	fmt.Fprintf(writer, `{"pop artists":"%s"}`, artists)
}
func getPopArtists(request *http.Request) (string, error) {
	popEndpoint, err := getPopEndpoint()
	if err != nil {
		return "-n/a-", err
	}

	client := xray.Client(&http.Client{})
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s", popEndpoint), nil)
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

	popArtists := strings.TrimSpace(string(body))
	if len(popArtists) < 1 {
		return "-n/a-", errors.New("Empty response from popArtists")
	}

	return popArtists, nil
}

type pingHandler struct{}

func (h *pingHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("ping requested, responding with HTTP 200")
	writer.WriteHeader(http.StatusOK)
}

func main() {
	log.Println("Starting server, listening on port " + getServerPort())

	xray.SetLogger(xraylog.NewDefaultLogger(os.Stderr, xraylog.LogLevelInfo))

	metalEndpoint, err := getMetalEndpoint()
	if err != nil {
		log.Fatalln(err)
	}
	popEndpoint, err := getPopEndpoint()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Using -metal- service at " + metalEndpoint)
	log.Println("Using -pop- service at " + popEndpoint)

	xraySegmentNamer := xray.NewFixedSegmentNamer(getXRAYAppName())

	http.Handle("/metal", xray.Handler(xraySegmentNamer, &metalHandler{}))
	http.Handle("/pop", xray.Handler(xraySegmentNamer, &popHandler{}))
	http.Handle("/ping", xray.Handler(xraySegmentNamer, &pingHandler{}))
	log.Fatal(http.ListenAndServe(":"+getServerPort(), nil))
}
