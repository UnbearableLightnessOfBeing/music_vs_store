package driver

import (
	"database/sql"
	db "music_vs_store/db/sqlc"
	"os"
)

func GetQueries() *db.Queries {
  conn, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_SOURCE"))
  if err != nil {
    panic(err)
  }

  err = conn.Ping()
  if err != nil {
    panic(err)
  }

  return db.New(conn)
}
