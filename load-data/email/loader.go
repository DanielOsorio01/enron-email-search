package email

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

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
