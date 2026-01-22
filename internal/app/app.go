package app

import (
	"database/sql"
	"log"
	"os"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

func NewApp() *pocketbase.PocketBase {
	sqlite_vec.Auto()

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DBConnect: func(dbPath string) (*dbx.DB, error) {
			db, err := sql.Open("sqlite3", dbPath)
			if err != nil {
				return nil, err
			}

			// Verify sqlite-vec registration
			var vecVer string
			err = db.QueryRow("SELECT vec_version()").Scan(&vecVer)
			if err != nil {
				log.Printf("vec_version() failed in DBConnect: %v (extension might not be loaded correctly)", err)
			} else {
				log.Printf("sqlite-vec successfully registered: %s", vecVer)
			}

			_, _ = db.Exec("PRAGMA journal_mode=WAL;")

			return dbx.NewFromDB(db, "sqlite3"), nil
		},
	})

	return app
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
