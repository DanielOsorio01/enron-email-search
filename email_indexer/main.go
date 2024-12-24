package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Email structure

type Email struct {
	MessageID               string    `json:"message_id"`
	Date                    time.Time `json:"date"`
	From                    string    `json:"from"`
	To                      []string  `json:"to"`
	Cc                      []string  `json:"cc"`
	Bcc                     []string  `json:"bcc"`
	Subject                 string    `json:"subject,omitempty"`
	MimeVersion             string    `json:"mime_version,omitempty"`
	ContentType             string    `json:"content_type,omitempty"`
	ContentTransferEncoding string    `json:"content_transfer_encoding,omitempty"`
	XFrom                   string    `json:"x_from,omitempty"`
	XTo                     []string  `json:"x_to,omitempty"`
	XCc                     []string  `json:"x_cc,omitempty"`
	XBcc                    []string  `json:"x_bcc,omitempty"`
	XFolder                 string    `json:"x_folder,omitempty"`
	XOrigin                 string    `json:"x_origin,omitempty"`
	XFileName               string    `json:"x_filename,omitempty"`
	Body                    string    `json:"body"`
}

// BulkV2Data structure for the bulkv2 API request
type BulkV2Data struct {
	Index   string  `json:"index"`
	Records []Email `json:"records"`
}

// PrintEmail prints the content of an Email struct.
func (e Email) String() string {
	return fmt.Sprintf("MessageID: %s\nDate: %s\nFrom: %s\nTo: %v\nCc: %v\nBcc: %v\nSubject: %s\nMimeVersion: %s\nContentType: %s\nContentTransferEncoding: %s\nXFrom: %s\nXTo: %v\nXCc: %v\nXBcc: %v\nXFolder: %s\nXOrigin: %s\nXFileName: %s\nBody: %s\n",
		e.MessageID, e.Date, e.From, e.To, e.Cc, e.Bcc, e.Subject, e.MimeVersion, e.ContentType, e.ContentTransferEncoding, e.XFrom, e.XTo, e.XCc, e.XBcc, e.XFolder, e.XOrigin, e.XFileName, e.Body)
}

// ParseEmail parses the content of an email file into an Email struct.
func ParseEmail(filePath string) (Email, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Email{}, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	var email Email

	var bodyLines []string
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	inHeaders := true // Flag to indicate whether we are parsing headers or body

	for scanner.Scan() {
		line := scanner.Text()

		// Detect empty line separating headers and body
		if inHeaders && strings.TrimSpace(line) == "" {
			inHeaders = false
			continue
		}
		// Parse headers
		if inHeaders {
			// Find the first colon in the line
			i := strings.Index(line, ":")
			if i != -1 {
				// Split the line into key and value
				key := strings.TrimSpace(line[:i])
				value := strings.TrimSpace(line[i+1:])

				switch strings.ToLower(key) {
				case "message-id":
					if email.MessageID == "" {
						email.MessageID = value
					}
				case "date":
					if email.Date.IsZero() {
						parsedDate, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700 (MST)", value)
						if err != nil {
							log.Fatalf("Error parsing date: %v", err)
							log.Fatalf("File: %s", filePath)
						}
						email.Date = parsedDate
					}
				case "from":
					if email.From == "" {
						email.From = value
					}
				case "to":
					if len(email.To) == 0 {
						email.To = strings.Split(value, ",")
					}
				case "cc":
					if len(email.Cc) == 0 {
						email.Cc = strings.Split(value, ",")
					}
				case "bcc":
					if len(email.Bcc) == 0 {
						email.Bcc = strings.Split(value, ",")
					}
				case "subject":
					if email.Subject == "" {
						email.Subject = value
					}
				case "mime-version":
					if email.MimeVersion == "" {
						email.MimeVersion = value
					}
				case "content-type":
					if email.ContentType == "" {
						email.ContentType = value
					}
				case "content-transfer-encoding":
					if email.ContentTransferEncoding == "" {
						email.ContentTransferEncoding = value
					}
				case "x-from":
					if email.XFrom == "" {
						email.XFrom = value
					}
				case "x-to":
					if len(email.XTo) == 0 {
						email.XTo = strings.Split(value, ",")
					}
				case "x-cc":
					if len(email.XCc) == 0 {
						email.XCc = strings.Split(value, ",")
					}
				case "x-bcc":
					if len(email.XBcc) == 0 {
						email.XBcc = strings.Split(value, ",")
					}
				case "x-folder":
					if email.XFolder == "" {
						email.XFolder = value
					}
				case "x-origin":
					if email.XOrigin == "" {
						email.XOrigin = value
					}
				case "x-filename":
					if email.XFileName == "" {
						email.XFileName = value
					}
				}
			}
		} else {
			// Parse body
			bodyLines = append(bodyLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return Email{}, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	email.Body = strings.Join(bodyLines, "\n")
	return email, nil
}

// IndexEmails processes all files in the directory tree using Depth-First Search.
func readEmails(root string) ([]Email, error) {
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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to file")

func main() {
	// Set up signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Caught termination signal, stopping profiler...")
		pprof.StopCPUProfile()
		os.Exit(1)
	}()
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
	// Define the root folder to start DFS
	rootFolder := "enron_mail_2011_04_02/maildir" // Change this to your desired folder

	// Read email files
	startTime := time.Now()
	emails, err := readEmails(rootFolder)
	duration := time.Since(startTime)
	if err != nil {
		fmt.Printf("Error during DFS: %v\n", err)
		return
	}
	fmt.Printf("%d emails indexed in %v.\n", len(emails), duration)

	// log.Printf("Email example: \n%s\n", emails[0].String())

	// Measure the time taken to send the emails to the bulkv2 API
	startTime = time.Now()
	err = postEmails(emails)
	duration = time.Since(startTime)
	if err != nil {
		fmt.Printf("Error sending emails to bulkv2 API: %v\n", err)
		return
	}
	fmt.Printf("Emails sent to bulkv2 API in %v.\n", duration)

}
