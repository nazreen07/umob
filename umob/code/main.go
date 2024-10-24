package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Function to fetch provider URLs from environment variables
func getProviderURLs() []string {
	return []string{
		os.Getenv("PROVIDER1_URL"),
		os.Getenv("PROVIDER2_URL"),
		os.Getenv("PROVIDER3_URL"),
	}
}

func ingestGBFSData() {
	providerURLs := getProviderURLs()

	totalBikes := 0

	// Fetch and update Prometheus metrics for each provider
	for _, providerURL := range providerURLs {
		freeBikeStatusURL, err := fetchFreeBikeStatusURL(providerURL)
		if err != nil {
			log.Printf("Error fetching free bike status URL from %s: %v", providerURL, err)
			continue
		}

		numBikes, err := fetchFreeBikeStatusData(freeBikeStatusURL)
		if err != nil {
			log.Printf("Error fetching free bike status data from %s: %v", freeBikeStatusURL, err)
			continue
		}

		fmt.Printf("Provider: %s, Available Bikes: %d\n", providerURL, numBikes)

		totalBikes += numBikes
	}

	fmt.Printf("Total Available Bikes: %d\n", totalBikes)
	log.Printf("Ingested data for %d providers. Total bikes available: %d", len(providerURLs), totalBikes)
}

func main() {
	// Start automated ingestion in the background
	go func() {
		for {
			ingestGBFSData()
			time.Sleep(1 * time.Minute)
		}
	}()

	// Create a new Gin router
	router := gin.Default()

	// Expose Prometheus metrics on /metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Run the server on port 8080
	router.Run(":8080")
}
