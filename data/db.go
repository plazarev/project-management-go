package data

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Debug = true

type DBConfig struct {
	Path  string
	Debug bool
	Reset bool
}

var db *gorm.DB

func Init(config DBConfig) {
	if config.Path == "" {
		panic("undefined path to the database")
	}

	var err error
	level := logger.Error
	if config.Debug {
		level = logger.Warn
	}

	db, err = gorm.Open(sqlite.Open(config.Path), &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})
	if err != nil {
		fmt.Printf("unable to connect to database: %s\n%s\n", config.Path, err)
		os.Exit(1)
	}

	Debug = config.Debug

	if err != nil {
		panic("failed to connect database")
	}

	must(db.AutoMigrate(&Item{}))
	must(db.AutoMigrate(&User{}))
	must(db.AutoMigrate(&Project{}))
	must(db.AutoMigrate(&ItemUser{}))
	must(db.AutoMigrate(&File{}))
	must(db.AutoMigrate(&Comment{}))
}

func GetDB() *gorm.DB {
	return db
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
