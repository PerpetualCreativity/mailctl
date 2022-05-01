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

type message struct {
	SeqNum	uint32
	Sender  string
	Subject string
}

func ListMessages(c *client.Client, folder string, number uint32) []message {
	mailbox, err := c.Select(folder, false)
	fc.ErrCheck(err, "Could not select mailbox")

	from := uint32(1)
	to := mailbox.Messages

	if mailbox.Messages > number {
		from = mailbox.Messages - number
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, number)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	msgs := []message{}
	for msg := range messages {
		sender := ""
		if s := msg.Envelope.From; len(s) > 0 {
			sender = msg.Envelope.From[0].PersonalName
		}
		msgs = append(msgs, message{
			SeqNum: msg.SeqNum,
			Sender: sender,
			Subject: msg.Envelope.Subject,
		})
	}
	fc.ErrCheck(<-done, "No messsages in this folder")

	return msgs
}
