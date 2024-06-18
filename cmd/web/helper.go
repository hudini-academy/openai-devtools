package main

import (
	"database/sql"
	"log"
	"os"
	"regexp"
	"unicode"
	"fmt"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var printAst = false

// mdToHTML converts markdown to HTML using the gomarkdown/markdown package.
// It supports common extensions and rendering options, and can print the AST if printAst is true.
func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	if printAst {
		fmt.Print("--- AST tree:\n")
		ast.Print(os.Stdout, doc)
		fmt.Print("\n")
	}

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)

}
// initLogger initializes two loggers for info and error messages,
// writing to info.log and error.log respectively.
func initLogger() (*log.Logger, *log.Logger) {
    // Create a file for writing information messages.
    infoFile, err := os.OpenFile("./info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }

    // Create a file for writing error messages.
    errFile, err := os.OpenFile("./error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }

    // Create a logger for writing information messages.
    // The destination is the infoFile, the prefix is "INFO", and the flags include the date, time, and the message.
    infoLog := log.New(infoFile, "INFO\t", log.Ldate|log.Ltime)

    // Create a logger for writing error messages.
    // The destination is the errFile, the prefix is "ERROR", and the flags include the date, time, the message, and the file and line number where the error occurred.
    errorLog := log.New(errFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

    // Return the two loggers.
    return infoLog, errorLog
}

// openDB opens a database connection using the given DSN, checks the connection,
// and returns the *sql.DB instance or an error if the connection fails.
func openDB(dsn string) (*sql.DB, error) {
    // Open a database connection using the provided DSN
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        // Return nil and the error if the connection fails
        return nil, err
    }

    // Ping the database to check if the connection is successful
    if err = db.Ping(); err != nil {
        // Return nil and the error if the ping fails
        return nil, err
    }

    // Return the sql.DB instance and nil if both the connection and the ping are successful
    return db, nil
}


// isValidEmail checks if the given email address is valid.
// It uses a regular expression pattern to validate the email format.
func isValidEmail(email string) bool {
	// Regular expression for validating email addresses
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}


// isValidPassword checks if a password meets the required criteria.
// The password must contain at least 8 characters, including at least one uppercase letter,
// one lowercase letter, one digit, and one special character.
// If the password is valid, it returns nil. Otherwise, it returns an error.
func isValidPassword(password string) error {
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if len(password) < 8 || !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return fmt.Errorf("password must contain 8 characters including an uppercase, lowercase, digit and special character")
	}

	return nil
}
