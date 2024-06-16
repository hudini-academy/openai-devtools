package main

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
func initLogger() (*log.Logger, *log.Logger) {
	//create variable f for info
	infoFile, err := os.OpenFile("./info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//create variable e for error.log
	errFile, err := os.OpenFile("./error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//created logger using log.New() for writing information message
	//destination,prefix message,flags for additional info(three parameters)
	infoLog := log.New(infoFile, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(errFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}
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
