package database

import (
	"errors"
	"fmt"

	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
)

func (db *DB) SaveVuz(vuz models.Vuz) (err error) {
	if vuz.Base.Name == "" {
		logger.Log.Printf("Невозможно сохранить пустой вуз!")
		return errors.New("Невозможно сохранить пустой вуз!")
	}
	smt := fmt.Sprintf(`INSERT INTO %s (institutionid, name, description, cost, budget_points, payment_points, budget_places, payment_places, img, logo, url) 
		VALUES ('%s', '%s', '%s', %d, '%f', '%f', %d, %d, '%s', '%s', '%s') ON CONFLICT DO NOTHING;`,
		db.TableVuz, vuz.VuzId, vuz.Base.Name, vuz.Description, vuz.Base.Cost, vuz.Base.Scores.PointsBudget, vuz.Base.Scores.PointsPayment, vuz.Base.Scores.PlacesBudget, vuz.Base.Scores.PlacesPayment, vuz.Base.Image, vuz.Base.Logo, vuz.Base.Url)

	tx, _ := db.Db.Begin()
	_, err = db.Db.Exec(smt)
	if err != nil {

		logger.Log.Printf("[Ошибка] Заведение - %s не было добавлено в базу по причине: %s", vuz.VuzId, err.Error())
		return err
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf("Спарсили заведение - %s", vuz.Base.Name)
	return
}
