package utils

import (
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func MoveMail(c *client.Client, id int, from string, to string) error {
	_, err := c.Select(from, false)
	if err != nil {
		return errors.New("from folder does not exist")
	}
	seqset := new(imap.SeqSet)
	seqset.AddNum(uint32(id))
	err = c.Move(seqset, to)
	if err != nil {
		return errors.New("could not move message")
	}
	return nil
}
