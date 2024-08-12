package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/MuhammadSaim/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type applictaion struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {

	// Define a new command-line flag with the name 'addr', a default
	// value of ":4000"
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command-line flag for the MySQL DSN string
	dsn := flag.String("dsn", "muhammadsaim:muhammadsaim@/golang_snippetbox?parseTime=true", "MySQL data source name")

	// Importantly, we use flag parse func to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr variable
	flag.Parse()

	// use the slog.New func to initialize a new structured logger, which
	// writes to the standard out stream and use the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// To keep the main() func tidy I've put the code for creating a connection
	// pool into the separate openDB func below and we pass the DSN using command-line
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// We also defer a call to db.Close(), so that the connection pool is closes
	// beofre the main func exit
	defer db.Close()

	// Initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a decoder instance
	formDecoder := form.NewDecoder()

	// Initialize a new instance of our application struct, containing the
	// dependencies for the time being just adding our logger
	app := &applictaion{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	// use the Info() method to log the starting server message at info
	logger.Info("Starting server", "addr", *addr)

	serverErr := http.ListenAndServe(*addr, app.routes())

	// And we also use the Error() method to lag any error message returned by
	// http.ListenAndServe() at Error. End of that terminate the application with os.Exit
	logger.Error(serverErr.Error())

	// exit the application
	os.Exit(1)
}

// The openDB func wraps sql.Open() and return the sql.DB connection pool
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
