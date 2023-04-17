package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/Rosya-edwica/postupi-online/pkg/common/database"
	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
	"github.com/Rosya-edwica/postupi-online/pkg/scraper"
)

func main() {
	start := time.Now().Unix()
	logger.Log.Println("Программа запущена")

	parserType := defineParserType()

	if err := initConfig(parserType); err != nil {
		panic(fmt.Sprintf("Не загрузился конфиг: %s", err))
	}
	db := initDatabase()
	scr := initScraper()

	scraper.Run(scr, db)
	fmt.Printf("Время выполнения программы: %d секунд.\n", time.Now().Unix()-start)
}

func initDatabase() *database.DB {
	db := &database.DB{
		ConnectionUrl:                viper.GetString("db.postgresUrl"),
		TableVuz:                     viper.GetString("db.table.vuz"),
		TableSpecialization:          viper.GetString("db.table.specialization"),
		TableProgram:                 viper.GetString("db.table.program"),
		TableProfession:              viper.GetString("db.table.profession"),
		TableContacts:                viper.GetString("db.table.contacts"),
		TableVuzToSpecialization:     viper.GetString("db.table.vuz_to_specialization"),
		TableSpecializationToProgram: viper.GetString("db.table.specialization_to_program"),
		TableProgramToProfession:     viper.GetString("db.table.program_to_profession"),
	}
	logger.Log.Println("Содержание объекта Базы данных:", db)
	return db
}

func initScraper() scraper.Scraper {
	scr := scraper.Scraper{
		Domain:         viper.GetString("domain"),
		FormEducations: strings.Split(viper.GetString("formEducations"), "|"),
	}
	logger.Log.Println("Содержание объекта парсера:", scr)
	return scr

}

func initConfig(typeParsing string) error {
	logger.Log.Println("Загружаем переменные окружения под парсер для", typeParsing)
	configPath := getConfigPath(typeParsing)
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func getConfigPath(typeParsing string) (path string) {
	if typeParsing == "college" {
		return "configs/college"
	}
	return "configs/vuz"
}

func defineParserType() (parserType string) {
	errorMessage := "Передайте вместе с командой `go run main.go` обязательный параметр `college` или `vuz`"
	command_args := os.Args
	if len(command_args) == 1 {
		panic(errorMessage)
	}
	switch command_args[1] {
	case "vuz":
		parserType = "vuz"
	case "college":
		parserType = "college"
	default:
		panic(errorMessage)
	}
	return
}
