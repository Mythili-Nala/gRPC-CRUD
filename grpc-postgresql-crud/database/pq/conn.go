package pq

import (
	"fmt"
	"log"
	"os"

	//database "grpc-postgre2/database"
	"strconv"

	database "github.com/alifudin-a/grpc-pg-crud/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	host := fmt.Sprintf(os.Getenv("db_host"))
	user := fmt.Sprintf(os.Getenv("db_user"))
	dbname := fmt.Sprintf(os.Getenv("db_name"))
	port, _ := strconv.Atoi(os.Getenv("db_port"))
	sslmode := "disable"

	urlConnection := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s", host, port, user, dbname, sslmode)

	log.Println("Connecting to Database Server " + fmt.Sprint(os.Getenv("db_host")) + ":" + fmt.Sprint(os.Getenv("db_port")) + "...")

	database.Register("postgres", urlConnection)
}
