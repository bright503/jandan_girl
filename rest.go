package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func getPosts(startId string) Response {
	// HTTP GET request
	url := "https://api.jandan.net/api/v1/comment/list/108629?start_id=" + startId

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}
	req.Header.Set("User-Agent", "Custom-User-Agent/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Parse JSON response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Failed to parse JSON: %s", body)
	}
	return response
}
