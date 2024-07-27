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

	// mysqlRepo "github.com/bxcodec/go-clean-arch/internal/repository/mysql"
	postgresRepo "github.com/bxcodec/go-clean-arch/internal/repository/postgresql"

	"github.com/bxcodec/go-clean-arch/article"
	"github.com/bxcodec/go-clean-arch/internal/rest"
	"github.com/bxcodec/go-clean-arch/internal/rest/middleware"
	"github.com/joho/godotenv"
)

const (
	defaultTimeout = 30
	defaultAddress = ":9090"
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
