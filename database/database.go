package database

import (
	"clinic_server/config"
	logs "clinic_server/logger"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func GetDBContext() *gorm.DB {
	logs.Init(zerolog.InfoLevel)

	db, err := gorm.Open(sqlite.Open(config.GetConfiguration().Database.ConnectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   true, // skip the snake_casing of names
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Panic().Err(err).Msg("failed to connect database")
		logs.Error("failed to connect database", err)
	}
	return db
}
