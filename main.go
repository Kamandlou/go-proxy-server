package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

var customTransport = http.DefaultTransport

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a new HTTP server with the handleRequest function as the handler
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	server := http.Server{
		Addr:    port,
		Handler: http.HandlerFunc(handleRequest),
	}

	// Start the server and log any errors
	log.Printf("Starting proxy server on %s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}
}

func handleRequest(w http.ResponseWriter, request *http.Request) {
	// Create a new HTTP request with the same method, URL, and body as the original request
	targetURL := request.URL
	proxyRequest, err := http.NewRequest(request.Method, targetURL.String(), request.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request.", http.StatusInternalServerError)
		return
	}

	// Copy the headers from the original request to the proxy request
	for name, values := range request.Header {
		for _, value := range values {
			proxyRequest.Header.Add(name, value)
		}
	}

	// Send the proxy request using the custom transport
	response, err := customTransport.RoundTrip(proxyRequest)
	if err != nil {
		http.Error(w, "Error sending proxy request.", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Copy the headers from the proxy response to the original response
	for name, values := range response.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set the status code of the original response to the status code of the proxy response
	w.WriteHeader(response.StatusCode)

	// Copy the body of the proxy response to the original response
	if _, err := io.Copy(w, response.Body); err != nil {
		http.Error(w, "Error copy body of the proxy response to the original request.", http.StatusInternalServerError)
		return
	}
}
