package main

// Import required packages
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

// Initialize the Open AI Client
var openAIClient *OpenAIClient

// `application` struct holds application-specific data and services including responses, OpenAI client,
// user and custom GPT models, error and info loggers, and session information.
type application struct {
	Response     string
	OpenAIClient *OpenAIClient
	users        *mysql.UserModel
	CustomGPT    *mysql.CustomGPTModel
	errorLog     *log.Logger
	infoLog      *log.Logger
	session      *sessions.Session
}

// main function is the entry point of the application.
// It initializes the application, sets up the database connection,
// configures the HTTP server, and starts listening for incoming requests.
func main() {
	// addr flag is used to specify the network address the server should listen on.
	// Default value is ":4000".
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Make the db user and password to load from a config file.
	// Replace DBUser, DBPass, and DBName with actual values.
	dsn := fmt.Sprintf("%s:%s@/%s?parseTime=true", DBUser, DBPass, DBName)

	// secret flag is used to specify the secret key for session management.
	// Default value is "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge".
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret Key")

	// Parse the command-line flags.
	flag.Parse()

	// Initialize the info and error loggers.
	infoLog, errorLog := initLogger()

	// Open the database connection using the provided DSN.
	db, er := openDB(dsn)
	if er != nil {
		// Log the error and terminate the application if the database connection fails.
		errorLog.Fatal(er)
	}
	defer db.Close() // Close the database connection when the main function exits.

	// Create a new session manager with the provided secret key.
	// Set the session lifetime to 12 hours.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	// Initialize the OpenAI client with the provided API key.
	openAIClient = NewOpenAIClient(ApiKey)

	// Create a new application instance with the initialized services and loggers.
	app := &application{
		users:        &mysql.UserModel{DB: db},
		CustomGPT:    &mysql.CustomGPTModel{DB: db},
		errorLog:     errorLog,
		infoLog:      infoLog,
		session:      session,
		OpenAIClient: openAIClient,
	}

	// Create a new HTTP server with the specified address, handler, and error logger.
	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	// Print a message indicating that the server is starting.
	fmt.Println("starting server on ", *addr)

	// Start listening for incoming requests and handle errors.
	err := srv.ListenAndServe()
	if err != nil {
		// Log the error and terminate the application if the server fails to start.
		errorLog.Fatal(err)
	}
}
