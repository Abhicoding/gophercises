package main

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	postgresUser = "POSTGRES_USER"
	postgresPwd  = "POSTGRES_PASSWORD"
	postgreDB    = "POSTGRES_DB"
	tableName    = "phone_numbers"
)

var (
	port     string
	address  = "127.0.0.1"
	user     string
	password string
	dbName   string
)

type conf struct {
	Services `yaml:"services"`
}

type Services struct {
	Db `yaml:"db"`
}

type Db struct {
	Environment []string `yaml:"environment"`
	Ports       string   `yaml:"port"`
}

func main() {

	user = "root"
	password = "root"
	dbName = "PhoneNumbers"
	port = "5432"

	//TODO: read config from yaml
	// var config conf
	// file, err := ioutil.ReadFile("./docker-compose.yml")
	// if err != nil {
	// 	fmt.Errorf("failed trying to read file docker-compose.yml err: %s", err.Error())
	// 	return
	// }

	// err = yaml.Unmarshal(file, &config)
	// if err != nil {
	// 	fmt.Errorf("failed trying to parse yaml err: %s", err.Error())
	// 	return
	// }
	// getConfig(config.Services.Db)

	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s sslmode=disable", address, port, user, password, dbName)

	db, err := gorm.Open(postgres.Open(psqlConn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("error connecting to database %v", err))
	}
	createEntries(db)
	fixEntries(db)
}

func getConfig(db Db) {
	for _, v := range db.Environment {
		dbConfig := strings.Split(v, "=")
		switch dbConfig[0] {
		case postgresUser:
			user = dbConfig[1]
		case postgresPwd:
			password = dbConfig[1]
		case postgreDB:
			dbName = dbConfig[1]
		}
	}
}

var phoneNumbers = []string{
	"1234567890",
	"123 456 7891",
	"(123) 456 7892",
	"(123) 456-7893",
	"123-456-7894",
	"123-456-7890",
	"1234567892",
	"(123)456-7892",
}

type phoneNumber struct {
	ID    int    `gorm:"id"`
	Value string `gorm:"number"`
}

func createEntries(db *gorm.DB) {
	db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
	db.Exec(fmt.Sprintf(`CREATE TABLE %s (
		id SERIAL,
		number VARCHAR(24)
	)`, tableName))
	for _, number := range phoneNumbers {
		db.Exec(fmt.Sprintf("INSERT INTO %s (number) VALUES ('%s')", tableName, number))
	}
}

func fixEntries(db *gorm.DB) {
	rows, err := db.Raw(fmt.Sprintf("SELECT * from %s", tableName)).Rows()
	if err != nil {
		panic(fmt.Errorf("select query err: %v", err))
	}
	defer rows.Close()
	var ph = phoneNumber{}
	for rows.Next() {
		rows.Scan(&ph.ID, &ph.Value)
		formattedNumber := formatNumber(ph.Value)
		if ph.Value == formattedNumber {
			continue
		}
		if exists := checkIfDuplicateExists(db, phoneNumber{
			ID:    ph.ID,
			Value: formattedNumber,
		}); exists {
			deleteFromDb(db, ph.ID)
			continue
		}
		ph.Value = formattedNumber
		updateInDb(db, ph)
	}
}

func formatNumber(unformattedNum string) string {
	unformattedNumByteStr := []byte(unformattedNum)
	var result []byte
	for _, v := range unformattedNumByteStr {
		if v >= 48 && v <= 57 {
			result = append(result, v)
		}
	}
	return string(result)
}

func checkIfDuplicateExists(db *gorm.DB, phone phoneNumber) bool {
	rows, e := db.Raw(fmt.Sprintf("SELECT * from %s where number='%s' and id <> %d", tableName, phone.Value, phone.ID)).Rows()
	if e != nil {
		panic(fmt.Errorf("error fetching from db: %v", e))
	}
	var phoneFromDb phoneNumber
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&phoneFromDb.ID, &phoneFromDb.Value)
		return phone.Value == phoneFromDb.Value
	}
	return false
}

func deleteFromDb(db *gorm.DB, ID int) {
	db.Exec(fmt.Sprintf("DELETE from %s where id=%d", tableName, ID))
}

func updateInDb(db *gorm.DB, phone phoneNumber) {
	db.Exec(fmt.Sprintf("UPDATE %s set number='%s' where id=%d", tableName, phone.Value, phone.ID))
}
