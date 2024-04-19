package model

import (
	"log"
	db "ovo-server/internal/database"
)

func Init() {
	log.Println("Automigrating database schema...")
	db.GetDB().AutoMigrate(&User{}, &Library{}, &Movie{}, &Show{})
}
