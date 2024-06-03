package mailer

import (
	"plex-tvtime-sync/pkg/lib"
	"text/template"

	"go.uber.org/fx"
)

type smtpMailer struct {
	logger    lib.Logger
	config    lib.Config
	templates *template.Template
}

type SMTPMailerParams struct {
	fx.In

	Logger lib.Logger
	Config lib.Config
}
