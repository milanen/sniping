package checker

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"sniping/helpers"
	"sync"
	"time"
)

type ExceptList struct {
    Users []string
    Mu sync.Mutex
}

func GetResponse(resp *http.Response) (map[string]any, error) {
    var response map[string]any
			err := json.NewDecoder(resp.Body).Decode(&response)
			if err != nil {
				log.Fatal(err)
			}
    return response, nil
}

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

func CreateReqBody(username string) map[string]string {
    return map[string]string{
        "username": username,
    }
}

func DoRequest(
    client *http.Client,
    method string,
    url string,
    payload any,
) (*http.Response, error) {

    var body io.Reader

    if payload != nil {
        jsonData, err := json.Marshal(payload)
        if err != nil {
            return nil, err
        }
        body = bytes.NewBuffer(jsonData)
    }

    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }

    // headers
	helpers.SetHeaders(req)

    return client.Do(req)
}

func HandleException(username string, list *ExceptList) {
    list.Mu.Lock()
    list.Users = append(list.Users, username)
    list.Mu.Unlock()
}

func ExecExceptionList(list *ExceptList, prefs helpers.Network) {
    	Discord(list.Users, prefs)
}