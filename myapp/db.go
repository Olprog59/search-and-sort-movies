package myapp

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func StartDb() {
	db, err := gorm.Open("sqlite3", "app.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Movie{}, &Serie{}, &Season{}, &File{})
}

func testDb(f func(dataBase *gorm.DB)) {
	db, err := gorm.Open("sqlite3", "app.db")
	db.LogMode(true)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	f(db)
}
