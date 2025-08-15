package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Word struct {
	ID   uint   `gorm:"primaryKey"`
	Word string `gorm:"unique"`
}

type URL struct {
	ID       uint   `gorm:"primaryKey"`
	URL      string `gorm:"unique"`
	NumWords int
}

type WordCount struct {
	ID     uint `gorm:"primaryKey"`
	Count  int
	WordID uint
	URLID  uint

	Word Word `gorm:"foreignKey:WordID"`
	URL  URL  `gorm:"foreignKey:URLID"`
}

func database() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error creating the database: ", err)
	}

	db.Exec("PRAGMA foreign_keys = ON;")

	err = db.AutoMigrate(&Word{}, &URL{}, &WordCount{})
	if err != nil {
		fmt.Println("Error creating tables: ", err)
	}

	return db
}
