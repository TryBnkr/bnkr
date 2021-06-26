package main

import (
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MohammedAl-Mahdawi/bnkr/app/dal"
	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"github.com/MohammedAl-Mahdawi/bnkr/utils"

	"strconv"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/app/services"
	"github.com/MohammedAl-Mahdawi/bnkr/config"
	"github.com/MohammedAl-Mahdawi/bnkr/config/database"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/middleware"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/password"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/render"

	"github.com/alexedwards/scs/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var app config.AppConfig
var session *scs.SessionManager
var errorLog *log.Logger

func main() {
	// what am I going to put in the session
	gob.Register(dal.User{})

	isProduction, _ := strconv.ParseBool(os.Getenv("PRODUCTION"))

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = isProduction

	app.Session = session
	app.Version = "1.0.4"

	app.Cron = cron.New()

	database.Connect()
	database.Migrate(&dal.User{}, &dal.Backup{}, &dal.Option{}, &dal.Job{}, &dal.Queue{})

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return
	}

	app.TemplateCache = tc

	// Mail setup
	mailChan := make(chan types.MailData, 100)
	app.MailChan = mailChan

	defer close(app.MailChan)

	fmt.Println("Staring mail listener...")
	listenForMail()

	// Queue setup
	queueChan := make(chan types.NewQueueDTO, 100)
	app.Queue = queueChan

	defer close(app.Queue)

	fmt.Println("Staring queue listener...")
	listenForQueues()

	// read flags
	useCache := flag.Bool("cache", isProduction, "Use template cache")

	flag.Parse()

	app.UseCache = *useCache

	if os.Getenv("SETUP") == "true" {
		setup()
	}

	app.CronIds = make(map[int]cron.EntryID)
	runJobs()

	repo := services.NewRepo(&app)
	services.NewHandlers(repo)

	middleware.NewMiddleware(&app)
	utils.NewUtils(&app)
	render.NewRenderer(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", os.Getenv("PORT")))

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func runJobs() {
	backups := &[]types.NewBackupDTO{}
	if err := dal.FindAllBackups(&backups).Error; err != nil {
		panic(err)
	}

	for _, b := range *backups {
		rBackup := b
		cron := render.ConstructCron(&rBackup)
		cronId, _ := app.Cron.AddFunc("CRON_TZ="+rBackup.Timezone+" "+cron, func() { services.Repo.CreateNewJob(&rBackup, false) })
		app.CronIds[int(rBackup.ID)] = cronId
	}

	app.Cron.Start()
}

func setup() {
	// Check if any user exists, if no user exists then create new one.
	result := dal.FindUser(&struct{ ID string }{})

	// If no user exist
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user := &dal.User{
			Name:     os.Getenv("USERNAME"),
			Password: password.Generate(os.Getenv("USERPASSWORD")),
			Email:    os.Getenv("USEREMAIL"),
		}

		// Create a user, if error return
		if err := dal.CreateUser(user); err.Error != nil {
			panic(err)
		}

		// Store the version
		o := &dal.Option{
			Name:  "VERSION",
			Value: app.Version,
		}

		if err := dal.CreateOption(o); err.Error != nil {
			panic(err)
		}
	}
}
