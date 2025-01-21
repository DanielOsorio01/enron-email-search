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
	"sync/atomic"
	"time"

	"github.com/DanielOsorio01/enron-email-search/load-data/email"
)

const (
	zincURL        = "http://localhost:4080"
	batchSize      = 5000
	fileQueueSize  = 530000
	emailQueueSize = 15000
	statusInterval = 5 // seconds
)

func main() {
	// Check if the user provided at least one argument (the folder name)
	if len(os.Args) < 2 {
		fmt.Println("Please provide the folder name as a positional argument.")
		os.Exit(1)
	}

	// read the folder name from the first argument
	rootFolder := os.Args[1]

	var workerCount = runtime.NumCPU()
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	var memprofile = flag.String("memprofile", "", "write memory profile to file")

	flag.CommandLine.Parse(os.Args[2:]) // Parse only the flags

	// Debugging: Print parsed values
	log.Println("Root folder:", rootFolder)
	log.Println("cpuprofile flag value:", *cpuprofile)
	log.Println("memprofile flag value:", *memprofile)
	log.Println("workerCount:", workerCount)

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

	var wg sync.WaitGroup       // WaitGroup to wait for all goroutines to complete
	var workerWg sync.WaitGroup // WaitGroup to wait for all workers to complete
	var client *email.ZincClient = email.NewZincClient(zincURL, "admin", "Complexpass#123")
	var sentEmails atomic.Uint64   // Counter for the number of emails sent
	var done = make(chan struct{}) // Signal to stop the status goroutine

	// Start a goroutine to print the status
	go func() {
		ticker := time.NewTicker(statusInterval * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Printf("[STATUS] Sent emails: %d\n", sentEmails.Load())
			case <-done:
				log.Println("[STATUS] Final count:", sentEmails.Load())
				return
			}
		}
	}()

	// Create a channel to receive email file paths
	emailFiles := make(chan string, fileQueueSize) // Buffer capacity of 5000 files
	// Create a channel to send parsed emails
	emailQueue := make(chan email.Email, emailQueueSize) // Buffer capacity of 2000 emails

	// Discover all email files in the root folder
	wg.Add(1)
	go email.DiscoverEmailFiles(rootFolder, emailFiles, &wg)
	log.Println("Discovering email files...")

	// Start the worker pool to process email files
	wg.Add(1) // Increment the counter for the worker pool
	for i := 0; i < workerCount; i++ {
		workerWg.Add(1)
		go email.ProcessEmailFiles(emailFiles, emailQueue, &workerWg)
	}

	// Goroutine that closes the emailQueue channel
	go func() {
		workerWg.Wait()   // Wait for all workers to complete
		close(emailQueue) // Close the emailQueue channel
		wg.Done()         // Decrement the counter for the worker pool
	}()

	// Start a batch processor to send emails to Zinc
	wg.Add(1)
	go email.SendEmailBatches(emailQueue, batchSize, &wg, client, &sentEmails)
	log.Println("Starting email batch processor...")

	// Wait for all goroutines to complete
	wg.Wait()
	done <- struct{}{} // Signal the status goroutine to stop
	time.Sleep(1 * time.Second)
}
