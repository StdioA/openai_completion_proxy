package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var (
	endpointBase       = flag.String("endpoint", "https://api.openai.com/v1", "API Endpoint Base")
	openAIChatEndpoint string
	listen             = flag.String("listen", ":8080", "Listen address")
)

func main() {
	flag.Parse()

	openAIChatEndpoint = strings.TrimSuffix(*endpointBase, "/") + "/chat/completions"
	r := mux.NewRouter()
	r.HandleFunc("/v1/completions", handleCompletions).Methods("POST")

	log.Printf("Starting server on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, r))
}

func handleCompletions(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Unmarshal the legacy completion request
	var legacyReq map[string]interface{}
	if err := json.Unmarshal(body, &legacyReq); err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	// Translate to chat completion request
	chatReq := map[string]interface{}{
		"model": legacyReq["model"],
		"messages": []map[string]string{
			{"role": "user", "content": legacyReq["prompt"].(string)},
		},
		"max_tokens":        legacyReq["max_tokens"],
		"temperature":       legacyReq["temperature"],
		"top_p":             legacyReq["top_p"],
		"n":                 legacyReq["n"],
		"stream":            legacyReq["stream"],
		"stop":              legacyReq["stop"],
		"presence_penalty":  legacyReq["presence_penalty"],
		"frequency_penalty": legacyReq["frequency_penalty"],
	}

	// Marshal the chat completion request
	chatReqBody, err := json.Marshal(chatReq)
	if err != nil {
		http.Error(w, "Failed to marshal chat completion request", http.StatusInternalServerError)
		return
	}

	// Forward the request to OpenAI's chat completion API
	req, err := http.NewRequest("POST", openAIChatEndpoint, bytes.NewBuffer(chatReqBody))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if len(r.Header["Authorization"]) > 0 {
		req.Header.Set("Authorization", r.Header["Authorization"][0])
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to OpenAI", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// Write the response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}
