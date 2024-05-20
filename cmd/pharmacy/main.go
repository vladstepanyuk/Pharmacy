package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"fyne.io/fyne/v2/app"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"pharmacy/interanl/pharmacy"
	"pharmacy/interanl/ui"
	"time"
)

func loadDBConn() (*sql.DB, error) {
	txtConf, ok := os.LookupEnv("PGCONN")
	if !ok {
		return nil, errors.New("PGCONN environment variable not set")
	}

	db, err := sql.Open("pgx", txtConf)
	if err != nil {
		return nil, err
	}

	return db, nil
}

var data = [][]string{[]string{"top left", "top right"},
	[]string{"bottom left", "bottom right"}}

var model *pharmacy.Model

func main() {

	db, err := loadDBConn()
	if err != nil {
		log.Fatal(err)
	}

	ctx, f := context.WithTimeout(context.Background(), time.Second)
	defer f()

	model, err = pharmacy.NewModel(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	defer model.Close()

	fmt.Println("ok")

	myApp := app.New()
	mainWindow := ui.NewMainWindow(myApp, model)
	mainWindow.ShowAndRun()
}
