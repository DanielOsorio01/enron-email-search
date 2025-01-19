package email

import (
	"fmt"
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

			// Save the email to the NDJSON file
			mu.Lock()
			_, err = fmt.Fprintf(file, "%s\n", email)
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
		return 0, fmt.Errorf("Error saving emails: %v\n", err)
	}

	return count, err
}
