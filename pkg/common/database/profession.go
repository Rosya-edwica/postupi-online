package database

import (
	"fmt"

	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
)

func (db *DB) SaveProfession(profession models.Profession, connectionBetweenProgramAndProfession int) {
	msgError2 := fmt.Sprintf("[Ошибка] Профессия - %s не была добавлена в базу по причине:", profession.Name)
	if profession.Name == "" {
		return
	}

	professionId := db.CheckProfessionExists(profession.Name)
	if professionId != 0 {
		db.SetConnectionBetweenProgramAndProfession(professionId, connectionBetweenProgramAndProfession)
		return
	}
	connection := db.Connect()
	defer connection.Close()

	smt := fmt.Sprintf(`INSERT INTO %s (name, image) VALUES ('%s', '%s')`, db.TableProfession, profession.Name, profession.Image)
	tx, _ := connection.Begin()
	_, err := tx.Exec(smt)
	if err != nil {
		logger.Log.Println(msgError2 + err.Error())
		return
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf("Успешно сохранили профессию %s\n", profession.Name)

	professionId = db.CheckProfessionExists(profession.Name)
	db.SetConnectionBetweenProgramAndProfession(professionId, connectionBetweenProgramAndProfession)
	return
}

func (db *DB) SetConnectionBetweenProgramAndProfession(professionId int, programAndProfessionId int) {
	connection := db.Connect()
	defer connection.Close()

	smt := fmt.Sprintf(`INSERT INTO %s(vuz_specialization_program, profession_id) VALUES (%d, %d)`, db.TableProgramToProfession, programAndProfessionId, professionId)
	tx, _ := connection.Begin()
	_, err := connection.Exec(smt)
	if err != nil {
		logger.Log.Printf("Не удалось наладить связь между профессией и программой - %s\n", err)
		return
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf("Установлена связь между профессией и программой [%d -> %d]", programAndProfessionId, professionId)
}

func (db *DB) CheckProfessionExists(professionName string) (professionId int) {
	connection := db.Connect()
	defer connection.Close()

	query := fmt.Sprintf(`SELECT id FROM %s WHERE name='%s'`, db.TableProfession, professionName)
	rows, err := connection.Query(query)
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&professionId)
		return professionId
	}
	return
}
