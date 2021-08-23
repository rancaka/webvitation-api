package util

import (
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db     *sqlx.DB
	dbOnce sync.Once
)

func getDB() *sqlx.DB {

	dbOnce.Do(func() {
		var (
			err    error
			config = GetConfig()
		)

		db, err = sqlx.Connect(
			"mysql",
			config.Storage.DB.DataSourceName,
		)
		if err != nil {
			fmt.Printf("[BOOTING][ERROR] when trying to connect to DB: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("[BOOTING] successfully connected to Database!")
	})

	return db
}
