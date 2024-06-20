package context

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type Db struct {
	db     *sql.DB
	config DbConfig
}

func NewDb(config DbConfig) *Db {
	return &Db{
		config: config,
	}
}

func (d *Db) Init() {
	log.Printf("Connecting to database %s:%s as user %s to database %s\n",
		d.config.Host, d.config.Port, d.config.User, d.config.Name)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.config.User, d.config.Password, d.config.Host, d.config.Port, d.config.Name, d.config.SSLMode)

	var err error
	d.db, err = sql.Open("postgres", connStr)

	if err == nil {
		err = d.migrate()
	}

	if err != nil {
		log.Fatal("Failed to initialize DB", err)
	}
}

func (d *Db) Close() {
	err := d.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (d *Db) Get() *sql.DB {
	return d.db
}

func (d *Db) migrate() error {
	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not start migration: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run up migrations: %w", err)
	}
	log.Println("Migrations ran successfully")

	return nil
}
