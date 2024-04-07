package main

import (
	"context"
	"database/sql"
	"fmt"
	db "music_vs_store/db/sqlc"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

const (
    dbDriver = "postgres"
    dbSource = "postgresql://admin:admin@localhost:1234/music_vs_store_db?sslmode=disable"
)

var queries *db.Queries

func main() {
  ctx := context.Background()

  conn, err := sql.Open(dbDriver, dbSource)
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  queries = db.New(conn)


  err = queries.DeleteUserByName(ctx, "aboba")
  if err != nil {
    panic(err)
  }
  err = queries.DeleteUserByName(ctx, "zeliboba")
  if err != nil {
    panic(err)
  }

  insertedUser, err := queries.CreateUser(ctx, db.CreateUserParams{
    Username: "aboba",
    Email: "aboba@gmail.com",
    Password: "password",
  })
  if err != nil {
    panic(err)
  }

  insertedUser2, err := queries.CreateUser(ctx, db.CreateUserParams{
    Username: "zeliboba",
    Email: "zeliboba@gmail.com",
    Password: "password",
  })
  if err != nil {
    panic(err)
  }

  fmt.Println("inserted user:", insertedUser)
  fmt.Println("inserted user 2:", insertedUser2)

  mux := http.NewServeMux()
  mux.HandleFunc("/", handler)

  fmt.Println("starting server at :8080")
  http.ListenAndServe(":8080", mux)
}


const templ = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Music Vs Store</title>
    <link rel="stylesheet" href="style.css">
  </head>
  <body>
    {{ range . }}
      <div style="margin-bottom: 10px;">
        <div style="font-size: 30px;">{{ .Username }}</div>
        <div style="color: red;">{{ .Email }}</div>
      </div>
    {{ end }}
  </body>
</html>
`

var tmpl = template.Must(template.New("aboba").Parse(templ))

func handler(w http.ResponseWriter, r *http.Request) {
  users, err := queries.ListUsrs(context.Background(), db.ListUsrsParams{
    Limit: 10,
    Offset: 0,
  })
  if err != nil {
    panic(err)
  }

  tmpl.Execute(w, users)
}
