package email

import (
	"log"
	"sync"
	"sync/atomic"
)

func SendEmailBatches(emailsQueue <-chan Email, batchSize int, wg *sync.WaitGroup, client *ZincClient, sentEmails *atomic.Uint64) {
	defer wg.Done()

	var emails []Email

	for email := range emailsQueue {
		emails = append(emails, email)
		if len(emails) == batchSize {
			if err := client.SendEmails(emails); err != nil {
				log.Printf("Error sending emails: %v", err)
			}
			emails = nil
			sentEmails.Add(uint64(batchSize))
		}
	}

	// Send the remaining emails
	if len(emails) > 0 {
		if err := client.SendEmails(emails); err != nil {
			log.Printf("Error sending emails: %v", err)
		}
		sentEmails.Add(uint64(len(emails)))
	}
}
