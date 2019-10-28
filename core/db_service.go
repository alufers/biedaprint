package core

import (
	"fmt"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//DBService handles connecting to the database
type DBService struct {
	app *App
	DB  *gorm.DB
}

//NewDBService constructs a DBService
func NewDBService(app *App) *DBService {
	return &DBService{
		app: app,
	}
}

//Init initializes the service
func (dbs *DBService) Init() error {
	var err error
	dbs.DB, err = gorm.Open("sqlite3", filepath.Join(dbs.app.GetDataPath(), "db.sqlite"))
	if err != nil {
		return fmt.Errorf("connecting to db: %v", err)
	}
	return nil
}

//Close finishes the DB connection
func (dbs *DBService) Close() error {
	return dbs.DB.Close()
}
