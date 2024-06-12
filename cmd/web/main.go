package main

import (
	"OpenAIDevTools/pkg/models/mysql"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/golangcollege/sessions"
	_ "github.com/go-sql-driver/mysql"
)

var client *OpenAIClient

func createConfig() {
	// Create OpenAI client
	client = NewOpenAIClient(ApiKey)
}

type application struct {
	Response string
	users    *mysql.UserModel
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
}

func main() {
	createConfig()
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := "root:root@/openaiusers?parseTime=true"
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret Key")
	flag.Parse()

	infoLog, errorLog := initLogger()
	db, er := openDB(dsn)
	if er != nil {
		errorLog.Fatal(er)
	}
	defer db.Close()
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		users: &mysql.UserModel{DB: db},
		session: session,

	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}
	infoLog.Println("starting server on :4000", addr)
	err := srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
