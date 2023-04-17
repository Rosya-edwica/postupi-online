package database

import (
	"errors"
	"fmt"

	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
)

func (db *DB) SaveSpecialization(spec models.Specialization) (err error) {
	msgError1 := fmt.Sprintf("Не удалось добавить пустую специализацию: %T", spec)
	msgError2 := fmt.Sprintf("[Ошибка] Специализация - %s не была добавлена в базу по причине:", spec.SpecId)
	msgSuccess := fmt.Sprintf("Специализация успешно добавлена - %s", spec.Base.Name)

	specId := db.CheckSpecializationExists(spec.SpecId)
	if specId != "" {
		db.SetConnectionBetweenVuzAndSpecialization(spec.VuzId, spec.SpecId)
		return
	}
	if spec.Base.Name == "" {
		logger.Log.Println(msgError1)
		return errors.New(msgError1)
	}

	connection := db.Connect()
	defer connection.Close()

	smt := fmt.Sprintf(`INSERT INTO %s (id, name, description, direction, cost, budget_points, payment_points, budget_places, payment_places, image, url)
	VALUES ('%s', '%s', '%s', '%s', %d, '%f', '%f', %d, %d, '%s', '%s')`,
		db.TableSpecialization, spec.SpecId, spec.Base.Name, spec.Description, spec.Base.Direction, spec.Base.Cost, spec.Base.Scores.PointsBudget, spec.Base.Scores.PointsPayment, spec.Base.Scores.PlacesBudget, spec.Base.Scores.PlacesPayment, spec.Base.Image, spec.Base.Url)
	tx, _ := connection.Begin()
	_, err = connection.Exec(smt)
	if err != nil {
		logger.Log.Println(msgError2 + err.Error())
		return nil
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf(msgSuccess)
	db.SetConnectionBetweenVuzAndSpecialization(spec.VuzId, spec.SpecId)
	return
}

func (db *DB) SetConnectionBetweenVuzAndSpecialization(vuzId string, specId string) {
	msgError1 := fmt.Sprintf("[Ошибка] Не удалось наладить связь между вузом и специализацией [%s -> %s]:", vuzId, specId)
	msgSuccess := fmt.Sprintf("Установлена связь между вузом и специализацией [%s -> %s]", vuzId, specId)

	connection := db.Connect()
	defer connection.Close()
	smt := fmt.Sprintf(`INSERT INTO %s(vuz_id, spec_id) VALUES ('%s', '%s')`, db.TableVuzToSpecialization, vuzId, specId)
	tx, _ := connection.Begin()
	_, err := connection.Exec(smt)
	if err != nil {
		logger.Log.Println(msgError1 + err.Error())
		return
	}
	err = tx.Commit()
	checkErr(err)

	logger.Log.Println(msgSuccess)

}

func (db *DB) CheckSpecializationExists(specId string) (id string) {
	connection := db.Connect()
	defer connection.Close()
	query := fmt.Sprintf(`SELECT name FROM %s WHERE id='%s'`, db.TableSpecialization, specId)
	rows, err := connection.Query(query)
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		return specId
	}
	return
}
