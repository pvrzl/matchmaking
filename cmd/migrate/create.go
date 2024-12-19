package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

func sanitizeMigrationName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	reg, _ := regexp.Compile("[^a-zA-Z0-9_]+")
	name = reg.ReplaceAllString(name, "")
	return name
}

func create() {
	fmt.Print("Enter the migration name: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Error reading input.")
		return
	}

	migrationName := scanner.Text()
	sanitizedMigrationName := sanitizeMigrationName(migrationName)

	timestamp := time.Now().Format("20060102150405")

	dir := sqlDir
	upFile := fmt.Sprintf("%s/%s_%s.up.sql", dir, timestamp, sanitizedMigrationName)
	downFile := fmt.Sprintf("%s/%s_%s.down.sql", dir, timestamp, sanitizedMigrationName)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	if _, err := os.Create(upFile); err != nil {
		fmt.Printf("Error creating .up.sql file: %v\n", err)
		return
	}

	if _, err := os.Create(downFile); err != nil {
		fmt.Printf("Error creating .down.sql file: %v\n", err)
		return
	}

	fmt.Println("Migration files created successfully:")
	fmt.Println(upFile)
	fmt.Println(downFile)
}
