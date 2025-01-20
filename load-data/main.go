package main

import (
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"

	"github.com/DanielOsorio01/enron-email-search/load-data/email"
)

const (
	zincURL        = "http://localhost:4080"
	batchSize      = 1000
	workerCount    = 16 // Adjust this based on the number of CPU cores
	fileQueueSize  = 5000
	emailQueueSize = 2000
)

func main() {
	// Check if the user provided at least one argument (the folder name)
	if len(os.Args) < 2 {
		fmt.Println("Please provide the folder name as a positional argument.")
		os.Exit(1)
	}

	// read the folder name from the first argument
	rootFolder := os.Args[1]

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	var memprofile = flag.String("memprofile", "", "write memory profile to file")

	flag.CommandLine.Parse(os.Args[2:]) // Parse only the flags

	// Debugging: Print parsed values
	fmt.Println("Root folder:", rootFolder)
	fmt.Println("cpuprofile flag value:", *cpuprofile)
	fmt.Println("memprofile flag value:", *memprofile)

	if *cpuprofile != "" {
		fmt.Println("Starting CPU profiling...")
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	// Start memory profiling if the flag is set
	if *memprofile != "" {
		fmt.Println("Starting memory profiling...")
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			runtime.GC() // Force garbage collection to capture all data
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal("could not write memory profile:", err)
			}
			f.Close()
			fmt.Println("Memory profiling complete.")
		}()
	}

	var wg sync.WaitGroup // WaitGroup to wait for all goroutines to complete
	var workerWg sync.WaitGroup // WaitGroup to wait for all workers to complete

	var client *email.ZincClient = email.NewZincClient(zincURL, "admin", "Complexpass#123")

	// Create a channel to receive email file paths
	emailFiles := make(chan string, fileQueueSize) // Buffer capacity of 5000 files
	// Create a channel to send parsed emails
	emailQueue := make(chan email.Email, emailQueueSize) // Buffer capacity of 2000 emails

	// Discover all email files in the root folder
	wg.Add(1)
	go email.DiscoverEmailFiles(rootFolder, emailFiles, &wg)
	fmt.Println("Discovering email files...")

	// Start the worker pool to process email files
	wg.Add(1) // Increment the counter for the worker pool
	for i := 0; i < workerCount; i++ {
		workerWg.Add(1)
		go email.ProcessEmailFiles(emailFiles, emailQueue, &workerWg)
		fmt.Printf("Starting worker %d\n", i+1)
	}

	// Goroutine that closes the emailQueue channel
	go func() {
		workerWg.Wait() // Wait for all workers to complete
		close(emailQueue) // Close the emailQueue channel
		wg.Done() // Decrement the counter for the worker pool
	}()

	// Start a batch processor to send emails to Zinc
	wg.Add(1)
	go email.SendEmailBatches(emailQueue, batchSize, &wg, client)
	fmt.Println("Starting email batch processor...")

	// Wait for all goroutines to complete
	wg.Wait()
	fmt.Println("All goroutines completed.")
}
