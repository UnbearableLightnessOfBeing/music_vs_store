package main

import (
	"context"
	"database/sql"
	"fmt"
	db "music_vs_store/db/sqlc"
	"net/http"
	"os"
	"text/template"

	"github.com/subosito/gotenv"

	_ "github.com/lib/pq"
)

var queries *db.Queries

func init() {
  gotenv.Load()
}
func main() {
  // ctx := context.Background()

  conn, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_SOURCE"))
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  queries = db.New(conn)

  // err = queries.DeleteUserByName(ctx, "aboba")
  // if err != nil {
  //   panic(err)
  // }
  // err = queries.DeleteUserByName(ctx, "zeliboba")
  // if err != nil {
  //   panic(err)
  // }

  // insertedUser, err := queries.CreateUser(ctx, db.CreateUserParams{
  //   Username: "aboba",
  //   Email: "aboba@gmail.com",
  //   Password: "password",
  // })
  // if err != nil {
  //   panic(err)
  // }

  // insertedUser2, err := queries.CreateUser(ctx, db.CreateUserParams{
  //   Username: "zeliboba",
  //   Email: "zeliboba@gmail.com",
  //   Password: "password",
  // })
  // if err != nil {
  //   panic(err)
  // }

  // fmt.Println("inserted user:", insertedUser)
  // fmt.Println("inserted user 2:", insertedUser2)

  mux := http.NewServeMux()
  mux.HandleFunc("/", handler)

  port := os.Getenv("SERVER_PORT")

  fmt.Println("starting server at " + port)
  http.ListenAndServe(port, mux)
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
