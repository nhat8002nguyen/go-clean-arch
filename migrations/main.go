package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Replace with your SQL driver
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func connectMysql() (*sql.DB, error) {
	//prepare database
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")
	// mysql connection str
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	dbConn, err := sql.Open("mysql", dsn)
	return dbConn, err
}

func connectPostgres() (*sql.DB, error) {
	// Retrieve environment variables
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")

	// Construct the PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
	logrus.Info(connStr)

	// Open the database connection
	dbConn, err := sql.Open("postgres", connStr)

	return dbConn, err
}

func main() {
	dbDriver := os.Getenv("DATABASE_DRIVER")
	var dbConn *sql.DB
	var err error
	if dbDriver == "mysql" {
		dbConn, err = connectMysql()
	} else if dbDriver == "postgres" {
		dbConn, err = connectPostgres()
	}
	if err != nil {
		log.Fatal("failed to open connection to database", err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal("failed to ping database ", err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection", err)
		}
	}()

	log.Default().Printf("connected %s database", strings.ToUpper(os.Getenv("DATABASE_DRIVER")))

	if len(os.Args) < 2 {
		log.Fatal("please provide a goose command")
	}

	goose.SetDialect("postgres") // Set the dialect for your database

	ctx := context.Background()
	if err := goose.RunContext(ctx, os.Args[1], dbConn, "migrations"); err != nil {
		log.Fatalf("failed to run goose command: %v", err)
	}
}
