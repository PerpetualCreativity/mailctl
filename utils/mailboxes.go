package utils

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func getMailboxes(c *client.Client) chan *imap.MailboxInfo {
	c.Unselect()
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	return mailboxes
}

func FindMailbox(c *client.Client, attr string, fallback string) string {
	mailboxes := getMailboxes(c)
	box := fallback
	for m := range mailboxes {
		for _, a := range m.Attributes {
			if a == attr {
				box = m.Name
				break
			}
		}
	}
	return box
}

func ListMailboxes(c *client.Client) []string {
	mailboxes := getMailboxes(c)
	var mailboxNames []string
	for m := range mailboxes {
		mailboxNames = append(mailboxNames, m.Name)
	}
	return mailboxNames
}
