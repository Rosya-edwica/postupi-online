package database

import (
	"database/sql"

	"github.com/Rosya-edwica/postupi-online/pkg/common/telegram"
)

type DB struct {
	ConnectionUrl                string
	TableVuz                     string
	TableSpecialization          string
	TableProgram                 string
	TableProfession              string
	TableContacts                string
	TableVuzToSpecialization     string
	TableProgramToProfession     string
	TableSpecializationToProgram string
}

func (d *DB) Connect() (connection *sql.DB) {
	connection, err := sql.Open("postgres", d.ConnectionUrl)
	checkErr(err)
	return
}

func checkErr(err error) {
	if err != nil {
		telegram.Mailing(err.Error())
		panic(err)
	}
}
