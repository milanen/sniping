package checker

import (
	"fmt"
	"log"
	"sniping/helpers"
)

func Discord(usernames []string, prefs helpers.Network) {
	api := helpers.Endpoint.Discord
	ExceptList := &ExceptList{}

	for _, username := range usernames {
		fmt.Printf("Checking Discord for username: %s\n", username)

		// create payload
		payload := CreateReqBody(username)

		// create client
		client, err := createCustomClient(helpers.GetProxy(), prefs.Timeout, prefs.Gateway)
		if err != nil {
			fmt.Printf("Request failed: Adding %s to exception list...\n", username)
			HandleException(username, ExceptList)
			continue
		}

		// make request
		resp, err := DoRequest(client, "POST", api, payload)
		if err != nil {
			fmt.Printf("Request failed: Adding %s to exception list...\n", username)
			HandleException(username, ExceptList)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			msg, err := GetResponse(resp)
			if err != nil {
    			log.Fatal(err)
			}

			if (msg["taken"] != true) {
				fmt.Printf("Username %s is available on Discord\n", username)
				helpers.SaveChecked(username, true)
				continue
			}
		}

		if resp.StatusCode == 429 {
			fmt.Printf("Rate limited! Adding %s to exception list...\n", username)
			HandleException(username, ExceptList)
			continue
		}

		fmt.Printf("Username %s is taken on Discord\n", username)
		helpers.SaveChecked(username, false)
	}
	if len(ExceptList.Users) > 0 { // prevent overflow
		fmt.Printf("\nProcessing %d usernames in the exception list...\n", len(ExceptList.Users))
    	ExecExceptionList(ExceptList, prefs)
	}
}
