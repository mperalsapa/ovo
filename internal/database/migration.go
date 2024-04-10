package database

import (
	"fmt"
	"ovo-server/internal/model"
)

func Migrate() {
	fmt.Println("Checking database")

	fmt.Println("Doing migration...")
	db.connection.AutoMigrate(model.User{})
	db.connection.AutoMigrate(model.Movie{})
}
