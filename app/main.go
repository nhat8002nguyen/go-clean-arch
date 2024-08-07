package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	// mysqlRepo "github.com/nhat8002nguyen/ecommerce-go-app/internal/repository/mysql"
	postgresRepo "github.com/nhat8002nguyen/ecommerce-go-app/internal/repository/postgresql"

	"github.com/joho/godotenv"
	"github.com/nhat8002nguyen/ecommerce-go-app/article"
	"github.com/nhat8002nguyen/ecommerce-go-app/internal/rest"
	"github.com/nhat8002nguyen/ecommerce-go-app/internal/rest/middleware"
)

const (
	defaultTimeout = 30
	defaultAddress = ":9090"
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
	//prepare database
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

	// prepare echo
	e := echo.New()
	e.Use(middleware.CORS)
	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Println("failed to parse timeout, using default timeout")
		timeout = defaultTimeout
	}
	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))

	// Prepare Repository
	var articleRepo article.ArticleRepository
	var authorRepo article.AuthorRepository
	if dbDriver == "mysql" {
		authorRepo = postgresRepo.NewAuthorRepository(dbConn)
		articleRepo = postgresRepo.NewArticleRepository(dbConn)
	} else if dbDriver == "postgres" {
		authorRepo = postgresRepo.NewAuthorRepository(dbConn)
		articleRepo = postgresRepo.NewArticleRepository(dbConn)
	}

	// Build service Layer
	svc := article.NewService(articleRepo, authorRepo)
	rest.NewArticleHandler(e, svc)

	// Start Server
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}
	log.Fatal(e.Start(address)) //nolint
}
