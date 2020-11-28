package archive

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/emersion/go-imap"
)

type ArchivedMail struct {
	ID      uint64 `gorm:"primary_key"`
	Hash    string
	Subject string
	Date    time.Time
	Body    string
}

func HashEnvelope(env *imap.Envelope) (string, error) {
	stringToHash := ""
	stringToHash += env.Date.String()
	stringToHash += env.Subject
	stringToHash += env.InReplyTo

	for _, sender := range env.Sender {
		stringToHash += sender.Address()
	}

	for _, from := range env.From {
		stringToHash += from.Address()
	}

	for _, cc := range env.Cc {
		stringToHash += cc.Address()
	}

	for _, bcc := range env.Bcc {
		stringToHash += bcc.Address()
	}

	sum := sha256.Sum256([]byte(stringToHash))

	return hex.EncodeToString(sum[:]), nil
}
