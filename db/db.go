package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

var Db *sql.DB //created outside to make it global.

// make sure your function start with uppercase to call outside of the directory.
func ConnectDatabase() {

	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.\n", k)
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep passwords and other secrets safe.
	var (
		dbUser                 = mustGetenv("POSTGRES_USER")               // e.g. 'my-db-user'
		dbPwd                  = mustGetenv("POSTGRES_PASSWORD")           // e.g. 'my-db-password'
		dbName                 = mustGetenv("POSTGRES_DB")                 // e.g. 'my-database'
		dbPort, _              = strconv.Atoi(mustGetenv("POSTGRES_PORT")) // e.g. 'my-database'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME")    // e.g. 'project:region:instance'
		usePrivate             = os.Getenv("POSTGRES_HOST")
	)

	dsn := fmt.Sprintf("user=%s password=%s database=%s, port=%d", dbUser, dbPwd, dbName, dbPort)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		fmt.Println("There is an error while connecting to the database ", err)
		panic(err)
	}
	var opts []cloudsqlconn.Option
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	}
	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		fmt.Println("There is an error while connecting to the cloudsql ", err)
		panic(err)
	}

	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName)
	}
	dbURI := stdlib.RegisterConnConfig(config)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		fmt.Println("sql.Open: %w", err)
		panic(err)
	}
	Db = dbPool

	// fmt.Printf("host=%s\nport=%d\nuser=%s\ndbname=%s\npass=%s", host, port, user, dbname, pass)

	// set up postgres sql to open it.
	// psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
	// 	host, port, user, dbname, pass)
	// psqlSetup := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s",
	// 	host, user, pass, port, dbname)

	// db, errSql := sql.Open("pgx", psqlSetup)
	// if errSql != nil {
	// 	fmt.Println("There is an error while connecting to the database ", errSql)
	// 	panic(errSql)
	// } else {
	// 	Db = db
	// 	fmt.Println("Successfully connected to database!")
	// }
}
