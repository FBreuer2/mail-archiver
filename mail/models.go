package mail

import (
	"github.com/emersion/go-imap/client"
	"github.com/jinzhu/gorm"
)

type SynchronizationConfiguration struct {
	WaitTimeInSeconds uint `gorm:"default:10"`
}

type MailAccount struct {
	gorm.Model
	URL      string
	Username string
	Password string
	Client   *client.Client
}
