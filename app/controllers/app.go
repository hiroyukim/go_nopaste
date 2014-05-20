package controllers

import (
    "os"
	"log"
    "time"
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/robfig/revel"
    "go_nopaste/app/routes"
    "github.com/nu7hatch/gouuid"
)

type App struct {
	*revel.Controller
}

func dbpath() string {
	return os.Getenv("GOPATH") +  "/src/github.com/hiroyukim/go_nopaste/db/nopaste.db"
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Post(body string) revel.Result {
    uuid,_ := uuid.NewV4()
    ctime  := time.Now()

    var uid string = uuid.String()

	db, err := sql.Open( revel.Config.String("db.driver"), dbpath() )
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO entry(entry_id, body, ctime) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

    _, err = stmt.Exec(uuid.String(),body,ctime.Second())
    if err != nil {
        log.Fatal(err)
    }
	tx.Commit()

	return c.Redirect(routes.App.Show(uid))
}


func (c App) Show(uid string) revel.Result {
	db, err := sql.Open("sqlite3", dbpath() )

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    stmt, err := db.Prepare("SELECT body FROM entry WHERE entry_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

    var body string
    err = stmt.QueryRow(uid).Scan(&body)
    if err != nil {
        log.Fatal(err)
    }

    return c.Render(body)
}

