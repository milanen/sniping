package checker

import (
	"fmt"
	"log"
	"sniping/helpers"
	"sync"
)

func Discord(usernames []string, prefs helpers.Network) {
	api := helpers.Endpoint.Discord
	ExceptList := &ExceptList{}

	sem := make(chan struct{}, prefs.Threads)
	var wg sync.WaitGroup

	for _, username := range usernames {
		u := username
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func(){ <-sem }()

			fmt.Printf("Checking Discord for username: %s\n", u)

			// create payload
			payload := CreateReqBody(u)

			// create client
			client, err := createCustomClient(helpers.GetProxy(), prefs.Timeout, prefs.Gateway)
			if err != nil {
				fmt.Printf("Request failed: Adding %s to exception list...\n", u)
				HandleException(u, ExceptList)
				return
			}

			// make request
			resp, err := DoRequest(client, "POST", api, payload)
			if err != nil {
				fmt.Printf("Request failed: Adding %s to exception list...\n", u)
				HandleException(u, ExceptList)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				msg, err := GetResponse(resp)
				if err != nil {
					log.Fatal(err)
				}
				if msg["taken"] != true {
					fmt.Printf("Username %s is available on Discord\n", u)
					helpers.SaveChecked(u, true)
					return
				}
			}

			if resp.StatusCode == 429 {
				fmt.Printf("Rate limited! Adding %s to exception list...\n", u)
				HandleException(u, ExceptList)
				return
			}

			fmt.Printf("Username %s is taken on Discord\n", u)
			helpers.SaveChecked(u, false)
		}()
	}

	wg.Wait()

	if len(ExceptList.Users) > 0 { // prevent overflow
		fmt.Printf("\nProcessing %d usernames in the exception list...\n", len(ExceptList.Users))
		ExecExceptionList(ExceptList, prefs)
	}
}