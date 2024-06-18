package main

// Import the Required packages
import (
	"database/sql"
	"log"
	"os"
	"regexp"
	"unicode"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"fmt"
)

var printAst = false

// mdToHTML converts markdown to HTML using github.com/gomarkdown/markdown and github.com/gomarkdown/markdown/html packages,
// with extensions: parser.CommonExtensions, parser.AutoHeadingIDs, parser.NoEmptyLineBeforeBlock,
// and flags: html.CommonFlags, html.HrefTargetBlank.
// Params: input byte slice of markdown content. Returns: byte slice of HTML content.
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

// initLogger initializes info and error loggers with file setup and formats;
// returns info and error loggers.
func initLogger() (*log.Logger, *log.Logger) {
	// Create variable infoFile for info log file
	infoFile, err := os.OpenFile("./info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Create variable errFile for error log file
	errFile, err := os.OpenFile("./error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Create info logger with prefix "INFO\t", date, time, and flags
	infoLog := log.New(infoFile, "INFO\t", log.Ldate|log.Ltime)

	// Create error logger with prefix "ERROR\t", date, time, file name and line number, and flags
	errorLog := log.New(errFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

// openDB opens database connection using provided DSN with sql.Open, 
// pings to verify connection; returns *sql.DB and error.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
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

// isValidPassword checks if password meets criteria: >=8 chars,
// 1 uppercase, 1 lowercase, 1 digit, 1 special char; returns nil if valid, error otherwise.
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
