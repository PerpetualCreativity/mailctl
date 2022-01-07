package utils

import (
	"strconv"

	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/spf13/viper"
)

// get login details from config
func loginDetails() (username string, password string,
	imap_server string, imap_port int,
	smtp_server string, smtp_port int,
) {

	ErrComp(viper.IsSet("accounts"), true, "Accounts not set in config file")

	ld := viper.GetStringMapString("accounts." + strconv.Itoa(viper.GetInt("default_account")-1))

	username = ld["username"]
	ErrNComp(username, "", "unable to read `username` field in config")
	password = ld["password"]
	ErrNComp(password, "", "unable to read `password` field in config")

	imap_server = ld["imap_server"]
	ErrNComp(imap_server, "", "unable to read `imap_server` field in config")
	imap_port, err := strconv.Atoi(ld["imap_port"])
	ErrNComp(imap_port, "", "unable to read `imap_port` field in config")
	ErrCheck(err, "`imap_port` field in config not an integer")


	smtp_server = ld["smtp_server"]
	ErrNComp(smtp_server, "", "unable to read `smtp_server` field in config")
	smtp_port, err = strconv.Atoi(ld["smtp_port"])
	ErrNComp(smtp_port, "", "unable to read `smtp_port` field in config")
	ErrCheck(err, "`smtp_port` field in config not an integer")

	return
}

// log in to IMAP
func ImapLogin() *client.Client {
	username, password, server, port, _, _ := loginDetails()
	c, err := client.DialTLS(server+":"+strconv.Itoa(port), nil)
	ErrCheck(err, "Could not connect to server")

	// Login
	ErrCheck(c.Login(username, password), "Could not log in to server")
	return c
}

// log in to SMTP
func SmtpLogin() *smtp.Client {
	username, password, _, _, server, port := loginDetails()
	c, err := smtp.Dial(server + ":" + strconv.Itoa(port))
	ErrCheck(err, "Could not connect to server")
	ErrCheck(c.StartTLS(nil), "Could not start TLS connection")
	_, state := c.TLSConnectionState()
	ErrNComp(state, false, "Could not start TLS connection")
	ErrCheck(c.Auth(sasl.NewLoginClient(username, password)), "Could not log in to server")
	return c
}

