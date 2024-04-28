package main

import (
	"log"
	"ovo-server/internal/config"
	"ovo-server/internal/database"
	"ovo-server/internal/model"
)

func main() {
	config.Init()
	database.Init()

	db, err := database.GetDB().DB()
	if err != nil {
		panic(err)
	}

	log.Printf("Connected to database %s, host %s, user %s", config.Variables.DatabaseName, config.Variables.DatabaseHost, config.Variables.DatabaseUsername)
	if db.Ping() != nil {
		panic("Could not connect to database")
	}

	log.Println("Dropping all tables")

	// Get all table names
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			panic(err)
		}
		tables = append(tables, tableName)
	}

	// Drop all tables
	for _, tableName := range tables {
		_, err := db.Exec("DROP TABLE IF EXISTS `" + tableName + "`")
		if err != nil {
			panic(err)
		}
	}

	log.Println("Tables dropped")

	log.Printf("Creating database %s", config.Variables.DatabaseName)
	sqlResult, err := db.Exec("CREATE DATABASE IF NOT EXISTS `" + config.Variables.DatabaseName + "`;")
	if err != nil {
		panic(err)
	}
	log.Printf("Result: %v", sqlResult)

	log.Println("Migrating")
	model.Init()

	log.Println("Done")
}
