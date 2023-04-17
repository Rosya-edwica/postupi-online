package database

import (
	"fmt"

	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
)

func (db *DB) SaveContacts(contacts models.Contacts) {
	smt := fmt.Sprintf(`INSERT INTO %s (vuz_id, address, email, phones, website) VALUES ('%s', '%s', '%s', '%s', '%s')`,
		db.TableContacts, contacts.VuzId, contacts.Address, contacts.Email, contacts.Phone, contacts.WebSite)

	connection := db.Connect()
	defer connection.Close()
	tx, _ := connection.Begin()
	_, err := connection.Exec(smt)
	if err != nil {
		logger.Log.Printf("Ошибка: Contact %s не был добавлен в базу - %s\n", contacts.VuzId, err)
		return
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf("Контакт успешно добавлен: %s\n", contacts.VuzId)
	return

}
