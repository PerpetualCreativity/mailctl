package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-smtp"
)

func Send(c *smtp.Client, to string, from string, subject string, body string) {
	toList := strings.Split(to, "; ")

	err := c.Mail(from, nil)
	fc.ErrCheck(err, "Error in setup of message")
	for _, t := range toList {
		err = c.Rcpt(t)
		fc.ErrCheck(err, "Error in setup of message")
	}
	wc, err := c.Data()
	fc.ErrCheck(err, "Error in setup of message")

	fmt.Fprintf(wc,
		"Date: %s\r\nFrom: %s\r\nSubject: %s\r\nTo: %s\r\n\r\n%s\r\n.\r\n",
		time.Now().Format("Mon, 2 Jan 2006 03:04:05 -0700"), from, subject, to, body,
	)
	wc.Close()
}
