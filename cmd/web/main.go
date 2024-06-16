package main

import (
	"OpenAIDevTools/pkg/models/mysql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

var openAIClient *OpenAIClient

type application struct {
	Response     string
	OpenAIClient *OpenAIClient
	users        *mysql.UserModel
	CustomGPT  *mysql.CustomGPTModel
	errorLog     *log.Logger
	infoLog      *log.Logger
	session      *sessions.Session
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Make the db user and password to load from a config file

	dsn := fmt.Sprintf("%s:%s@/%s?parseTime=true", DBUser, DBPass, DBName)
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

	openAIClient = NewOpenAIClient(ApiKey)

	app := &application{
		users:        &mysql.UserModel{DB: db},
		CustomGPT:  &mysql.CustomGPTModel{DB: db},
		errorLog:     errorLog,
		infoLog:      infoLog,
		session:      session,
		OpenAIClient: openAIClient,
		infoLog:      infoLog,
		errorLog:     errorLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}
	fmt.Println("starting server on ", *addr)
	err := srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
