package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func getConnectionStr() string {
	var connectStr = "c:/tfs/Learning Projects/Clinic_1/Database/Clinic.db"
	return connectStr
}

func GetDBContext() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(getConnectionStr()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   true, // skip the snake_casing of names
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
