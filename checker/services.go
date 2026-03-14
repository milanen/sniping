package checker

import (
	"encoding/json"
	"fmt"
	"log"
	"sniping/helpers"
)

func Discord(usernames []string, prefs helpers.Network) {
	api := helpers.Endpoint.Discord

	for _, username := range usernames {
		fmt.Printf("Checking Discord for username: %s\n", username)

		// create payload
		payload := CreateReqBody(username)

		// create client
		client, err := createCustomClient(helpers.GetProxy(), prefs.Timeout, prefs.Gateway)
		if err != nil {
			fmt.Println("Request failed:", err)
			continue
		}

		// make request
		resp, err := DoRequest(client, "POST", api, payload)
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