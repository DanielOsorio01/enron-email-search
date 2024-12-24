package email

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// BulkV2Data structure for the bulkv2 API request
type BulkV2Data struct {
	Index   string  `json:"index"`
	Records []Email `json:"records"`
}

func postEmails(emails []Email) error {
	var wg sync.WaitGroup
	concurrencyLimit := 2

	sem := make(chan struct{}, concurrencyLimit)

	// Send multiple requests in parallel
	thread := func(index string, records []Email) {
		defer wg.Done()

		// Create a buffer to temporarily store the JSON payload
		var buf bytes.Buffer
		data := BulkV2Data{
			Index:   index,
			Records: records,
		}

		// Directly encode JSON without GZIP compression
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(data); err != nil {
			log.Printf("Error encoding JSON: %v", err)
			return
		}

		// Create a POST request with the plain JSON data
		req, err := http.NewRequest("POST", "http://localhost:4080/api/_bulkv2", &buf)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return
		}

		req.SetBasicAuth("admin", "Complexpass#123")
		req.Header.Set("Content-Type", "application/json")

		// Acquire a semaphore slot before making the request
		sem <- struct{}{}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error sending request: %v", err)
			<-sem // Release the semaphore slot
			return
		}
		defer resp.Body.Close()

		// Print the response status only on error
		if resp.StatusCode != http.StatusOK {
			log.Printf("Request failed with status: %s", resp.Status)
		}

		// Release the semaphore slot after the request is completed
		<-sem
	}

	// Split the emails into chunks and process in parallel
	chunkSize := 50 // Try smaller chunks
	for i := 0; i < len(emails); i += chunkSize {
		end := i + chunkSize
		if end > len(emails) {
			end = len(emails)
		}
		wg.Add(1)
		go thread("enron_emails", emails[i:end])
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}
