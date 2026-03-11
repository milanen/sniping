package helpers

import (
	"bufio"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func RandomizeUserAgent() string {
	rand.Seed(time.Now().UnixNano())

	userAgents := []string {
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120 Safari/537.36 Edg/120.0.0.0",
	}
	return (userAgents[rand.Intn(len(userAgents))])
}

func SetHeaders(req *http.Request) {
    req.Header.Set("User-Agent", RandomizeUserAgent())
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
}

func LoadUsernames() []string {
    file, err := os.Open("config/users.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

	var users []string
    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
		users = append(users, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

	return users
}