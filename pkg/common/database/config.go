package database

import (
	"database/sql"

	"github.com/Rosya-edwica/postupi-online/pkg/common/telegram"
	_ "github.com/lib/pq"
)

type DB struct {
	Db		*sql.DB
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


func (d *DB) Close() {
	d.Close()
}

func checkErr(err error) {
	if err != nil {
		telegram.Mailing(err.Error())
		panic(err)
	}
}
