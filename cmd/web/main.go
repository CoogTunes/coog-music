package main

import (
	"database/sql"
	"fmt"

	"github.com/CoogTunes/coog-music/internal/render"

	"github.com/CoogTunes/coog-music/internal/config"
	"github.com/CoogTunes/coog-music/internal/driver"
	"github.com/CoogTunes/coog-music/internal/handlers"

	// "github.com/DeLuci/coog-music/internal"
	// "github.com/DeLuci/coog-music/cmd/web/routers"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

var app config.AppConfig
var session *scs.SessionManager

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer func(SQL *sql.DB) {
		err := SQL.Close()
		if err != nil {
			panic(err)
		}
	}(db.SQL)

	fmt.Println("Connecting to localhost:8080!!!")
	srv := &http.Server{
		// Addr:    app.ServerAddress,
		Addr:    ":8080",
		Handler: routes(),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	log.Println("Connecting to database...")
	dbSource := config.GetDBEnvVariable("DB_SOURCE")
	db, err := driver.ConnectSQL(dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("connected to database!!!")

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	return db, nil
}
