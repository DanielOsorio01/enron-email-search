package models

import (
	"fmt"
	"time"
)

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

// PrintEmail prints the content of an Email struct.
func (e Email) String() string {
	return fmt.Sprintf("MessageID: %s\nDate: %s\nFrom: %s\nTo: %v\nCc: %v\nBcc: %v\nSubject: %s\nMimeVersion: %s\nContentType: %s\nContentTransferEncoding: %s\nXFrom: %s\nXTo: %v\nXCc: %v\nXBcc: %v\nXFolder: %s\nXOrigin: %s\nXFileName: %s\nBody: %s\n",
		e.MessageID, e.Date, e.From, e.To, e.Cc, e.Bcc, e.Subject, e.MimeVersion, e.ContentType, e.ContentTransferEncoding, e.XFrom, e.XTo, e.XCc, e.XBcc, e.XFolder, e.XOrigin, e.XFileName, e.Body)
}
