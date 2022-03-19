package utils

import (
	"io"
	"io/ioutil"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

func GetMessage(c *client.Client, id int, folder string) (string, string) {
	fc.ErrExp(c.Unselect(), client.ErrNoMailboxSelected, "Could not unselect current folder")
	status, err := c.Select(folder, false)
	fc.ErrCheck(err, "Could not select folder")
	fc.ErrNComp(status.Messages, 0, "No messages in folder")

	seqset := new(imap.SeqSet)
	seqset.AddNum(uint32(id))

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	section := &imap.BodySectionName{}
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}, messages)
	}()

	fc.ErrCheck(<-done, "Failed to get message parts")

	msg := <-messages

	subject := func() string {
		defer func() {
			fc.ErrCheck(recover(), "Subject not defined")
		}()
		return msg.Envelope.Subject
	}()

	r := func() imap.Literal {
		defer func() {
			fc.ErrCheck(recover(), "Invalid ID")
		}()
		return msg.GetBody(section)
	}()

	fc.ErrNComp(r, nil, "Server did not return a message body")

	mr, err := mail.CreateReader(r)
	fc.ErrCheck(err, "Could not read message")

	var body string
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		fc.ErrCheck(err, "Could not read message")
		switch p.Header.(type) {
		case *mail.InlineHeader:
			b, err := ioutil.ReadAll(p.Body)
			fc.ErrCheck(err, "Could not read message")
			body = string(b)
			break
		}
	}
	return subject, body
}
