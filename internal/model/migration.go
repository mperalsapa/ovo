package model

import (
	"fmt"
	db "ovo-server/internal/database"
)

func Init() {
	fmt.Println("migrating from model")
	db.Migrate(&User{})
	db.Migrate(&Movie{})
}
