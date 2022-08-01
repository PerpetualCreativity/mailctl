package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/spf13/viper"
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

func ListMailboxes(c *client.Client, configIndex ...int) []string {
	mailboxes := getMailboxes(c)

	ci := viper.GetInt("default_account")
	if len(configIndex) > 0 {
		ci = configIndex[0]
	}
	ignoreList := viper.GetStringSlice("accounts." + strconv.Itoa(ci-1) + ".ignore_mailboxes")

	var mailboxNames []string
	for m := range mailboxes {
		exclude := false
		for _, i := range ignoreList {
			if strings.Contains(m.Name, i) {
				exclude = true
				break
			}
		}
		if !exclude {
			mailboxNames = append(mailboxNames, m.Name)
		}
	}

	return mailboxNames
}

type Message struct {
	SeqNum  uint32
	Sender  string
	Subject string
}

func ListMessages(c *client.Client, folder string, number uint32, offset uint32) ([]Message, error) {
	mailbox, err := c.Select(folder, false)
	if err != nil {
		return nil, errors.New("could not select mailbox")
	}

	from := offset + 1
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

	var msgs []Message
	for msg := range messages {
		sender := ""
		if s := msg.Envelope.From; len(s) > 0 {
			sender = msg.Envelope.From[0].PersonalName
			if sender == "" {
				sender = msg.Envelope.From[0].Address()
			}
		} else {
			sender = "[no sender]"
		}
		msgs = append(msgs, Message{
			SeqNum:  msg.SeqNum,
			Sender:  sender,
			Subject: msg.Envelope.Subject,
		})
	}
	<-done

	return msgs, nil
}
