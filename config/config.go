package config

import (
	"html/template"
	"log"

	"github.com/MohammedAl-Mahdawi/bnkr/app/types"

	"github.com/alexedwards/scs/v2"
	"github.com/robfig/cron/v3"
)

// AppConfig holds the application config
type AppConfig struct {
	Session       *scs.SessionManager
	TemplateCache map[string]*template.Template
	ErrorLog      *log.Logger
	UseCache      bool
	Version       string
	DbUri         string
	MailChan      chan types.MailData
	Queue         chan types.NewQueueDTO
	Cron          *cron.Cron
	CronIds       map[int]cron.EntryID
}
