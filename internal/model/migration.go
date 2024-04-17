package model

import (
	"log"
	db "ovo-server/internal/database"
)

func Init() {
	log.Println("Migrating database schema...")
	log.Println("Migrating User schema...")
	db.Migrate(&User{})
	log.Println("Migrating Movie schema...")
	db.Migrate(&Movie{})
}
