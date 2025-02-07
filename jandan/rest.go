package jandan

import (
	"encoding/json"
	"io"
	"jandan_girl/models"
	"log"
	"net/http"
)

func GetPosts(startId string) (models.Response, error) {
	// HTTP GET request
	url := "https://api.jandan.net/api/v1/comment/list/108629?start_id=" + startId

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		return models.Response{}, err
	}
	req.Header.Set("User-Agent", "Custom-User-Agent/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to make HTTP request: %v", err)
		return models.Response{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return models.Response{}, err
	}

	// Parse JSON response
	var response models.Response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Failed to parse JSON: %s", body)
		return models.Response{}, err
	}
	return response, nil
}
