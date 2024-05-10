package main

import (
	"ovo-server/internal/config"
	"ovo-server/internal/database"
	"ovo-server/internal/model"
)

func main() {
	config.Init()
	database.Init()
	var libraries []model.Library

	database.GetDB().Find(&libraries)

	for _, library := range libraries {
		for _, item := range library.GetItems() {
			// Fetch runtime from file
			if item.ItemType == model.ItemTypeMovie || item.ItemType == model.ItemTypeShow {
				item.UpdateItemRuntime()
			}
		}
	}
}
