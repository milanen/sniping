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

func createCustomClient(proxy string, timeout int, useTransport bool) (*http.Client, error) {
	// create and return http client
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	if useTransport {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}

		// configure the http.Transport
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}

	return client, nil
}

func Discord(usernames []string, prefs helpers.Network) {
	api := "https://discord.com/api/v9/unique-username/username-attempt-unauthed"

	for _, username := range usernames {
		fmt.Printf("Checking Discord for username: %s\n", username)

		// creating payload
		payload := map[string]string {
			"username": username,
		}

		jsonData, _ := json.Marshal(payload)

		// creating client
		client, err := createCustomClient(helpers.GetProxy(), prefs.Timeout, prefs.Gateway)

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
				helpers.SaveChecked(username, true)
				continue
			}

		fmt.Printf("Username %s is taken on Discord\n", username)
		helpers.SaveChecked(username, false)
		}
	}
}