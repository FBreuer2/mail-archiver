package mail

import (
	"errors"
	"log"

	"github.com/FBreuer2/mail-archiver/archive"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func (mA *MailAccount) Login() error {
	// Connect to server
	c, err := client.DialTLS(mA.URL, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Connected to " + mA.URL)

	if err := c.Login(mA.Username, mA.Password); err != nil {
		log.Println(err)
		return err
	}

	log.Println("Logged into " + mA.Username + " on " + mA.URL)

	mA.Client = c

	return nil
}

func (mA *MailAccount) SynchronizeAccount(db *ArchiveDatabase) {

}

func (mA *MailAccount) GetMailboxes() []string {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- mA.Client.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")

	mailboxNames := make([]string, 1)
	for m := range mailboxes {
		mailboxNames = append(mailboxNames, m.Name)
	}

	if err := <-done; err != nil {
		log.Println(err)
	}

	return mailboxNames
}

func (mA *MailAccount) CheckMailbox(mailbox string, db *ArchiveDatabase) error {
	mbox, err := mA.Client.Select(mailbox, false)
	if err != nil {
		log.Println(err)
		return err
	}

	// Get the whole message RAW
	items := []imap.FetchItem{imap.FetchRFC822, imap.FetchEnvelope}

	from := uint32(1)
	to := mbox.Messages

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	totalAmount := to - from

	messageChannel := make(chan *imap.Message, 10)
	doneChannel := make(chan error, 1)
	go func() {
		doneChannel <- mA.Client.Fetch(seqset, items, messageChannel)
	}()

	for {
		select {
		case msg := <-messageChannel:
			if msg == nil {
				log.Println("Server didn't return message")
				break
			}

			hash, err := archive.HashEnvelope(msg.Envelope)

			if err != nil {
				log.Println("Server couldn't hash")
				break
			}

			if db.HasMail(hash) {
				log.Println("Still ", totalAmount, " mails to be archived.")
				totalAmount--
				break
			} else {
				emailText, err := getEmailAsString(msg)
				if err != nil {
					log.Println("Server couldn't get email as string")
					break
				}

				db.ArchiveMail(&archive.ArchivedMail{
					Hash:    hash,
					Subject: msg.Envelope.Subject,
					Date:    msg.Envelope.Date,
					Body:    emailText,
				})

				log.Println("Still ", totalAmount, " mails to be archived.")
				totalAmount--
			}

			if err != nil {
				log.Println(err)
			}

		case err := <-doneChannel:
			if err != nil {
				log.Println(err.Error())
				return err
			}
			return nil

		}
	}
}

func getEmailAsString(message *imap.Message) (string, error) {
	complete := ""

	for _, value := range message.Body {
		len := value.Len()
		buf := make([]byte, len)
		n, err := value.Read(buf)
		if err != nil {
			log.Println(err)
			return "", err
		}

		if n != len {
			log.Println("Didn't read correct length")
			return "", errors.New("Didn't read correct length")
		}

		complete += string(buf)
	}

	return complete, nil
}

func (mA *MailAccount) FetchRaw() {

}

func (mA *MailAccount) Logout() {
	mA.Client.Logout()
}
