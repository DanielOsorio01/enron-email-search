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

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	// Start memory profiling if the flag is set
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		runtime.GC() // Force a garbage collection to get accurate stats
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Read email files
	startTime := time.Now()
	emails, err := email.LoadEmails(rootFolder)
	duration := time.Since(startTime)
	if err != nil {
		fmt.Printf("Error indexing emails: %v\n", err)
		return
	}
	fmt.Printf("%d emails indexed in %v.\n", len(emails), duration)

	// log.Printf("Email example: \n%s\n", emails[0].String())

	// Measure the time taken to send the emails to the bulkv2 API
	startTime = time.Now()
	err = email.PostEmails(emails)
	duration = time.Since(startTime)
	if err != nil {
		fmt.Printf("Error sending emails to bulkv2 API: %v\n", err)
		return
	}
	fmt.Printf("Emails sent to bulkv2 API in %v.\n", duration)

}
