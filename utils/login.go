package utils

import (
	"errors"
	"strconv"

	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/spf13/viper"
)

type logins struct {
	username   string
	password   string
	imapServer string
	imapPort   int
	smtpServer string
	smtpPort   int
}

// get login details from config
func loginDetails(configIndex int) (logins, error) {
	l := logins{}

	if !viper.IsSet("accounts") {
		return l, errors.New("accounts not set in config file")
	}

	ld := viper.GetStringMapString("accounts." + strconv.Itoa(configIndex))

	l.username = ld["username"]
	if l.username == "" {
		return l, errors.New("unable to read `username` field in config")
	}
	l.password = ld["password"]
	if l.password == "" {
		return l, errors.New("unable to read `password` field in config")
	}

	l.imapServer = ld["imap_server"]
	if l.imapServer == "" {
		return l, errors.New("unable to read `imap_server` field in config")
	}

	var err error
	l.imapPort, err = strconv.Atoi(ld["imap_port"])
	if err != nil {
		return l, errors.New("`imap_port` field in config not an integer")
	}

	l.smtpServer = ld["smtp_server"]
	if l.smtpServer == "" {
		return l, errors.New("unable to read `password` field in config")
	}
	l.smtpPort, err = strconv.Atoi(ld["smtp_port"])
	if err != nil {
		return l, errors.New("`smtp_port` field in config not an integer")
	}

	return l, nil
}

// log in to IMAP
func ImapLogin(account ...int) (*client.Client, error) {
	index := viper.GetInt("default_account") - 1
	if len(account) > 0 {
		index = account[0]
	}
	logins, err := loginDetails(index)
	if err != nil {
		return nil, err
	}
	c, err := client.DialTLS(logins.imapServer+":"+strconv.Itoa(logins.imapPort), nil)
	if err != nil {
		return nil, errors.New("could not connect to server")
	}

	// Login
	err = c.Login(logins.username, logins.password)
	if err != nil {
		return nil, errors.New("could not login to server")
	}
	return c, nil
}

// log in to SMTP
func SmtpLogin(account ...int) (*smtp.Client, error) {
	index := viper.GetInt("default_account") - 1
	if len(account) > 0 {
		index = account[0]
	}
	logins, err := loginDetails(index)
	if err != nil {
		return nil, err
	}
	c, err := smtp.Dial(logins.smtpServer + ":" + strconv.Itoa(logins.smtpPort))
	if err != nil {
		return nil, errors.New("could not connect to server")
	}
	if err := c.StartTLS(nil); err != nil {
		return nil, errors.New("could not start TLS connection")
	}
	_, state := c.TLSConnectionState()
	if !state {
		return nil, errors.New("could not start TLS connection")
	}
	if err := c.Auth(sasl.NewLoginClient(logins.username, logins.password)); err != nil {
		return nil, errors.New("could not log in to server")
	}
	return c, nil
}
