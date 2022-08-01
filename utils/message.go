package utils

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

func GetMessage(c *client.Client, id int, folder string) (string, string, error) {
	if err := c.Unselect(); err != nil && !errors.Is(err, client.ErrNoMailboxSelected) {
		return "", "", errors.New("could not unselect current folder")
	}
	status, err := c.Select(folder, false)
	if err != nil {
		return "", "", errors.New("could not select folder")
	}
	if status.Messages == 0 {
		return "", "", errors.New("no messages in folder")
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(uint32(id))

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	section := &imap.BodySectionName{}
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}, messages)
	}()

	if <-done != nil {
		return "", "", errors.New("failed to get message parts")
	}

	msg := <-messages

	r := msg.GetBody(section)

	if r == nil {
		return "", "", errors.New("server did not return a message body")
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		return "", "", errors.New("could not read message")
	}

	var body string
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "", errors.New("could not read message")
		}
		switch p.Header.(type) {
		case *mail.InlineHeader:
			b, err := ioutil.ReadAll(p.Body)
			if err != nil {
				return "", "", errors.New("could not read message")
			}
			body = string(b)
			break
		}
	}
	return msg.Envelope.Subject, body, nil
}

func NewMessage(c *client.Client, from *mail.Address, to []*mail.Address, body string) error {
	h := mail.Header{}
	h.SetDate(time.Now())
	h.SetAddressList("From", []*mail.Address{from})
	h.SetAddressList("To", to)

	var msg bytes.Buffer
	mw, _ := mail.CreateWriter(&msg, h)
	tw, _ := mw.CreateInline()
	var th mail.InlineHeader
	th.Set("Content-Type", "text/plain")
	w, _ := tw.CreatePart(th)
	io.WriteString(w, body)
	w.Close()
	tw.Close()

	draftsBox := FindMailbox(c, "\\Drafts", "Drafts")
	err := c.Append(draftsBox, nil, time.Now(), &msg)
	if err != nil {
		return errors.New("could not add new message template to drafts")
	}
	return nil
}
