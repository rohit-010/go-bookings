package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rohit-010/go-bookings/internal/config"
	"github.com/rohit-010/go-bookings/internal/driver"
	"github.com/rohit-010/go-bookings/internal/handlers"
	"github.com/rohit-010/go-bookings/internal/helpers"
	"github.com/rohit-010/go-bookings/internal/models"
	"github.com/rohit-010/go-bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8081"

var app config.AppConfig
var session *scs.SessionManager // variable shadowing if below var is used with :=
// example session := scs.New() but we r using session = scs.New()

// main application function
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database
	log.Println("Connecting to database....")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=rohit password=")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	log.Println("Connected to database....")
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err)
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}

// addValues adds two integers and returns a sum
// func addValues(x, y int) int {
// 	return x + y
// }

// func Divide(w http.ResponseWriter, r *http.Request) {
// 	f, err := divideValues(100.0, 0)
// 	if err != nil {
// 		fmt.Fprintf(w, "Cannot divide by 0")
// 		return
// 	}

// 	fmt.Fprintf(w, fmt.Sprintf("%f divided by %f is %f", 100.0, 100.0, f))
// }

// func divideValues(x, y float32) (float32, error) {
// 	if y <= 0 {
// 		err := errors.New("cannot divide by 0")
// 		return 0, err
// 	}
// 	result := x / y
// 	return result, nil
// }
