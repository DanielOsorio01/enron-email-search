package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/DanielOsorio01/enron-email-search/load-data/email"
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
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Read email files
	fmt.Printf("Indexing emails in %s...\n", rootFolder)
	startTime := time.Now()
	emails, err := email.LoadEmails(rootFolder)
	duration := time.Since(startTime)
	if err != nil {
		fmt.Printf("Error indexing emails: %v\n", err)
		return
	}
	fmt.Printf("%d emails indexed in %v.\n", len(emails), duration)

	// log.Printf("Email example: \n%s\n", emails[0].String())

	fmt.Println("Posting emails to the database...")
	startTime = time.Now()
	err = email.PostEmails(emails)
	duration = time.Since(startTime)
	if err != nil {
		fmt.Printf("Error sending emails to bulkv2 API: %v\n", err)
		return
	}
	fmt.Printf("Emails sent to Zincsearch database in %v.\n", duration)

}
