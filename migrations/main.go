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

func isEcsEnvironment() bool {
	_, exists := os.LookupEnv("AWS_EXECUTION_ENV")
	return exists
}

func isRunningInDocker() bool {
	// Check for the existence of /.dockerenv file
	_, err := os.Stat("/.dockerenv")
	if err == nil {
		return true
	}

	// Check if cgroups indicate a containerized environment
	cgroup, err := os.ReadFile("/proc/self/cgroup")
	if err != nil {
		return false
	}

	return strings.Contains(string(cgroup), "docker")
}

func init() {
	if isEcsEnvironment() || isRunningInDocker() {
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func connectMysql() (*sql.DB, error) {
	// Retrieve environment variables
	dbHost, ok := os.LookupEnv("DATABASE_HOST")
	if !ok {
		log.Fatal("missing env var DATABASE_HOST")
	}
	dbPort, ok := os.LookupEnv("DATABASE_PORT")
	if !ok {
		log.Fatal("missing env var DATABASE_PORT")
	}
	dbUser, ok := os.LookupEnv("DATABASE_USER")
	if !ok {
		log.Fatal("missing env var DATABASE_USER")
	}
	dbPass, ok := os.LookupEnv("DATABASE_PASS")
	if !ok {
		log.Fatal("missing env var DATABASE_USER")
	}
	dbName, ok := os.LookupEnv("DATABASE_NAME")
	if !ok {
		log.Fatal("missing env var DATABASE_NAME")
	}

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
	dbHost, ok := os.LookupEnv("DATABASE_HOST")
	if !ok {
		log.Fatal("missing env var DATABASE_HOST")
	}
	dbPort, ok := os.LookupEnv("DATABASE_PORT")
	if !ok {
		log.Fatal("missing env var DATABASE_PORT")
	}
	dbUser, ok := os.LookupEnv("DATABASE_USER")
	if !ok {
		log.Fatal("missing env var DATABASE_USER")
	}
	dbPass, ok := os.LookupEnv("DATABASE_PASS")
	if !ok {
		log.Fatal("missing env var DATABASE_USER")
	}
	dbName, ok := os.LookupEnv("DATABASE_NAME")
	if !ok {
		log.Fatal("missing env var DATABASE_NAME")
	}

	// Construct the PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
	logrus.Info(connStr)

	// Open the database connection
	dbConn, err := sql.Open("postgres", connStr)

	return dbConn, err
}

func main() {
	dbDriver, ok := os.LookupEnv("DATABASE_DRIVER")
	if !ok {
		log.Fatal("missing env var DATABASE_DRIVER")
	}
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
	if err := goose.RunContext(ctx, os.Args[1], dbConn, "."); err != nil {
		log.Fatalf("failed to run goose command: %v", err)
	}
}
