package helpers

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Gateway struct {
	Host string
	Port string
	User string
	Pass string
}

func GetProxy() string {
	data, _ := os.ReadFile("config/proxies.yaml")
	var raw map[string]Gateway
	yaml.Unmarshal(data, &raw)
	url := BuildUrl(raw["proxies"])
	return url
}

func BuildUrl(proxy Gateway) string {
	if proxy.Host == "" || proxy.Port == "" {
		return "";
	}

	url := "http://"
	// structuring user:pass@host:port
	if proxy.User != "" || proxy.Pass != "" {
		url += fmt.Sprintf("%s:%s@", proxy.User, proxy.Pass)
	}
	
	url += fmt.Sprintf("%s:%s", proxy.Host, proxy.Port)
	return url
}