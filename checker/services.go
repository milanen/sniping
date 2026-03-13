package checker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sniping/helpers"
	"time"
)

func createCustomClient(proxy string) (*http.Client, error) {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %w", err)
	}

	// configure the http.Transport
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	// create and return http client
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: transport,
	}

	return client, nil
}

func Discord(usernames []string) {
	api := "https://discord.com/api/v9/unique-username/username-attempt-unauthed"

	for _, username := range usernames {
		fmt.Printf("Checking Discord for username: %s\n", username)

		// creating payload
		payload := map[string]string {
			"username": username,
		}

		jsonData, _ := json.Marshal(payload)

		// creating client
		client, err := createCustomClient(helpers.GetProxy())

		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			return
		}

		req, err := http.NewRequest("POST", api, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Request creation failed:", err)
			continue
		}

		// headers
		helpers.SetHeaders(req)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Request failed:", err)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			var response map[string]any
			err := json.NewDecoder(resp.Body).Decode(&response)
			if err != nil {
				log.Fatal(err)
			}

			if (response["taken"] != true) {
				fmt.Printf("Username %s is available on Discord\n", username)
				continue
			}

		fmt.Printf("Username %s is taken on Discord\n", username)
		}
	}
}