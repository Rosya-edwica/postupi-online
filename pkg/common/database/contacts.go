package database

import (
	"fmt"

	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
)

func (db *DB) SaveContacts(contacts models.Contacts) {
	smt := fmt.Sprintf(`INSERT INTO %s (institutionid, address, email, phones, website) VALUES ('%s', '%s', '%s', '%s', '%s')`,
		db.TableContacts, contacts.VuzId, contacts.Address, contacts.Email, contacts.Phone, contacts.WebSite)
	tx, _ := db.Db.Begin()
	_, err := db.Db.Exec(smt)
	if err != nil {
		logger.Log.Printf("Ошибка: Contact %s не был добавлен в базу - %s\n", contacts.VuzId, err)
		return
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf("Контакт успешно добавлен: %s\n", contacts.VuzId)
	return

}
