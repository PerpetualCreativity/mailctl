package utils

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func MoveMail(c *client.Client, id int, from string, to string) {
	_, err := c.Select(from, false)
	ErrCheck(err, "From folder does not exist")
	seqset := new(imap.SeqSet)
	seqset.AddNum(uint32(id))
	err = c.Move(seqset, to)
	ErrCheck(err, "Could not move message")
}
