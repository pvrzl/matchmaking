package main

import (
	"app/internal/config"
	"app/pkg/database"
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func setup() (*migrate.Migrate, error) {
	cfg := config.Get()
	db, err := database.NewPostgresDB(cfg.Database.Write)
	if err != nil {
		return nil, fmt.Errorf("failed to start db %w", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to start driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", sqlDir), cfg.Database.Write.Name, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to start migration: %w", err)
	}

	return m, nil
}

func up() {
	m, err := setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer m.Close()

	prev, _, _ := m.Version()
	fmt.Println("previous version:", prev)

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println(err)
			return
		}
		if _, ok := err.(migrate.ErrDirty); ok {
			fmt.Println("please fix the dirty version first")
			return
		}
		fmt.Printf("failed while doing migration: %s, please use previous version(%d) to fix it\n", err, prev)
		return
	}
	v, _, _ := m.Version()
	fmt.Println("new version:", v)
}

func down() {
	m, err := setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer m.Close()

	prev, _, _ := m.Version()
	fmt.Println("previous version:", prev)
	if prev == 0 {
		fmt.Println("this is the lowest version that can be migrated down to")
		return
	}

	err = m.Steps(-1)
	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println(err)
			return
		}
		if _, ok := err.(migrate.ErrDirty); ok {
			fmt.Println("please fix the dirty version first")
			return
		}
		fmt.Printf("failed while doing migration: %s, please use previous version(%d) to fix it\n", err, prev)
		return
	}
	v, _, _ := m.Version()
	fmt.Println("new version:", v)
}

func fix() {
	fmt.Print("Enter the migration version: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Error reading input.")
		os.Exit(1)
	}

	migrationVersionStr := scanner.Text()
	migrationVersion, err := strconv.Atoi(migrationVersionStr)
	if err != nil {
		fmt.Println("failed to read migration version:", err)
		return
	}

	m, err := setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer m.Close()

	err = m.Force(migrationVersion)
	if err != nil {
		fmt.Println("failed while doing migration: ", err)
		return
	}

	v, _, _ := m.Version()
	fmt.Println("new version:", v)
}
