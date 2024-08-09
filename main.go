package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/thiagodsantos/gomockserver/pkg/config"
	"github.com/thiagodsantos/gomockserver/pkg/request"
)

func startServer(port string, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		request.HandleRequest(w, r)
	})

	log.Printf("Mock server running on port %s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}
}

func main() {
	err := config.LoadConfig(os.Getenv("CONFIG_FILE"))
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	var wg sync.WaitGroup

	for hosts := range config.GetAllServicesConfig() {
		urlParsed, err := url.Parse(hosts)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return
		}

		wg.Add(1)
		go startServer(urlParsed.Port(), &wg)
	}

	wg.Wait()
}
