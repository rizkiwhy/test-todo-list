package database

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"rizkiwhy/test-todo-list/util/config"
)

func MySQLConnection() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config("MYSQL_USER"),
		config.Config("MYSQL_ROOT_PASSWORD"),
		config.Config("MYSQL_HOST"),
		config.Config("MYSQL_PORT"),
		config.Config("MYSQL_DATABASE_NAME"))

	fmt.Println("dsn: ", dsn)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Error().Err(err).Msg("[database][MySQLConnection] Failed to connect to MySQL")
		return nil, err
	}

	db.Logger = logger.Default.LogMode(logger.Info)

	return
}
