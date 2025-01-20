package email

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

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

func ProcessEmailFiles(emailFiles <-chan string, emailQueue chan<- Email, wg *sync.WaitGroup) {
	// Reads email file paths from the emailFiles channel
	// and sends parsed emails to the emailQueue channel.

	defer wg.Done() // Decrement the counter when the function completes

	// Read email file paths from the channel
	for path := range emailFiles {
		email, err := ParseEmail(path)
		if err != nil {
			log.Printf("Warning: Failed to parse file %s: %v\n", path, err)
			continue
		}

		emailQueue <- email // Send the parsed email to the emailQueue channel
	}
	log.Println("Finished processing email files.")
}
