package requestsender

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	requestCount   int
	requestCountMu sync.Mutex
	started        bool
	startedMu      sync.Mutex
)

func StartRequester() error {
	startedMu.Lock()
	defer startedMu.Unlock()

	if started {
		return nil
	}

	targetURL, rate, err := loadEnv()
	if err != nil {
		return err
	}

	if rate <= 0 {
		return fmt.Errorf("invalid REQUESTS_PER_SECOND value: %d", rate)
	}

	started = true
	go runRequester(targetURL, rate)
	return nil
}

func RequestsSent() int {
	requestCountMu.Lock()
	defer requestCountMu.Unlock()
	return requestCount
}

func loadEnv() (string, int, error) {
	file, err := os.Open(".env")
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var targetURL string
	var rate int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "TARGET_URL":
			targetURL = value
		case "REQUESTS_PER_SECOND":
			rate, err = strconv.Atoi(value)
			if err != nil {
				return "", 0, fmt.Errorf("invalid REQUESTS_PER_SECOND value: %w", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", 0, err
	}

	if targetURL == "" {
		return "", 0, fmt.Errorf("TARGET_URL is required in .env")
	}

	if rate <= 0 {
		return "", 0, fmt.Errorf("REQUESTS_PER_SECOND must be greater than 0")
	}

	return targetURL, rate, nil
}

func runRequester(targetURL string, rate int) {
	interval := time.Second / time.Duration(rate)
	if interval <= 0 {
		interval = time.Millisecond
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	client := &http.Client{Timeout: 10 * time.Second}

	for range ticker.C {
		resp, err := client.Get(targetURL)
		if err != nil {
			log.Printf("request failed: %v", err)
			continue
		}
		resp.Body.Close()

		requestCountMu.Lock()
		requestCount++
		requestCountMu.Unlock()
	}
}
