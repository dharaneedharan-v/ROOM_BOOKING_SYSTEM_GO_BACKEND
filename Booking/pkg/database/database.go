
package database

import (
	"database/sql"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DBConnector interface {
	EstablishConnection(dbURL string) (*Db, error)
}

type Db struct {
	Gorm  *gorm.DB
	SqlDb *sql.DB
}

type DBService struct{}

func (s *DBService) EstablishConnection(dbURL string) (*Db, error) {
	var dialector gorm.Dialector
	
	if strings.HasPrefix(dbURL, "sqlite://") {
		dbPath := strings.TrimPrefix(dbURL, "sqlite://")
		dialector = sqlite.Open(dbPath)
		// print("-------------" , dialector)
	} else if strings.HasPrefix(dbURL, "sqlserver://") {
		dialector = sqlserver.Open(dbURL)
	}else if strings.HasPrefix(dbURL,"postgresql://")  {    // Change 
			dialector = postgres.Open(dbURL)

	}else {
		dialector = sqlserver.Open(dbURL)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &Db{
		Gorm:  db,
		SqlDb: sqlDB,
	}, nil
}