package email

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// IndexEmails processes all files in the directory tree using Depth-First Search.
func LoadEmails(root string) ([]Email, error) {
	var emails []Email
	var mu sync.Mutex
	var wg sync.WaitGroup

	// thread function to parse email files concurrently
	thread := func(path string) {
		defer wg.Done()

		// Parse the email file
		email, err := ParseEmail(path)
		if err != nil {
			fmt.Printf("Warning: Failed to parse file %s: %v\n", path, err)
			return // Continue processing other files
		}

		mu.Lock()
		emails = append(emails, email)
		mu.Unlock()
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %v", path, err)
		}

		// Skip directories; process files only
		if info.IsDir() {
			return nil
		}

		// Process the file in a separate goroutine
		wg.Add(1)
		go thread(path)
		return nil
	})

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return emails, nil
}

func SaveEmails(rootFolder string) (uint64, error) {
	// This function reads all email files in the root folder and
	// saves them to a NDJSON file.

	var err error
	var count uint64
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create a new file to save the emails
	file, err := os.Create("emails.ndjson")
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Loop through all files in the root folder
	err = filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %v", path, err)
		}

		// Skip directories; process files only
		if info.IsDir() {
			return nil
		}

		// Process the file in a separate goroutine
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			// Parse the email file
			email, err := ParseEmail(path)
			if err != nil {
				fmt.Printf("Warning: Failed to parse file %s: %v\n", path, err)
				return // Continue processing other files
			}

			jsonData, err := json.Marshal(email)

			if err != nil {
				fmt.Printf("Warning: Failed to marshal email to JSON: %v\n", err)
				return
			}

			// Save the email to the NDJSON file
			mu.Lock()
			_, err = fmt.Fprintf(file, "%s\n", jsonData)
			mu.Unlock()
			if err != nil {
				fmt.Printf("Warning: Failed to save email to file: %v\n", err)
			} else {
				count++
			}
		}(path)

		return nil
	})

	wg.Wait()

	if err != nil {
		return 0, fmt.Errorf("error saving emails: %v", err)
	}

	return count, err
}

func DiscoverEmailFiles(rootFolder string, filepaths chan<- string, wg *sync.WaitGroup) {
	defer close(filepaths) // Close the channel when this function returns

	err := filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %v", path, err)
		}

		if !info.IsDir() {
			filepaths <- path // Send the file path to the channel
		}

		return nil
	})

	if err != nil {
		log.Printf("Error discovering email files: %v\n", err)
	}

	log.Println("Finished discovering email files.")

	wg.Done() // Decrement the WaitGroup counter
}
