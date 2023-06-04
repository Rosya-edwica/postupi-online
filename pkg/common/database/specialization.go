package database

import (
	"errors"
	"fmt"

	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
)

func (db *DB) SaveSpecialization(spec models.Specialization) (err error) {
	if spec.Base.Name == "" {
		logger.Log.Printf("Не удалось добавить пустую специализацию: %T", spec)
		return errors.New("Не удалось добавить пустую специализацию")
	}
	
	smt := fmt.Sprintf(`INSERT INTO %s (specid, institutionid, name, description, direction, cost, budget_points, payment_points, budget_places, payment_places, img, url)
	VALUES ('%s', '%s', '%s', '%s', '%s', %d, '%f', '%f', %d, %d, '%s', '%s') ON CONFLICT DO NOTHING`,
		db.TableSpecialization, spec.SpecId, spec.VuzId, spec.Base.Name, spec.Description, spec.Base.Direction, spec.Base.Cost, spec.Base.Scores.PointsBudget, spec.Base.Scores.PointsPayment, spec.Base.Scores.PlacesBudget, spec.Base.Scores.PlacesPayment, spec.Base.Image, spec.Base.Url)
	tx, _ := db.Db.Begin()
	_, err = db.Db.Exec(smt)
	if err != nil {
		logger.Log.Printf("[Ошибка] Специализация - %s не была добавлена в базу по причине: %s", spec.SpecId, err.Error())
		return nil
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf("Специализация успешно добавлена - %s", spec.Base.Name)
	return
}
