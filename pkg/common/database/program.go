package database

import (
	"errors"
	"fmt"
	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
	"github.com/lib/pq"
	"strings"
)

// TODO: Константы  с успешными и провальными сообщениями

func (db *DB) SaveProgram(program models.Program) (connectionId int, err error) {
	msgError1 := fmt.Sprintf("Невозможно сохранить пустую программу: %T", program)
	msgError2 := fmt.Sprintf("[Ошибка] Программа - %s не была добавлена в базу по причине:", program.ProgramId)
	msgSuccess := fmt.Sprintf("Программа успешно добавлена - %s", program.Base.Name)

	programId := db.CheckProgramExists(program.ProgramId)
	if programId != "" {
		return db.SetConnectionBetweenSpecializationAndProgram(program.VuzId, program.SpecId, program.ProgramId), nil
	}

	if program.Base.Name == "" {
		logger.Log.Printf(msgError1)
		return 0, errors.New(msgError1)
	}

	connection := db.Connect()
	defer connection.Close()

	smt := fmt.Sprintf(`INSERT INTO %s (id, name, description, direction, form, subjects, cost, has_professions,
		 budget_points, payment_points, budget_places, payment_places, image, url) 
		 VALUES ('%s', '%s', '%s', '%s', '%s', '%s', %d, %t, '%f', '%f', %d, %d, '%s', '%s' )`,
		db.TableProgram, program.ProgramId, program.Base.Name, program.Base.Description, program.Base.Direction, program.Form, pq.Array(makeArrayForPostgreSQL(program.Exams)),
		program.Base.Cost, program.HasProfessions, program.Base.Scores.PointsBudget, program.Base.Scores.PointsPayment, program.Base.Scores.PlacesBudget,
		program.Base.Scores.PlacesPayment, program.Base.Image, program.Base.Url)
	tx, _ := connection.Begin()
	_, err = connection.Exec(smt)
	if err != nil {
		logger.Log.Println(msgError2 + err.Error())
		return 0, nil
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Println(msgSuccess)
	return db.SetConnectionBetweenSpecializationAndProgram(program.VuzId, program.SpecId, program.ProgramId), nil

}

func (db *DB) SetConnectionBetweenSpecializationAndProgram(vuzId string, specId string, programId string) (connectionBetweenSpecializationAndProgram int) {
	// msgError1 := fmt.Sprintf("[Ошибка] Не удалось наладить связь между специализацией и программой [%s -> %s]:", spec.VuzId, spec.SpecId)
	// msgSuccess := fmt.Sprintf("Установлена связь между вузом и специализацией [%s -> %s]", spec.VuzId, spec.SpecId)

	connectionId := db.GetConnectionBetweenVuzAndSpecialization(vuzId, specId)
	connection := db.Connect()
	defer connection.Close()

	smt := fmt.Sprintf(`INSERT INTO %s(vuz_specialization, program_id) VALUES (%d, '%s')`, db.TableSpecializationToProgram, connectionId, programId)
	tx, _ := connection.Begin()
	_, err := connection.Exec(smt)
	if err != nil {
		logger.Log.Printf("[Ошибка] Не удалось наладить связь между специализацией и программой [%d -> %s]: %s\n", connectionId, programId, err)
		return
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf("Установлена связь между специализацией и программой [%d -> %s]", connectionId, programId)
	connectionBetweenSpecializationAndProgram = db.GetConnectionBetweenSpecializationAndProgram(connectionId, programId)
	return
}

func (db *DB) GetConnectionBetweenVuzAndSpecialization(vuzId string, specId string) (id int) {
	connection := db.Connect()
	defer connection.Close()

	query := fmt.Sprintf(`SELECT id FROM %s WHERE vuz_id='%s' AND spec_id='%s'`, db.TableVuzToSpecialization, vuzId, specId)
	rows, err := connection.Query(query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id)
		return id
	}
	logger.Log.Printf("[Ошибка] Не удалось найти подобную связь [%s -> %s]", vuzId, specId)
	return 0
}

func (db *DB) GetConnectionBetweenSpecializationAndProgram(vuzAndSpecId int, programId string) (id int) {
	connection := db.Connect()
	defer connection.Close()

	query := fmt.Sprintf(`SELECT id FROM %s WHERE vuz_specialization=%d AND program_id='%s'`, db.TableSpecializationToProgram, vuzAndSpecId, programId)
	rows, err := connection.Query(query)
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&id)
		return id
	}
	logger.Log.Printf("Не удалось найти подобную связь между программой и специализацией [%s->%d]", programId, vuzAndSpecId)
	return 0
}

func (db *DB) CheckProgramExists(programId string) (id string) {
	connection := db.Connect()
	defer connection.Close()

	query := fmt.Sprintf(`SELECT name FROM %s WHERE id='%s'`, db.TableProgram, programId)
	rows, err := connection.Query(query)
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		return programId
	}
	return id
}

func makeArrayForPostgreSQL(array []string) string {
	if len(array) > 0 {
		return "{" + strings.Join(array, ",") + "}"
	} else {
		return ""
	}
}
