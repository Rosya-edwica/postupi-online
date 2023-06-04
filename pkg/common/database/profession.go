package database

import (
	"fmt"

	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
)

func (db *DB) SaveProfession(profession models.Profession) {
	msgError2 := fmt.Sprintf("[Ошибка] Профессия - %s не была добавлена в базу по причине:", profession.Name)
	if profession.Name == "" {
		return
	}
	smt := fmt.Sprintf(`INSERT INTO %s (programid, name, img) VALUES ('%d', '%s', '%s')`, db.TableProfession, profession.ProgramId, profession.Name, profession.Image)
	tx, _ := db.Db.Begin()
	_, err := tx.Exec(smt)
	if err != nil {
		logger.Log.Println(msgError2 + err.Error())
		return
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf("Успешно сохранили профессию %s\n", profession.Name)
	return
}
