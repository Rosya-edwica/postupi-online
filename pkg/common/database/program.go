package database

import (
	"errors"
	"fmt"
	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
	"strings"
)


func (db *DB) SaveProgram(program models.Program) (programId int, err error) {
	if program.Base.Name == "" {
		logger.Log.Printf("Невозможно сохранить пустую программу: %T", program)
		return 0, errors.New("Невозможно сохранить пустую программу:" + err.Error())
	}

	smt := fmt.Sprintf(`INSERT INTO %s (specid, institutionid, name, description, direction, form, subjects, cost, has_professions,
		 budget_points, payment_points, budget_places, payment_places, img, url) 
		 VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %t, '%f', '%f', %d, %d, '%s', '%s' ) RETURNING programid`,
		db.TableProgram, program.SpecId, program.VuzId, program.Base.Name, program.Base.Description, program.Base.Direction, program.Form, strings.Join(program.Exams, "|"),
		program.Base.Cost, program.HasProfessions, program.Base.Scores.PointsBudget, program.Base.Scores.PointsPayment, program.Base.Scores.PlacesBudget,
		program.Base.Scores.PlacesPayment, program.Base.Image, program.Base.Url)

	tx, _ := db.Db.Begin()
	err = db.Db.QueryRow(smt).Scan(&programId)
	if err != nil {
		logger.Log.Printf("[Ошибка] Программа - %s не была добавлена в базу по причине: %s", program.ProgramId, err.Error())
		return 0, nil
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf("Программа успешно добавлена - %s", program.Base.Name)
	return programId, nil

}
