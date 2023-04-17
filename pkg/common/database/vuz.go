package database

import (
	"errors"
	"fmt"

	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
)

func (db *DB) SaveVuz(vuz models.Vuz) (err error) {
	msgError1 := "Невозможно сохранить пустой вуз!"
	msgError2 := fmt.Sprintf("[Ошибка] Заведение - %s не было добавлено в базу по причине:", vuz.VuzId)
	msgSucces := fmt.Sprintf("Спарсили заведение - %s", vuz.Base.Name)

	if vuz.Base.Name == "" {
		logger.Log.Printf(msgError1)
		return errors.New(msgError1)
	}
	smt := fmt.Sprintf(`INSERT INTO %s (id, name, description, city, cost, budget_points, payment_points, budget_places, payment_places, image, logo, url) 
		VALUES ('%s', '%s', '%s', '%s', %d, '%f', '%f', %d, %d, '%s', '%s', '%s')`,
		db.TableVuz, vuz.VuzId, vuz.Base.Name, vuz.Description, vuz.City, vuz.Base.Cost, vuz.Base.Scores.PointsBudget, vuz.Base.Scores.PointsPayment, vuz.Base.Scores.PlacesBudget, vuz.Base.Scores.PlacesPayment, vuz.Base.Image, vuz.Base.Logo, vuz.Base.Url)
	connection := db.Connect()
	defer connection.Close()

	tx, _ := connection.Begin()
	_, err = connection.Exec(smt)
	if err != nil {

		logger.Log.Println(msgError2 + err.Error())
		return nil
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf(msgSucces)
	return
}
