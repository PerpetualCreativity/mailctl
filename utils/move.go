package utils

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func MoveMail(c *client.Client, id int, from string, to string) {
	_, err := c.Select(from, false)
	fc.ErrCheck(err, "From folder does not exist")
	seqset := new(imap.SeqSet)
	seqset.AddNum(uint32(id))
	err = c.Move(seqset, to)
	fc.ErrCheck(err, "Could not move message")
}
