package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getConnectionStr() string {
	var connectStr = "c:/tfs/Learning Projects/Clinic_1/Database/Clinic.db"
	return connectStr
}

func GetDBContext() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(getConnectionStr()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
