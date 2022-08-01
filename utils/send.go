package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-smtp"
)

func Send(c *smtp.Client, to string, from string, subject string, body string) error {
	toList := strings.Split(to, "; ")

	err := c.Mail(from, nil)
	if err != nil {
		return errors.New("error in message setup")
	}
	for _, t := range toList {
		err = c.Rcpt(t)
		if err != nil {
			return errors.New("error in message setup")
		}
	}
	wc, err := c.Data()
	if err != nil {
		return errors.New("error in message setup")
	}

	fmt.Fprintf(wc,
		"Date: %s\r\nFrom: %s\r\nSubject: %s\r\nTo: %s\r\n\r\n%s\r\n.\r\n",
		time.Now().Format("Mon, 2 Jan 2006 03:04:05 -0700"), from, subject, to, body,
	)
	wc.Close()

	return nil
}
