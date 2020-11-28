package main

import (
	"flag"
	"log"

	"github.com/FBreuer2/mail-archiver/mail"
)

func main() {

	var configurationDatabasePath string
	flag.StringVar(&configurationDatabasePath, "databasePath", "./archiver.sqlite", "The path to the configuration database.")

	mailAccountDatabase, err := mail.OpenDatabase(configurationDatabasePath)

	if err != nil {
		log.Println("Couldn't open database!")
		return
	}

	yourMailAccount := mail.MailAccount{URL: "url:port", Username: "username", Password: "password"}

	yourMailAccount.Login()

	yourMailAccount.CheckMailbox("INBOX", mailAccountDatabase)

	defer yourMailAccount.Logout()

	//sigs := make(chan os.Signal, 1)
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//sig := <-sigs

	defer mailAccountDatabase.Close()

}
