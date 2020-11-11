package databae

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	driver     string
	dataSource string
)

// SQLi interface for wrapping sqlx.DB and sqlx.Tx
type SQLi interface {
	PreparedNamedContext(context.Context, string) (*sqlx.NamedStmt, error)
	MustExecContext(context.Context, string, ...interface{}) sql.Result
	SelectContext(context.Context, interface{}, string, ...interface{}) error
	GetContext(context.Context, interface{}, string, ...interface{}) error
	NamedExecContext(context.Context, string, interface{}) (sql.Result, error)
	sqlx.ExtContext
	sqlx.PreparerContext
}

// SqlxDB is struct for Sqlx Connection
type SqlxDB struct{}

//Register for registering database
func Register(driverName string, dataSourceName string) {
	driver = driverName
	dataSource = dataSourceName
}

// Open for open connection database
func Open() (*sql.DB, error) {
	return sql.Open(driver, dataSource)
}

func (d *SqlxDB) buildConnection() (*sqlx.DB, error) {
	db, err := Open()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("DB server is connected")

	return sqlx.NewDb(db, driver), nil
}

// GetSqlxConnection create connection for sqlx
func GetSqlxConnection() (*sqlx.DB, error) {
	db := SqlxDB{}
	return db.buildConnection()
}
