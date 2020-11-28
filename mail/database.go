package mail

import (
	"log"

	"github.com/FBreuer2/mail-archiver/archive"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ArchiveDatabase struct {
	database *gorm.DB
}

func OpenDatabase(databasePath string) (*ArchiveDatabase, error) {
	db, err := gorm.Open("sqlite3", databasePath)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&MailAccount{})
	db.AutoMigrate(&SynchronizationConfiguration{})
	db.AutoMigrate(&archive.ArchivedMail{})

	return &ArchiveDatabase{
		database: db,
	}, nil
}

func (aD *ArchiveDatabase) HasMail(hash string) bool {
	var archivedEmail archive.ArchivedMail
	aD.database.Where("hash = ?", hash).First(&archivedEmail)

	return archivedEmail.Hash == hash
}

func (aD *ArchiveDatabase) ArchiveMail(archivedMail *archive.ArchivedMail) {
	aD.database.Create(archivedMail)
}

func (aD *ArchiveDatabase) Close() {
	aD.database.Close()
}
